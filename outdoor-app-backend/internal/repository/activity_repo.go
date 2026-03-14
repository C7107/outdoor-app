package repository

import (
	"errors"
	"outdoor-app-backend/internal/database"
	"outdoor-app-backend/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateActivity(activity *model.Activity) error {
	return database.DB.Create(activity).Error
}

// GetActivityList 获取瀑布流列表 (多条件组合查询)
func GetActivityList(city string, minFitness int, status string, ageLimit string, page, pageSize int) ([]model.Activity, error) {
	var activities []model.Activity
	query := database.DB.Model(&model.Activity{})

	if city != "" {
		query = query.Where("city = ?", city)
	}
	if minFitness > 0 {
		query = query.Where("min_fitness <= ?", minFitness)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if ageLimit != "" {
		query = query.Where("age_limit = ?", ageLimit)
	}

	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Preload("Initiator").
		Find(&activities).Error
	return activities, err
}

func GetActivityByID(id uint) (*model.Activity, error) {
	var activity model.Activity
	err := database.DB.Preload("Initiator").First(&activity, id).Error
	return &activity, err
}

func DeleteActivity(id uint) error {
	return database.DB.Delete(&model.Activity{}, id).Error
}

// =====================================
// 🌟 核心难点：带【悲观锁】的审核机制，防并发超卖
// =====================================
func AuditMember(memberID uint, status string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 获取申请记录
		var member model.ActivityMember
		if err := tx.First(&member, memberID).Error; err != nil {
			return err
		}

		if status == "approved" {
			var activity model.Activity
			// 🔒 悲观锁：clause.Locking{Strength: "UPDATE"} 会在底层执行 SELECT ... FOR UPDATE
			// 这意味着同一瞬间，其他企图审核该活动的人，必须排队等我这个事务执行完！
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&activity, member.ActivityID).Error; err != nil {
				return err
			}
			//.Clauses() 是 GORM 提供的方法，用于向生成的 SQL 语句中添加自定义的子句（clause）。它可以接受多个实现了 clause.Interface 的参数
			//clause.Locking 是 GORM 定义的锁子句结构体，用于表示数据库的锁选项
			//Strength: "UPDATE" 指定锁的强度为 UPDATE

			// 2. 严格校验名额 (由于锁住了，此时的 CurrentMembers 是绝对准确的)
			if activity.CurrentMembers >= activity.MaxMembers {
				return errors.New("超卖拦截：名额已满，无法通过")
			}

			// 3. 人数 + 1
			newCount := activity.CurrentMembers + 1
			updateData := map[string]interface{}{"current_members": newCount}

			// 4. 如果加上这1个刚好满了，自动把活动改状态
			if newCount >= activity.MaxMembers {
				updateData["status"] = "full"
			}

			if err := tx.Model(&activity).Updates(updateData).Error; err != nil {
				return err
			}
		}

		// 5. 更新申请记录状态为通过/拒绝
		if err := tx.Model(&member).Update("status", status).Error; err != nil {
			return err
		}

		return nil
	})
} //这里的锁是自动释放的，无需显式编写解锁代码。锁的生命周期与数据库事务完全绑定：事务提交（COMMIT）或回滚（ROLLBACK）时，该事务持有的所有锁都会被自动释放

// 模糊搜索Title和Destination（支持分页）
func SearchActivities(keyword string, page, pageSize int) ([]model.Activity, error) {
	var activities []model.Activity

	like := "%" + keyword + "%"

	err := database.DB.
		Model(&model.Activity{}).
		Where("title LIKE ? OR destination LIKE ?", like, like).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Preload("Initiator").
		Find(&activities).Error

	return activities, err
}

// 报名行为
func CreateMember(member *model.ActivityMember) error {
	return database.DB.Create(member).Error
}

// 获取我发布的活动（分页）
func GetActivitiesByInitiator(userID uint, page, pageSize int) ([]model.Activity, error) {
	var activities []model.Activity

	err := database.DB.
		Where("initiator_id = ?", userID).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Preload("Initiator").
		Find(&activities).Error

	return activities, err
}

// 查询某个活动的所有成员（分页）
func GetMembersByActivityID(activityID uint, page, pageSize int) ([]model.ActivityMember, error) {
	var members []model.ActivityMember

	err := database.DB.
		Where("activity_id = ?", activityID).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Preload("User").
		Find(&members).Error

	return members, err
}

// 根据状态筛选成员（分页）
func GetMembersByActivityAndStatus(activityID uint, status string, page, pageSize int) ([]model.ActivityMember, error) {
	var members []model.ActivityMember

	query := database.DB.Where("activity_id = ?", activityID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Preload("User").
		Find(&members).Error

	return members, err
}

// 根据ID获取报名记录
func GetMemberByID(id uint) (*model.ActivityMember, error) {
	var member model.ActivityMember

	err := database.DB.
		Preload("User").
		First(&member, id).Error

	return &member, err
}

// 更新报名状态
func UpdateMemberStatus(id uint, status string) error {
	return database.DB.
		Model(&model.ActivityMember{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// 活动参与人数 +1
func IncreaseActivityMembers(activityID uint) error {
	return database.DB.
		Model(&model.Activity{}).
		Where("id = ?", activityID).
		Update("current_members", gorm.Expr("current_members + 1")).Error
}

// 根据用户ID查询报名记录（防止重复报名）
func GetMemberByUserAndActivity(userID uint, activityID uint) (*model.ActivityMember, error) {
	var member model.ActivityMember

	err := database.DB.
		Where("user_id = ? AND activity_id = ?", userID, activityID).
		First(&member).Error

	if err != nil {
		return nil, err
	}

	return &member, nil
}

// 查询用户参与的活动（分页）
func GetActivitiesJoinedByUser(userID uint, page, pageSize int) ([]model.ActivityMember, error) {
	var members []model.ActivityMember

	err := database.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Preload("Activity").
		Find(&members).Error

	return members, err
}
