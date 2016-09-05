package server

import (
	"errors"
)

//异常
var (
	ERROR_SERVER_ALREADY_START = errors.New("server already start")
	ERROR_COMMND_NOT_FOUND     = errors.New("commnd not found")
	ERROR_COMMND_PARAM_ERROR   = errors.New("commnd param error")
	ERROR_COMMND_NO_LOGIN      = errors.New("you have no connect")
	ERROR_ITEM_NOT_EXIST       = errors.New("item not exist")
	ERROR_AUTHORITY_NO_PWD     = errors.New("system use password")
	ERROR_AUTHORITY_PWD_ERROR  = errors.New("password error")
	ERROR_PORT_ERROR           = errors.New("port error")
)

const (
	MESSAGE_SUCCESS = "SUCCESS"
	MESSAGE_ERROR   = "ERROR"
	MESSAGE_PONG    = "PONG"
	MESSAGE_EXIT    = "Bye"
)
