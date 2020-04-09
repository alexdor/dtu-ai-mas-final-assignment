package communication

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alexdor/dtu-ai-mas-final-assigment/config"
)

var stdin_reader *bufio.Reader

func Init() {
	stdin_reader = bufio.NewReader(os.Stdin)
	SendMessage(config.Config.Name)
}

func Log(message ...interface{}) {
	fmt.Fprintln(os.Stderr, message...)
}
func Error(message ...error) {
	fmt.Fprintln(os.Stderr, "Error :", message)
}

func ReadNextMessages() (string, error) {
	return stdin_reader.ReadString('\n')
}

func SendMessage(message ...interface{}) {
	fmt.Fprintln(os.Stdout, message...)
}
