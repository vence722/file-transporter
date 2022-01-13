package utils

import (
	"bufio"
	"os"
	"runtime"
	"strings"
)

func ReadCliInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.ReplaceAll(text, GetLineSeparator(), "")
}

func GetLineSeparator() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
