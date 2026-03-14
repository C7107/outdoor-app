package service

import (
	"encoding/json"
	"outdoor-app-backend/internal/dto"
	"outdoor-app-backend/internal/model"
	"outdoor-app-backend/internal/repository"
)

// GetUserProfile 获取个人资料
func GetUserProfile(userID uint) (*dto.UserProfileRes, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserProfileRes{
		ID:               user.ID,
		Email:            user.Email,
		Nickname:         user.Nickname,
		Avatar:           user.Avatar,
		Signature:        user.Signature,
		FitnessLevel:     user.FitnessLevel,
		EmergencyContact: user.EmergencyContact,
	}, nil
}

func GetUserPasswordbyID(userID uint) (string, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return "", err
	}

	return user.PasswordHash, nil
}

func UpdateUserPasswordbyID(userID uint, newPassword string) error {
	return repository.UpdateUser(userID, map[string]interface{}{"password_hash": newPassword})
}

// UpdateUserProfile 修改个人资料
func UpdateUserProfile(userID uint, req *dto.UserUpdateReq) error {

	updateData := map[string]interface{}{}

	if req.Nickname != "" {
		updateData["nickname"] = req.Nickname
	}

	if req.Avatar != "" {
		updateData["avatar"] = req.Avatar
	}

	if req.Signature != "" {
		updateData["signature"] = req.Signature
	}

	if req.FitnessLevel != 0 {
		updateData["fitness_level"] = req.FitnessLevel
	}

	if req.EmergencyContact != "" {
		updateData["emergency_contact"] = req.EmergencyContact
	}

	return repository.UpdateUser(userID, updateData)
}

// GetMyPublished 获取我发布的活动瀑布流
func GetMyPublished(userID uint, page, pageSize int) ([]dto.ActivityDetailRes, error) {
	activities, err := repository.GetActivitiesByInitiator(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	var res []dto.ActivityDetailRes
	for _, a := range activities {
		var images []string
		json.Unmarshal([]byte(a.Images), &images)

		res = append(res, dto.ActivityDetailRes{
			ID:             a.ID,
			Title:          a.Title,
			City:           a.City,
			Destination:    a.Destination,
			GatherTime:     a.GatherTime,
			GatherPlace:    a.GatherPlace,
			FeeType:        a.FeeType,
			GroupQrCode:    a.GroupQrCode, // 发起人自己看自己的活动，一定有二维码
			CoverImage:     a.CoverImage,
			Description:    a.Description,
			Images:         images,
			MinFitness:     a.MinFitness,
			AgeLimit:       a.AgeLimit,
			MaxMembers:     a.MaxMembers,
			CurrentMembers: a.CurrentMembers,
			Status:         a.Status,
			Initiator: dto.UserBasicInfo{
				ID:       a.Initiator.ID,
				Nickname: a.Initiator.Nickname,
				Avatar:   a.Initiator.Avatar,
			},
			HasApplied:  true,
			ApplyStatus: "approved", // 发起人自己就是默认通过的
		})
	}
	return res, nil
}

// 🌟 核心逻辑：获取我参与的活动瀑布流
// 这里底层查出来的是 []model.ActivityMember (报名记录)
func GetMyJoined(userID uint, page, pageSize int) ([]dto.ActivityDetailRes, error) {
	members, err := repository.GetActivitiesJoinedByUser(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	var res []dto.ActivityDetailRes
	for _, m := range members {
		a := m.Activity // 取出关联的活动信息

		var images []string
		json.Unmarshal([]byte(a.Images), &images)

		// 脱敏逻辑：只有当该条报名记录的状态是 approved 时，才下发二维码
		qrCode := ""
		if m.Status == "approved" {
			qrCode = a.GroupQrCode
		}

		res = append(res, dto.ActivityDetailRes{
			ID:             a.ID,
			Title:          a.Title,
			City:           a.City,
			Destination:    a.Destination,
			GatherTime:     a.GatherTime,
			GatherPlace:    a.GatherPlace,
			FeeType:        a.FeeType,
			GroupQrCode:    qrCode,
			CoverImage:     a.CoverImage,
			Description:    a.Description,
			Images:         images,
			MinFitness:     a.MinFitness,
			AgeLimit:       a.AgeLimit,
			MaxMembers:     a.MaxMembers,
			CurrentMembers: a.CurrentMembers,
			Status:         a.Status,
			// 发起人信息如果不强关联，可以不用查完整，或者底层加上 Preload("Activity.Initiator")
			Initiator: dto.UserBasicInfo{
				ID: a.InitiatorID,
			},
			HasApplied:  true,
			ApplyStatus: m.Status, // 👈 最关键的字段：前端靠这个显示 "待审核" 或 "已通过"
		})
	}
	return res, nil
}

// GetMyFavorites 获取我的线路收藏
func GetMyFavorites(userID uint) ([]model.Route, error) {
	return repository.GetUserFavorites(userID)
}

// GetMyMessages 获取我的通知
func GetMyMessages(userID uint, page, pageSize int) ([]model.Message, error) {
	return repository.GetUserMessages(userID, page, pageSize)
}
