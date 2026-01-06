package utils

import (
	"web-chat/pkg/code"
	"web-chat/pkg/http"
	"web-chat/pkg/regexp"

	"github.com/bwmarrin/snowflake"
	"github.com/redis/go-redis/v9"
)

type Utils struct {
	SnowFlake      *snowflake.Node
	Regexp         *regexp.Handler
	RequestHandler *http.RequestHandler
	Code           *code.Manager
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
	}
}
