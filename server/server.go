package server

import (
	"bufio"
	"file-transporter/common/constants"
	"fmt"
	"net"
	"strings"
)

var (
	userStore = NewUserStore()
)

func StartFileTransporterServer(hostPort string) error {
	l, err := net.Listen("tcp", hostPort)
	if err != nil {
		return err
	}
	defer l.Close()
	fmt.Println("[INFO] File transporter server is started on", hostPort)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("[ERROR] Failed to accept connection, error:", err.Error())
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	clientReader := bufio.NewReader(conn)
	clientWriter := bufio.NewWriter(conn)

	// Read login type
	loginType, err := clientReader.ReadByte()
	if err != nil {
		fmt.Println("[ERROR] Client login failed, disconnecting")
		conn.Close()
		return
	}

	if loginType == constants.LoginTypeFileReceiver {
		handleFileReceiverConnection(clientReader, clientWriter, conn)
	} else if loginType == constants.LoginTypeCommandLine {
		handleCommandLineConnection(clientReader, clientWriter, conn)
	} else {
		fmt.Println("[ERROR] Login type", loginType, "is invalid, disconnecting")
		conn.Close()
		return
	}

}

func handleFileReceiverConnection(clientReader *bufio.Reader, clientWriter *bufio.Writer, conn net.Conn) {
	// Read username
	username, err := clientReader.ReadString(constants.CommandDelimiter)
	if err != nil {
		fmt.Println("[ERROR] Client login failed, disconnecting")
		conn.Close()
		return
	}
	// Trim delimiter
	username = username[:len(username)-1]

	// Store username
	if !userStore.AddUser(username, conn) {
		fmt.Println("[ERROR] Username already taken, disconnecting")
		clientWriter.WriteString("username already taken")
		clientWriter.WriteByte(constants.CommandDelimiter)
		conn.Close()
		return
	} else {
		fmt.Println("[INFO] User", username, "logged in")
	}

	// Send login success reply
	clientWriter.WriteString("OK")
	clientWriter.WriteByte(constants.CommandDelimiter)
	clientWriter.Flush()
}

func handleCommandLineConnection(clientReader *bufio.Reader, clientWriter *bufio.Writer, conn net.Conn) {
	// Read username
	username, err := clientReader.ReadString(constants.CommandDelimiter)
	if err != nil {
		fmt.Println("[ERROR] Client login failed, disconnecting")
		conn.Close()
		return
	}
	// Trim delimiter
	username = username[:len(username)-1]

	// Check login status
	if _, ok := userStore.GetUser(username); !ok {
		fmt.Println("[ERROR] Client login failed, disconnecting")
		conn.Close()
		return
	}

	// Send login success reply
	clientWriter.WriteString("OK")
	clientWriter.WriteByte(constants.CommandDelimiter)
	clientWriter.Flush()

	// Handle client actions
	for {
		action, err := clientReader.ReadByte()
		if err != nil {
			fmt.Println("[ERROR] Failed to read client action, disconnecting")
			conn.Close()
			userStore.RemoveUser(username)
			return
		}
		switch action {
		case constants.ActionListOnlineUsers:
			fmt.Println("[INFO] User", username, "requires online user list")
			usersList := userStore.ListUsers()
			clientWriter.WriteString(strings.Join(usersList, constants.StringSeparator))
			clientWriter.WriteByte(constants.CommandDelimiter)
			clientWriter.Flush()
		case constants.ActionSendFile:
			handleSendFile()
		case constants.ActionLogout:
			fmt.Println("[INFO] User", username, "logged out")
			conn.Close()
			userStore.RemoveUser(username)
			return
		default:
			fmt.Println("[ERROR] Client action", action, "invalid, disconnecting")
			conn.Close()
			userStore.RemoveUser(username)
			return
		}
	}
}

func handleSendFile() {

}
