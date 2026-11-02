package config

import "errors"

var ErrBotTokenEmpty = errors.New("bot token is empty")
var ErrRedisAddrEmpty = errors.New("redis addr is empty")
