package server

import (
	"errors"
)

//异常
var (
	ERROR_SERVER_ALREADY_START     = errors.New("server already start")
	ERROR_SERVER_CONNECT_TYPE      = errors.New("connect type is not support")
	ERROR_COMMAND_NOT_FOUND        = errors.New("command not found")
	ERROR_COMMAND_PARAM_ERROR      = errors.New("command param error")
	ERROR_COMMAND_NO_LOGIN         = errors.New("you have no connect")
	ERROR_ITEM_NOT_EXIST           = errors.New("item not exist")
	ERROR_TABLE_NOT_EXIST          = errors.New("table not exist")
	ERROR_AUTHORITY_NO_PWD         = errors.New("system use password")
	ERROR_AUTHORITY_PWD_ERROR      = errors.New("password error")
	ERROR_PORT_ERROR               = errors.New("port error")
	ERROR_PROTOCOL_ERROR           = errors.New("protocol error")
	ERROR_SYSTEM                   = errors.New("system error")
	ERROR_COMMAND_NOT_SUPPORT_DATA = errors.New("command not support for the data type")
)

const (
	MESSAGE_SUCCESS                  = "SUCCESS"
	MESSAGE_ERROR                    = "ERROR"
	MESSAGE_PONG                     = "PONG"
	MESSAGE_EXIT                     = "Bye"
	MESSAGE_NO_PWD                   = "NO_PWD"
	MESSAGE_PWD_ERROR                = "PWD_ERROR"
	MESSAGE_PORT_ERROR               = "PORT_ERROR"
	MESSAGE_PROTOCOL_ERROR           = "PROTOCOL_ERROR"
	MESSAGE_COMMAND_PARAM_ERROR      = "COMMAND_PARAM_ERROR"
	MESSAGE_ITEM_NOT_EXIST           = "ITEM_NOT_EXIST"
	MESSAGE_TABLE_NOT_EXIST          = "TABLE_NOT_EXIST"
	MESSAGE_COMMAND_NOT_FOUND        = "COMMAND_NOT_FOUND"
	MESSAGE_COMMAND_NO_LOGIN         = "COMMAND_NO_LOGIN"
	MESSAGE_COMMAND_NOT_SUPPORT_DATA = "COMMAND_NOT_SUPPORT_DATA"
)
