package service

import (
	"encoding/json"
	"errors"
	"outdoor-app-backend/internal/dto"
	"outdoor-app-backend/internal/model"
	"outdoor-app-backend/internal/repository"
	"outdoor-app-backend/pkg/ws"
)

// GetActivityList 获取瀑布流活动列表 (多条件组合筛选)
func GetActivityList(req *dto.ActivityFilterReq, page, pageSize int) ([]dto.ActivityDetailRes, error) {
	activities, err := repository.GetActivityList(req.City, req.MinFitness, req.Status, req.AgeLimit, page, pageSize)
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
			GroupQrCode:    "", // 列表页绝对不返回二维码
			CoverImage:     a.CoverImage,
			Description:    a.Description,
			Images:         images,
			MinFitness:     a.MinFitness,
			AgeLimit:       a.AgeLimit,
			MaxMembers:     a.MaxMembers,
			CurrentMembers: a.CurrentMembers,
			Status:         a.Status,
			WeatherAlert:   a.WeatherAlert,
			Initiator: dto.UserBasicInfo{
				ID:       a.Initiator.ID,
				Nickname: a.Initiator.Nickname,
				Avatar:   a.Initiator.Avatar,
			},
		})
	}
	return res, nil
}

// GetActivityDetail 获取活动详情 (包含隐私隔离逻辑)
func GetActivityDetail(userID, activityID uint) (*dto.ActivityDetailRes, error) {
	activity, err := repository.GetActivityByID(activityID)
	if err != nil {
		return nil, errors.New("活动不存在")
	}

	var images []string
	json.Unmarshal([]byte(activity.Images), &images)

	// 1. 判断当前登录用户的身份和报名状态
	var applyStatus string
	hasApplied := false

	member, err := repository.GetMemberByUserAndActivity(userID, activityID)
	if err == nil {
		hasApplied = true
		applyStatus = member.Status
	}

	// 2. 🌟 隐私脱敏逻辑：控制二维码的展示
	// 默认不展示
	qrCode := ""
	// 只有这两种情况才展示二维码：
	// 情况A: 你是发起人本人
	// 情况B: 你报名了，并且审核已通过 (approved)
	if activity.InitiatorID == userID || applyStatus == "approved" {
		qrCode = activity.GroupQrCode
	}

	// 3. 组装返回
	res := &dto.ActivityDetailRes{
		ID:             activity.ID,
		Title:          activity.Title,
		City:           activity.City,
		Destination:    activity.Destination,
		GatherTime:     activity.GatherTime,
		GatherPlace:    activity.GatherPlace,
		FeeType:        activity.FeeType,
		GroupQrCode:    qrCode, // 脱敏后的二维码
		CoverImage:     activity.CoverImage,
		Description:    activity.Description,
		Images:         images,
		MinFitness:     activity.MinFitness,
		AgeLimit:       activity.AgeLimit,
		MaxMembers:     activity.MaxMembers,
		CurrentMembers: activity.CurrentMembers,
		Status:         activity.Status,
		WeatherAlert:   activity.WeatherAlert,
		Initiator: dto.UserBasicInfo{
			ID:       activity.Initiator.ID,
			Nickname: activity.Initiator.Nickname,
			Avatar:   activity.Initiator.Avatar,
		},
		HasApplied:  hasApplied,
		ApplyStatus: applyStatus,
	}

	return res, nil
}

// CreateActivity 发起活动 (所有人均可发起)
func CreateActivity(userID uint, req *dto.ActivityCreateReq) error {
	imagesJSON, _ := json.Marshal(req.Images)

	activity := &model.Activity{
		InitiatorID:    userID,
		Title:          req.Title,
		City:           req.City,
		Destination:    req.Destination,
		GatherTime:     req.GatherTime,
		GatherPlace:    req.GatherPlace,
		FeeType:        req.FeeType,
		GroupQrCode:    req.GroupQrCode,
		CoverImage:     req.CoverImage,
		Description:    req.Description,
		Images:         string(imagesJSON),
		MinFitness:     req.MinFitness,
		AgeLimit:       req.AgeLimit,
		MaxMembers:     req.MaxMembers,
		CurrentMembers: 1, // 发起人自己占1个名额
		Status:         "enrolling",
	}
	return repository.CreateActivity(activity)
}

// DeleteActivity 删除活动 (只能发起人删除)
func DeleteActivity(userID, activityID uint) error {
	activity, err := repository.GetActivityByID(activityID)
	if err != nil {
		return errors.New("活动不存在")
	}

	// 越权校验
	if activity.InitiatorID != userID {
		return errors.New("无权删除他人的活动")
	}

	return repository.DeleteActivity(activityID)
}

// ApplyActivity 用户报名活动 (带各种业务拦截)
func ApplyActivity(userID, activityID uint, req *dto.ActivityApplyReq) error {
	activity, err := repository.GetActivityByID(activityID)
	if err != nil {
		return errors.New("活动不存在")
	}

	if activity.Status != "enrolling" {
		return errors.New("活动当前不可报名")
	}
	if activity.InitiatorID == userID {
		return errors.New("不能报名自己发起的活动")
	}
	if activity.CurrentMembers >= activity.MaxMembers {
		return errors.New("名额已满")
	}

	if _, err := repository.GetMemberByUserAndActivity(userID, activityID); err == nil {
		return errors.New("您已提交过申请，请勿重复报名")
	}

	user, _ := repository.GetUserByID(userID)
	if user.FitnessLevel < activity.MinFitness {
		return errors.New("您的体能等级未达到门槛")
	}

	member := &model.ActivityMember{
		ActivityID:    activityID,
		UserID:        userID,
		EquipmentNote: req.EquipmentNote,
		Status:        "pending",
	}

	msg := &model.Message{
		UserID:    activity.InitiatorID,
		Type:      "system",
		Title:     "新报名通知",
		Content:   "有新用户申请加入【" + activity.Title + "】，请审核。",
		RelatedID: activityID, // 🌟 把活动 ID 传过去
	}

	repository.CreateMessage(msg)
	repository.CreateMember(member)
	ws.SendMessageToUser(activity.InitiatorID, msg)

	return nil
}

// AuditActivityMember 发起人审核报名
func AuditActivityMember(initiatorID uint, req *dto.ActivityAuditReq) error {
	member, err := repository.GetMemberByID(req.MemberID)
	if err != nil {
		return errors.New("申请记录不存在")
	}

	activity, err := repository.GetActivityByID(member.ActivityID)
	if err != nil {
		return errors.New("活动不存在")
	}

	if activity.InitiatorID != initiatorID {
		return errors.New("无权审核此活动")
	}
	if member.Status != "pending" {
		return errors.New("该申请已处理")
	}

	// 执行审核
	if err := repository.AuditMember(req.MemberID, req.Status); err != nil {
		return errors.New("审核操作失败: " + err.Error())
	}

	resultStr := "已通过"
	if req.Status == "rejected" {
		resultStr = "已被婉拒"
	}

	// 构造消息
	msg := &model.Message{
		UserID:  member.UserID,
		Type:    "activity",
		Title:   "审核结果通知",
		Content: "您报名的活动【" + activity.Title + "】" + resultStr,
	}

	// 存数据库
	repository.CreateMessage(msg)
	// websocket 推送
	ws.SendMessageToUser(member.UserID, msg)
	return nil
}

// GetActivityMembers 获取某个活动的报名人员列表 (带越权校验)
func GetActivityMembers(initiatorID, activityID uint, status string, page, pageSize int) ([]dto.ActivityMemberRes, error) {
	// 1. 越权校验：不是发起人不能看别人的报名列表
	activity, err := repository.GetActivityByID(activityID)
	if err != nil {
		return nil, errors.New("活动不存在")
	}
	if activity.InitiatorID != initiatorID {
		return nil, errors.New("无权查看此活动的报名人员")
	}

	// 2. 底层查数据
	membersData, err := repository.GetMembersByActivityAndStatus(activityID, status, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 3. 组装 DTO
	var res []dto.ActivityMemberRes
	for _, m := range membersData {
		res = append(res, dto.ActivityMemberRes{
			MemberID:      m.ID,
			EquipmentNote: m.EquipmentNote,
			Status:        m.Status,
			User: dto.UserBasicInfo{
				ID:       m.User.ID,
				Nickname: m.User.Nickname,
				Avatar:   m.User.Avatar,
			},
			EmergencyContact: m.User.EmergencyContact,
		})
	}
	return res, nil
}

// SearchActivityByLocation 搜索活动 (复用 ActivityDetailRes 格式)
func SearchActivityByLocation(keyword string, page, pageSize int) ([]dto.ActivityDetailRes, error) {
	// 1. 调用 Repository 获取底层数据
	activities, err := repository.SearchActivitiesByLocation(keyword, page, pageSize)
	if err != nil {
		return nil, err
	}

	var res []dto.ActivityDetailRes

	// 如果没有搜到数据，返回空数组而不是 null (防呆设计)
	if len(activities) == 0 {
		return []dto.ActivityDetailRes{}, nil
	}

	// 2. 遍历数据，组装成安全的 DTO
	for _, a := range activities {
		// 解析图片 JSON 数组
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
			GroupQrCode:    "", // 🌟 列表页绝对不返回二维码
			CoverImage:     a.CoverImage,
			Description:    a.Description,
			Images:         images,
			MinFitness:     a.MinFitness,
			AgeLimit:       a.AgeLimit,
			MaxMembers:     a.MaxMembers,
			CurrentMembers: a.CurrentMembers,
			Status:         a.Status,
			WeatherAlert:   a.WeatherAlert,
			Initiator: dto.UserBasicInfo{
				ID:       a.Initiator.ID,
				Nickname: a.Initiator.Nickname,
				Avatar:   a.Initiator.Avatar,
			},
			HasApplied:  false, // 列表页不查个人的报名状态，省性能
			ApplyStatus: "",
		})
	}

	return res, nil
}
