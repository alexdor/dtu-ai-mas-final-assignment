package communication

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alexdor/dtu-ai-mas-final-assignment/config"
)

var stdinReader *bufio.Scanner

func Init() {
	stdinReader = bufio.NewScanner(os.Stdin)

	SendMessage(config.Config.Name)
}

func Log(message ...interface{}) {
	fmt.Fprintln(os.Stderr, message...)
}
func Error(message ...error) {
	fmt.Fprintln(os.Stderr, "Error :", message)
}

func SendComment(message ...interface{}) {
	fmt.Println("#", message)
}

func SendMessage(message ...interface{}) (string, error) {
	fmt.Fprintln(os.Stdout, message...)
	return ReadNextMessages()
}

func ReadNextMessages() (string, error) {
	stdinReader.Scan()
	return stdinReader.Text(), stdinReader.Err()
}
