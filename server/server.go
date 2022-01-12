package server

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"file-transporter/common/constants"
	"fmt"
	"io"
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
	clientWriter.WriteString(constants.ResponseOK)
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
	clientWriter.WriteString(constants.ResponseOK)
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
			handleSendFile(clientReader, clientWriter, conn)
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

func handleSendFile(clientReader *bufio.Reader, clientWriter *bufio.Writer, conn net.Conn) {
	// Read target usernam
	targetUsername, err := clientReader.ReadString(constants.CommandDelimiter)
	if err != nil {
		fmt.Println("[ERROR] Failed to read target username:", err.Error())
		clientWriter.WriteString(constants.ResponseErrReadTargetUserName)
		clientWriter.WriteByte(constants.CommandDelimiter)
		clientWriter.Flush()
		return
	}
	targetUsername = targetUsername[:len(targetUsername)-1]

	// Read filename
	fileName, err := clientReader.ReadString(constants.CommandDelimiter)
	if err != nil {
		fmt.Println("[ERROR] Failed to read file name:", err.Error())
		clientWriter.WriteString(constants.ResponseErrReadFileName)
		clientWriter.WriteByte(constants.CommandDelimiter)
		clientWriter.Flush()
		return
	}
	fileName = fileName[:len(fileName)-1]

	// Read file size
	fileSizeBuf := make([]byte, 8)
	_, err = io.ReadFull(clientReader, fileSizeBuf)
	if err != nil {
		fmt.Println("[ERROR] Failed to read file size:", err.Error())
		clientWriter.WriteString(constants.ResponseErrReadFileSize)
		clientWriter.WriteByte(constants.CommandDelimiter)
		clientWriter.Flush()
		return
	}
	fileSize, _ := binary.ReadVarint(bytes.NewReader(fileSizeBuf))

	fmt.Println("[INFO] targetUserName:", targetUsername, "fileName:", fileName, "fileSize:", fileSize)

	// Check target user exists
	targetConn, ok := userStore.GetUser(targetUsername)
	if !ok {
		fmt.Println("[ERROR] Target user", targetUsername, "doesn't exist")
		clientWriter.WriteString(constants.ResponseErrUserNotExist)
		clientWriter.WriteByte(constants.CommandDelimiter)
		clientWriter.Flush()
		return
	}

	// Send file info to target user
	if !checkTargetReady(targetConn, fileName, fileSize) {
		fmt.Println("[ERROR] Target user", targetUsername, "not ready")
		clientWriter.WriteString(constants.ResponseErrUserNotReady)
		clientWriter.WriteByte(constants.CommandDelimiter)
		clientWriter.Flush()
		return
	}

	// Send OK response to client
	clientWriter.WriteString(constants.ResponseOK)
	clientWriter.WriteByte(constants.CommandDelimiter)
	clientWriter.Flush()

	// Transfer file data
	if err := transferFileData(conn, targetConn, fileSize); err != nil {
		fmt.Println("[ERROR] Target user", targetUsername, "receiving failed")
		clientWriter.WriteString(constants.ResponseErrUserReceiveFailed)
		clientWriter.WriteByte(constants.CommandDelimiter)
		clientWriter.Flush()
		return
	}

	clientWriter.WriteString(constants.ResponseOK)
	clientWriter.WriteByte(constants.CommandDelimiter)
	clientWriter.Flush()
}

func checkTargetReady(conn net.Conn, fileName string, fileSize int64) bool {
	clientReader := bufio.NewReader(conn)
	clientWriter := bufio.NewWriter(conn)

	// Write file name
	clientWriter.WriteString(fileName)
	clientWriter.WriteByte(constants.CommandDelimiter)

	// Write file size
	fileSizeBuf := make([]byte, 8)
	binary.PutVarint(fileSizeBuf, fileSize)
	clientWriter.Write(fileSizeBuf)
	clientWriter.Flush()

	// Read response
	resp, err := clientReader.ReadString(constants.CommandDelimiter)
	if err != nil {
		fmt.Println("[ERROR] Failed to send action:", err.Error())
		return false
	}
	resp = resp[:len(resp)-1]
	if constants.ResponseOK != resp {
		fmt.Println("Non-OK login response:", resp)
		return false
	}

	return true
}

func transferFileData(fromConn net.Conn, toConn net.Conn, fileSize int64) error {
	fromReader := bufio.NewReader(fromConn)
	toWriter := bufio.NewWriter(toConn)

	buf := make([]byte, constants.DefaultBufferSize)
	byteTransferred := int64(0)
	for {
		bytesRead, err := fromReader.Read(buf)
		if err != nil {
			return err
		}
		_, err = toWriter.Write(buf[:bytesRead])
		if err != nil {
			return err
		}
		byteTransferred += int64(bytesRead)
		if byteTransferred >= fileSize {
			break
		}
	}

	return nil
}
