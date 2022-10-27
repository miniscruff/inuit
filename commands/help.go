package commands

import (
	"fmt"
	"strings"
)

func helpCommand() *Command {
	return &Command{
		Key: "help",
		Help: func() string {
			return "print help on running a command"
		},
		Run: helpAction,
	}
}

func helpAction(editor Editor, args []string) (string, error) {
	if len(args) == 0 {
		var builder strings.Builder
		longestKey := 0
		for _, cmd := range editor.Commands().commands {
			lKey := len(cmd.Key)
			if lKey > longestKey {
				longestKey = lKey
			}
		}

		for _, cmd := range editor.Commands().commands {
			WriteFormat(
				&builder,
				"%v%v %v",
				cmd.Key,
				strings.Repeat(" ", longestKey-len(cmd.Key)),
				cmd.Help(),
			)
		}
		return builder.String(), nil
	}

	cmd, _ := FindCommand(nil, editor.Commands().commands, args)
	if cmd == nil {
		return "", fmt.Errorf("%w: %v", errCommandNotFound, strings.Join(args, " "))
	}

	return cmd.Help(), nil
}
