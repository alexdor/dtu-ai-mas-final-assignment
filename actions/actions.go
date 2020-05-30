package actions

import (
	"strings"

	"github.com/alexdor/dtu-ai-mas-final-assignment/communication"
	"github.com/alexdor/dtu-ai-mas-final-assignment/config"
)

type (
	Direction  = byte
	Action     = []byte
	PushOrPull func(agentDirection, boxDirecation Direction) Action
)

var (
	move       Action = []byte("Move")
	pull       Action = []byte("Pull")
	push       Action = []byte("Push")
	NoOpAction Action = []byte("NoOp")

	North Direction = 'N'
	West  Direction = 'W'
	South Direction = 'S'
	East  Direction = 'E'

	SingleAgentEnd = byte('\n')
	MultiAgentEnd  = byte(';')
)

func Move(direction Direction, endWith byte) Action {
	return append(move, '(', direction, ')', endWith)
}

func NoOp(endWith byte) Action {
	return append(NoOpAction, endWith)
}

type PullOrPush = func(agentDirection, boxDirecation Direction, endWith byte) Action

func Pull(agentDirection, boxDirecation Direction, endWith byte) Action {
	return append(pull, '(', agentDirection, ',', boxDirecation, ')', endWith)
}

func Push(agentDirection, boxDirecation Direction, endWith byte) Action {
	return append(push, '(', agentDirection, ',', boxDirecation, ')', endWith)
}

func ExecuteActions(actions Action) bool {
	var actionRow string
	for _, action := range strings.Split(string(actions), "\n") {
		if len(action) == 0 {
			continue
		}
		actionRow = strings.TrimRight(action, ";")
		res, err := communication.SendMessage(actionRow)
		if err != nil {
			communication.Log(actionRow)
			communication.Error(err)
			return false
		}

		for _, msg := range strings.Split(res, ";") {
			if msg != config.ServersTrueValue {
				communication.Log("\n")
				communication.Log("Server didn't accept the following action:")
				communication.Log(actionRow)
				communication.Log(res, "\n")
				return false
			}
		}
	}

	return true
}
