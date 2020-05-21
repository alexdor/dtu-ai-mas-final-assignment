package actions

import (
	"strings"

	"github.com/alexdor/dtu-ai-mas-final-assignment/communication"
	"github.com/alexdor/dtu-ai-mas-final-assignment/config"
)

type (
	Direction = byte
	Action    = []byte

	PushOrPull func(agentDirection, boxDirecation Direction) Action
)

var (
	move Action = []byte("Move")
	pull Action = []byte("Pull")
	push Action = []byte("Push")
	NoOp Action = []byte("NoOp")

	North Direction = 'N'
	West  Direction = 'W'
	South Direction = 'S'
	East  Direction = 'E'
)

func Move(direction Direction) Action {
	return append(move, '(', direction, ')', '\n')
}

type PullOrPush = func(agentDirection, boxDirecation Direction) Action

func Pull(agentDirection, boxDirecation Direction) Action {
	return append(pull, '(', agentDirection, ',', boxDirecation, ')', '\n')
}

func Push(agentDirection, boxDirecation Direction) Action {
	return append(push, '(', agentDirection, ',', boxDirecation, ')', '\n')
}

func ExecuteActions(actions Action) bool {
	stringActions := string(actions)

	for _, action := range strings.Split(stringActions, ";") {
		action = strings.Split(action, "!")[0]
		communication.SendMessage(strings.TrimRight(action, ";"))

		res, err := communication.ReadNextMessages()
		if err != nil {
			communication.Error(err)
			return false
		}

		for _, msg := range strings.Split(res, ";") {
			if msg != config.ServersTrueValue {
				return false
			}
		}
	}

	return true
}
