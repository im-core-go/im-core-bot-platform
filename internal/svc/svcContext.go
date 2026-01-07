package svc

import (
	"web-chat/configs"
	"web-chat/internal/dao"
	"web-chat/pkg/auth"
	"web-chat/pkg/infra"
	"web-chat/pkg/utils"
)

type Context struct {
	Config configs.Config
	Dao    *dao.Dao
	Utils  *utils.Utils
	Infra  *infra.Infra
	Auth   *auth.JwtHandler
}

func NewContext(cfg configs.Config) *Context {
	infraSvc := infra.NewInfra(cfg)
	return &Context{
		Config: cfg,
		Utils:  utils.NewUtils(infraSvc.Redis),
		Dao:    dao.NewDao(infraSvc.DB),
		Infra:  infraSvc,
		Auth:   auth.NewJwtHandler(),
	}
}
