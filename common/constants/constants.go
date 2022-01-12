package constants

const (
	CommandDelimiter = '\n'
	StringSeparator  = ","

	ResponseOK                    = "OK"
	ResponseErrReadTargetUserName = "ERR_01"
	ResponseErrReadFileName       = "ERR_002"
	ResponseErrReadFileSize       = "ERR_003"

	LoginTypeFileReceiver byte = 1
	LoginTypeCommandLine  byte = 2

	ActionListOnlineUsers byte = 1
	ActionSendFile        byte = 2
	ActionLogout          byte = 3
)
