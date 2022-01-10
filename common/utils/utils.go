package utils

import (
	"bufio"
	"os"
	"strings"
)

func ReadCliInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.ReplaceAll(text, "\n", "")
}
