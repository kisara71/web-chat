package utils

import (
	"web-chat/pkg/regexp"

	"github.com/bwmarrin/snowflake"
)

type Utils struct {
	SnowFlake *snowflake.Node
	Regexp    *regexp.Handler
}

func NewUtils() *Utils {
	snowFlake, err := snowflake.NewNode(0)
	if err != nil {
		panic(err)
	}
	return &Utils{
		Regexp:    regexp.NewHandler(),
		SnowFlake: snowFlake,
	}
}
