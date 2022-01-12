package client

import (
	"bufio"
	"errors"
	"file-transporter/common/constants"
	"file-transporter/common/utils"
	"fmt"
	"github.com/vence722/convert"
	"net"
	"strings"
	"time"
)

func StartFileTransporterClient(serverHostPort string, username string) error {
	// Create connection for receiving files
	fileReceiverConn, err := net.Dial("tcp", serverHostPort)
	if err != nil {
		return err
	}

	// Login file receiver
	if err := loginFileReceiver(fileReceiverConn, username); err != nil {
		return err
	}
	// Handle file receiver connection
	go func() {
		if err := handleFileReceiverConnection(fileReceiverConn, username); err != nil {
			fmt.Println("[ERROR] Error from file receiver connection:", err.Error())
		}
	}()

	// Create connection for command-line
	cmdConn, err := net.Dial("tcp", serverHostPort)
	if err != nil {
		return err
	}
	return handleCommandLineConnection(cmdConn, fileReceiverConn, username)
}

func loginFileReceiver(conn net.Conn, username string) error {
	serverReader := bufio.NewReader(conn)
	serverWriter := bufio.NewWriter(conn)

	// Send login type
	serverWriter.WriteByte(constants.LoginTypeFileReceiver)

	// Send username
	serverWriter.WriteString(username)
	serverWriter.WriteByte(constants.CommandDelimiter)
	serverWriter.Flush()

	// Read server response
	loginResp, err := serverReader.ReadString(constants.CommandDelimiter)
	if err != nil {
		return err
	}
	loginResp = loginResp[:len(loginResp)-1]
	if "OK" != loginResp {
		return errors.New("Non-OK login response: " + loginResp)
	}

	return nil
}

func handleFileReceiverConnection(conn net.Conn, username string) error {
	for {
		time.Sleep(1 * time.Second)
	}
}

func handleCommandLineConnection(conn net.Conn, fileReceiverConn net.Conn, username string) error {
	serverReader := bufio.NewReader(conn)
	serverWriter := bufio.NewWriter(conn)

	// Send login type
	serverWriter.WriteByte(constants.LoginTypeCommandLine)

	// Send username
	serverWriter.WriteString(username)
	serverWriter.WriteByte(constants.CommandDelimiter)
	serverWriter.Flush()

	// Read server response
	loginResp, err := serverReader.ReadString(constants.CommandDelimiter)
	if err != nil {
		return err
	}
	loginResp = loginResp[:len(loginResp)-1]
	if "OK" != loginResp {
		return errors.New("Non-OK login response: " + loginResp)
	}

	// Handle user actions
	for {
		// Print Menu
		printMenu()
		action := utils.ReadCliInput()
		fmt.Println()
		switch action {
		case convert.Int2Str(constants.ActionListOnlineUsers):
			if err := handleListOnlineUsers(serverReader, serverWriter); err != nil {
				fmt.Println("[ERROR] Failed to list online users:", err.Error())
			}
		case convert.Int2Str(constants.ActionSendFile):
			handleSendFile(serverReader, serverWriter)
		case convert.Int2Str(constants.ActionLogout):
			handleLogout(conn, fileReceiverConn, serverWriter)
			return nil
		default:
			fmt.Println("[ERROR] Invalid action", action)
		}
	}
}

func printMenu() {
	fmt.Println("Choose Action:")
	fmt.Println("(1) List online users")
	fmt.Println("(2) Transfer file")
	fmt.Println("(3) Logout")
	fmt.Printf("Please input your action >>> ")
}

func handleListOnlineUsers(serverReader *bufio.Reader, serverWriter *bufio.Writer) error {
	// Send action
	serverWriter.WriteByte(constants.ActionListOnlineUsers)
	serverWriter.Flush()

	// Read response
	usersList, err := serverReader.ReadString(constants.CommandDelimiter)
	if err != nil {
		return err
	}
	usersList = usersList[:len(usersList)-1]

	// Print response
	fmt.Println("Online Users:")
	users := strings.Split(usersList, constants.StringSeparator)
	for _, user := range users {
		fmt.Println(user)
	}
	fmt.Println()

	return nil
}

func handleSendFile(serverReader *bufio.Reader, serverWriter *bufio.Writer) {
	//// Read user input
	//fmt.Println("Input target username:")
	//targetUsername := utils.ReadCliInput()
	//fmt.Println("Input file path:")
	//filePath := utils.ReadCliInput()
	//var fileSize int64
	//if fi, err := os.Stat(filePath); err != nil {
	//	fmt.Println("File path is not valid, command will end")
	//	return
	//} else if fi.IsDir() {
	//	fmt.Println("File path is a directory, and a file is expected, command will end")
	//	return
	//} else {
	//	fileSize = fi.Size()
	//}
	//fileToSend, err := os.Open(filePath)
	//if err != nil {
	//	fmt.Println("Failed to open file:", err.Error())
	//	return
	//}
	//
	//// Send action
	//serverWriter.WriteByte(constants.ActionSendFile)
	//
	//// Send filename
	//serverWriter.WriteString(fileToSend.Name())
	//serverWriter.WriteByte(constants.CommandDelimiter)
	//
	//// Send file size
	//fileSizeBuf := make([]byte, 8)
	//binary.PutVarint(fileSizeBuf, fileSize)
	//serverWriter.Write(fileSizeBuf)
	//serverWriter.Flush()
	//
	//// Read response
	//resp, err := serverReader.ReadString(constants.CommandDelimiter)
	//if err != nil {
	//	fmt.Println("Failed to send action:", err.Error())
	//	return
	//}
	//resp = resp[:len(resp)-1]
	//if "OK" != resp {
	//	fmt.Println("Non-OK login response:", resp)
	//	return
	//}
}

func handleLogout(conn net.Conn, fileReceiverConn net.Conn, serverWriter *bufio.Writer) {
	serverWriter.WriteByte(constants.ActionLogout)
	serverWriter.Flush()
	conn.Close()
	fileReceiverConn.Close()
	fmt.Println("[INFO] Logged out")
}
