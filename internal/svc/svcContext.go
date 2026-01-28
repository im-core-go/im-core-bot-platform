package svc

import (
	"github.com/im-core-go/im-core-bot-platform/configs"
	"github.com/im-core-go/im-core-bot-platform/internal/dao"
	"github.com/im-core-go/im-core-bot-platform/pkg/auth"
	"github.com/im-core-go/im-core-bot-platform/pkg/infra"
	"github.com/im-core-go/im-core-bot-platform/pkg/utils"
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
