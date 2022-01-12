package constants

const (
	CommandDelimiter = '\n'
	StringSeparator  = ","

	LoginTypeFileReceiver byte = 1
	LoginTypeCommandLine  byte = 2

	ActionListOnlineUsers byte = 1
	ActionSendFile        byte = 2
	ActionLogout          byte = 3
)
