package user

import (
	"fmt"
	"web-chat/api/http_model"
)

func (l *logicImpl) GetUserInfo(userID int64) (*http_model.UserInfoResp, error) {
	if userID == 0 {
		return nil, fmt.Errorf("user id is required")
	}
	entity, err := l.svcCtx.Dao.UserDao.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return &http_model.UserInfoResp{
		ID:        entity.ID,
		NickName:  entity.NickName,
		Email:     entity.Email,
		Phone:     entity.Phone,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}, nil
}
