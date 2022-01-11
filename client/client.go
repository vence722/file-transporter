package client

import (
	"bufio"
	"errors"
	"file-transporter/common/constants"
	"file-transporter/common/utils"
	"fmt"
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
		case "1":
			if err := handleListOnlineUsers(serverReader, serverWriter); err != nil {
				fmt.Println("[ERROR] Failed to list online users:", err.Error())
			}
		case "2":
		case "3":
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
	serverWriter.WriteByte(1)
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

func handleLogout(conn net.Conn, fileReceiverConn net.Conn, serverWriter *bufio.Writer) {
	serverWriter.WriteByte(3)
	serverWriter.Flush()
	conn.Close()
	fileReceiverConn.Close()
	fmt.Println("[INFO] Logged out")
}
