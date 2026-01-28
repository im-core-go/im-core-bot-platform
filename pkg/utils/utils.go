package utils

import (
	"github.com/im-core-go/im-core-bot-platform/pkg/code"
	"github.com/im-core-go/im-core-bot-platform/pkg/http"
	"github.com/im-core-go/im-core-bot-platform/pkg/regexp"
	"github.com/im-core-go/im-core-bot-platform/pkg/uuid"

	"github.com/bwmarrin/snowflake"
	"github.com/redis/go-redis/v9"
)

type Utils struct {
	SnowFlake      *snowflake.Node
	Regexp         *regexp.Handler
	RequestHandler *http.RequestHandler
	Code           *code.Manager
	UUID           *uuid.Wrap
}

func NewUtils(redisCmd redis.Cmdable) *Utils {
	snowFlake, err := snowflake.NewNode(0)
	if err != nil {
		panic(err)
	}
	return &Utils{
		Regexp:         regexp.NewHandler(),
		SnowFlake:      snowFlake,
		Code:           code.NewManager(redisCmd),
		RequestHandler: http.NewRequestHandler(),
		UUID:           uuid.NewWrap(),
	}
}
