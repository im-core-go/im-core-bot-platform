package utils

import (
	"web-chat/pkg/http"
	"web-chat/pkg/regexp"
	"web-chat/pkg/uuid"

	"github.com/bwmarrin/snowflake"
)

type Utils struct {
	SnowFlake      *snowflake.Node
	Regexp         *regexp.Handler
	RequestHandler *http.RequestHandler
<<<<<<< Updated upstream
=======
	Code           *code.Manager
	UUID           *uuid.Wrap
>>>>>>> Stashed changes
}

func NewUtils() *Utils {
	snowFlake, err := snowflake.NewNode(0)
	if err != nil {
		panic(err)
	}
	return &Utils{
		Regexp:         regexp.NewHandler(),
		SnowFlake:      snowFlake,
		RequestHandler: http.NewRequestHandler(),
		UUID:           uuid.NewWrap(),
	}
}
