package examples

import (
	"log"

	"github.com/unix-streamdeck/api/v3"
)

type CommandHandler struct{}

func (c *CommandHandler) Input(fields map[string]any, handlerType api.HandlerType, info api.StreamDeckInfo, event api.InputEvent, host api.IHost) {
	if !shouldAct(event) {
		return
	}

	command, ok := fieldString(fields, "command")
	if !ok {
		log.Println("[WARN] Command handler: no command specified")
		return
	}

	host.RunCommand(command)
}

func RegisterCommand() api.Module {
	return api.Module{
		Name: "Command",
		NewCombined: func() api.Handler {
			return &CommandHandler{}
		},
		CombinedFields: []api.Field{
			{Title: "Command", Name: "command", Type: api.Text},
		},
	}
}
