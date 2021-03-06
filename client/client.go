package client

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"file-transporter/common/constants"
	"file-transporter/common/utils"
	"fmt"
	"github.com/vence722/convert"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
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
		if err := handleFileReceiverConnection(fileReceiverConn); err != nil {
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
	if constants.ResponseOK != loginResp {
		return errors.New("Non-OK login response: " + loginResp)
	}

	return nil
}

func handleFileReceiverConnection(conn net.Conn) error {
	serverReader := bufio.NewReader(conn)
	serverWriter := bufio.NewWriter(conn)
	for {
		// Read filename
		fileName, err := serverReader.ReadString(constants.CommandDelimiter)
		if err != nil {
			return err
		}
		fileName = fileName[:len(fileName)-1]

		// Read file size
		fileSizeBuf := make([]byte, 8)
		_, err = io.ReadFull(serverReader, fileSizeBuf)
		if err != nil {
			return err
		}
		fileSize, _ := binary.ReadVarint(bytes.NewReader(fileSizeBuf))

		// Create file (create to the same directory to the executable)
		f, err := os.Create(fileName)
		if err != nil {
			return err
		}

		// Write OK response
		serverWriter.WriteString(constants.ResponseOK)
		serverWriter.WriteByte(constants.CommandDelimiter)
		serverWriter.Flush()

		fmt.Println("\n\n[INFO] Start receiving file", fileName)

		// Receive file data
		buf := make([]byte, constants.DefaultBufferSize)
		bytesReceived := int64(0)
		for {
			bytesRead, err := serverReader.Read(buf)
			if err != nil {
				return err
			}
			_, err = f.Write(buf[:bytesRead])
			if err != nil {
				return err
			}
			bytesReceived += int64(bytesRead)
			if bytesReceived >= fileSize {
				break
			}
		}
		fmt.Println("\n\n[INFO] File", fileName, "received successfully!")
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
	if constants.ResponseOK != loginResp {
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
		case "":
			continue
		default:
			fmt.Println("[ERROR] Invalid action", action)
		}
	}
}

func printMenu() {
	fmt.Println()
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

	return nil
}

func handleSendFile(serverReader *bufio.Reader, serverWriter *bufio.Writer) {
	// Read user input
	fmt.Println("Input target username:")
	targetUsername := utils.ReadCliInput()
	fmt.Println("Input file path:")
	filePath := utils.ReadCliInput()
	var fileSize int64
	if fi, err := os.Stat(filePath); err != nil {
		fmt.Println("[ERROR] File path is not valid, command will end")
		return
	} else if fi.IsDir() {
		fmt.Println("[ERROR] File path is a directory, and a file is expected, command will end")
		return
	} else {
		fileSize = fi.Size()
	}
	fileToSend, err := os.Open(filePath)
	if err != nil {
		fmt.Println("[ERROR] Failed to open file:", err.Error())
		return
	}

	// Send action
	serverWriter.WriteByte(constants.ActionSendFile)

	// Send target username
	serverWriter.WriteString(targetUsername)
	serverWriter.WriteByte(constants.CommandDelimiter)

	// Send filename
	serverWriter.WriteString(filepath.Base(fileToSend.Name()))
	serverWriter.WriteByte(constants.CommandDelimiter)

	// Send file size
	fileSizeBuf := make([]byte, 8)
	binary.PutVarint(fileSizeBuf, fileSize)
	serverWriter.Write(fileSizeBuf)
	serverWriter.Flush()

	// Read response
	resp, err := serverReader.ReadString(constants.CommandDelimiter)
	if err != nil {
		fmt.Println("[ERROR] Failed to send action:", err.Error())
		return
	}
	resp = resp[:len(resp)-1]
	if constants.ResponseOK != resp {
		fmt.Println("[ERROR] Non-OK login response:", resp)
		return
	}

	fmt.Println("[INFO] Start sending file", filepath.Base(fileToSend.Name()))

	// Send file
	buf := make([]byte, constants.DefaultBufferSize)
	bytesTransferred := int64(0)
	for {
		bytesRead, err := fileToSend.Read(buf)
		if err != nil {
			fmt.Println("[ERROR] Failed to read file:", err.Error())
			return
		}
		_, err = serverWriter.Write(buf[:bytesRead])
		if err != nil {
			fmt.Println("[ERROR] Failed to write file to server:", err.Error())
			return
		}
		bytesTransferred += int64(bytesRead)
		if bytesTransferred >= fileSize {
			break
		}
	}

	fmt.Println("[INFO] File", filepath.Base(fileToSend.Name()), "is sent successfully!")
}

func handleLogout(conn net.Conn, fileReceiverConn net.Conn, serverWriter *bufio.Writer) {
	serverWriter.WriteByte(constants.ActionLogout)
	serverWriter.Flush()
	conn.Close()
	fileReceiverConn.Close()
	fmt.Println("[INFO] Logged out")
}
