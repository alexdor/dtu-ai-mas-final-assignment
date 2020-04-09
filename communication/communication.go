package communication

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alexdor/dtu-ai-mas-final-assigment/config"
)

var stdinReader *bufio.Reader

func Init() {
	stdinReader = bufio.NewReader(os.Stdin)
	SendMessage(config.Config.Name)
}

func Log(message ...interface{}) {
	fmt.Fprintln(os.Stderr, message...)
}
func Error(message ...error) {
	fmt.Fprintln(os.Stderr, "Error :", message)
}

func ReadNextMessages() (string, error) {
	return stdinReader.ReadString('\n')
}

func SendMessage(message ...interface{}) {
	fmt.Fprintln(os.Stdout, message...)
}
