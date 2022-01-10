package main

import (
	"file-transporter/client"
	"file-transporter/server"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		showUsage()
		os.Exit(1)
	}
	mode := os.Args[1]
	if mode == "server" {
		if len(os.Args) < 3 {
			showUsage()
			os.Exit(1)
		}
		hostPort := os.Args[2]
		if err := server.StartFileTransporterServer(hostPort); err != nil {
			fmt.Println("[ERROR] Failed to start server:", err.Error())
			os.Exit(1)
		}
	} else if mode == "client" {
		if len(os.Args) < 4 {
			showUsage()
			os.Exit(1)
		}
		serverHostPort := os.Args[2]
		username := os.Args[3]
		if err := client.StartFileTransporterClient(serverHostPort, username); err != nil {
			fmt.Println("[ERROR] Client received error:", err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Println("[ERROR] mode can only be 'server' or 'client'")
		os.Exit(1)
	}
}

func showUsage() {
	fmt.Println("Usage: ./file-transporter <mode> <...args>")
	fmt.Println("  <mode> can be 'server' or 'client'")
	fmt.Println("  For 'server' mode: ./file-transporter server <host:port>")
	fmt.Println("  For 'client' mode: ./file-transporter client <server_host:server_port> <username>")
}
