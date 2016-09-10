package server

import (
	"errors"
)

//异常
var (
	ERROR_SERVER_ALREADY_START = errors.New("server already start")
	ERROR_SERVER_CONNECT_TYPE  = errors.New("connect type is not support")
	ERROR_COMMND_NOT_FOUND     = errors.New("commnd not found")
	ERROR_COMMND_PARAM_ERROR   = errors.New("commnd param error")
	ERROR_COMMND_NO_LOGIN      = errors.New("you have no connect")
	ERROR_ITEM_NOT_EXIST       = errors.New("item not exist")
	ERROR_TABLE_NOT_EXIST      = errors.New("table not exist")
	ERROR_AUTHORITY_NO_PWD     = errors.New("system use password")
	ERROR_AUTHORITY_PWD_ERROR  = errors.New("password error")
	ERROR_PORT_ERROR           = errors.New("port error")
	ERROR_PROTOCOL_ERROR       = errors.New("protocol error")
	ERROR_SYSTEM               = errors.New("system error")
)

const (
	MESSAGE_SUCCESS            = "SUCCESS"
	MESSAGE_ERROR              = "ERROR"
	MESSAGE_PONG               = "PONG"
	MESSAGE_EXIT               = "Bye"
	MESSAGE_NO_PWD             = "NO_PWD"
	MESSAGE_PWD_ERROR          = "PWD_ERROR"
	MESSAGE_PORT_ERROR         = "PORT_ERROR"
	MESSAGE_PROTOCOL_ERROR     = "PROTOCOL_ERROR"
	MESSAGE_COMMND_PARAM_ERROR = "COMMND_PARAM_ERROR"
	MESSAGE_ITEM_NOT_EXIST     = "ITEM_NOT_EXIST"
	MESSAGE_TABLE_NOT_EXIST    = "TABLE_NOT_EXIST"
	MESSAGE_COMMND_NOT_FOUND   = "COMMND_NOT_FOUND"
	MESSAGE_COMMND_NO_LOGIN    = "COMMND_NO_LOGIN"
)
