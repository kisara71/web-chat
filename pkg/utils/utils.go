package utils

import (
	"web-chat/pkg/http"
	"web-chat/pkg/regexp"

	"github.com/bwmarrin/snowflake"
)

type Utils struct {
	SnowFlake      *snowflake.Node
	Regexp         *regexp.Handler
	RequestHandler *http.RequestHandler
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
	}
}
