package commands

import (
	"errors"
	"strings"
)

var (
	errCommandNotFound = errors.New("command not found")
)

type CommandAction func(editor Editor, args []string) (string, error)
type Validation func(editor Editor, args []string) error

type Editor interface {
	Visual() *SceneVisual
	SetVisual(visual *SceneVisual)
	SceneData() *SceneData
	Path() string
	Commands() *Commands
}

type Command struct {
	Key         string
	Help        func() string
	Suggestions func(editor Editor, partial []string) []string
	Validations []Validation
	Run         CommandAction
	Subcommands []*Command
}

type Commands struct {
	editor   Editor
	commands []*Command
}

func NewCommands(editor Editor) *Commands {
	return &Commands{
		editor:   editor,
		commands: buildCommands(),
	}
}

func (c *Commands) Run(text string) (string, error) {
	split := strings.Split(text, " ")
	cmd, args := FindCommand(nil, c.commands, split)
	if cmd == nil {
		return "", errCommandNotFound
	}

	for _, v := range cmd.Validations {
		if err := v(c.editor, args); err != nil {
			return "", err
		}
	}

	return cmd.Run(c.editor, args)
}

func (c *Commands) BuildSuggestions(text string) []string {
	if text == "" {
		return nil
	}

	split := strings.Split(strings.Trim(text, " "), " ")
	cmd, partial := FindCommand(nil, c.commands, split)

	var search []*Command
	if cmd == nil {
		search = c.commands
	} else {
		search = cmd.Subcommands
	}

	if len(search) > 0 {
		var lastWord string
		if len(partial) > 0 {
			lastWord = partial[len(partial)-1]
		}
		return Filter(lastWord, search, func(cmd *Command) string {
			return cmd.Key
		})
	}

	if cmd != nil && cmd.Suggestions != nil {
		return cmd.Suggestions(c.editor, partial)
	}

	return nil
}

func buildCommands() []*Command {
	return []*Command{
		cdCommand(),
		helpCommand(),
		lsCommand(),
		setCommand(),
		writeCommand(),
	}
}
