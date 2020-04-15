package actions

import (
	"strings"

	"github.com/alexdor/dtu-ai-mas-final-assignment/communication"
	"github.com/alexdor/dtu-ai-mas-final-assignment/config"
)

type Direction = string
type Action = string

const (
	move Action = "Move"
	pull Action = "Pull"
	push Action = "Push"
	NoOp Action = "NoOp"
)

const (
	North Direction = "N"
	West  Direction = "W"
	South Direction = "S"
	East  Direction = "E"
)

func Move(direction Direction) Action {
	return move + "(" + direction + ")"
}

func Pull(agentDirection, boxDirecation Direction) Action {
	return pull + "(" + agentDirection + "," + boxDirecation + ")"
}

func Push(agentDirection, boxDirecation Direction) Action {
	return push + "(" + agentDirection + "," + boxDirecation + ")"
}

func ExecuteActions(actions ...Action) bool {
	var b strings.Builder

	for _, action := range actions {
		if _, err := b.WriteString(action); err != nil {
			communication.Error(err)
			return false
		}

		if _, err := b.WriteRune(';'); err != nil {
			communication.Error(err)
			return false
		}
	}

	communication.SendMessage(strings.TrimRight(b.String(), ";"))

	res, err := communication.ReadNextMessages()

	return err == nil && res == config.ServersTrueValue
}
