package constants

const (
	CommandDelimiter = '\n'
	StringSeparator  = ","

	ResponseOK                    = "OK"
	ResponseErrReadTargetUserName = "ERR_READ_USERNAME"
	ResponseErrReadFileName       = "ERR_READ_FILE_NAME"
	ResponseErrReadFileSize       = "ERR_READ_FILE_SIZE"
	ResponseErrUserNotExist       = "ERR_USER_NOT_EXIST"
	ResponseErrUserNotReady       = "ERR_USER_NOT_READY"
	ResponseErrUserReceiveFailed  = "ERR_USER_RECEIVE_FAILED"

	LoginTypeFileReceiver byte = 1
	LoginTypeCommandLine  byte = 2

	ActionListOnlineUsers byte = 1
	ActionSendFile        byte = 2
	ActionLogout          byte = 3

	DefaultBufferSize = 1 * 1024 * 1024
)
