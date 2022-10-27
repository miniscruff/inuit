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
		// setCommand(),
		writeCommand(),
	}
}

/*
func cdCommand() *Command {
	return &Command{
		Key: "cd",
		Help: func() string {
			return "change active object up or down"
		},
		Suggestions: func(scene *EditorScene, partial []string) []string {
			if len(partial) != 1 {
				return nil
			}

			var search []*internal.SceneVisual
			if scene.activeVisual == nil {
				search = scene.sceneData.Visuals
			} else {
				search = scene.activeVisual.Children
			}

			return filter(partial[0], search, func(child *internal.SceneVisual) string {
				return child.Name
			})
		},
		Run: func(scene *EditorScene, args []string) (string, error) {
			if len(args) == 0 {
				scene.activeVisual = nil
			} else if args[0] == ".." {
				if scene.activeVisual != nil {
					scene.activeVisual = scene.activeVisual.Parent
				}
			} else {
				var search []*internal.SceneVisual
				if scene.activeVisual == nil {
					search = scene.sceneData.Visuals
				} else {
					search = scene.activeVisual.Children
				}

				for _, child := range search {
					if child.Name == args[0] {
						scene.activeVisual = child
						break
					}
				}
			}
			return "", nil
		},
	}
}

	visible = true/false
	parent = name of our new parent
	sprite/font/text depending on type
	x = x pos
	y = y pos
	width
	height
	anchor <x y>
	anchor x value
	anchor y value
		etc

func setCommand() *Command {
	return &Command{
		Key: "set",
		Help: func() string {
			return "modify a value of our visual"
		},
		Subcommands: []*Command{
			{
				Key: "name",
				Help: func() string {
					return "change the name of our visual"
				},
				Suggestions: func(scene *EditorScene, partial []string) []string {
					return nil
				},
				Run: func(scene *EditorScene, args []string) (string, error) {
					if len(args) != 1 {
						return "", errors.New("need 1 arg to set name")
					}

					if scene.activeVisual == nil {
						return "", errors.New("need active to set name")
					}

					scene.activeVisual.Name = args[0]
					return "", nil
				},
			},
			{
				Key: "visible",
				Help: func() string {
					return "turn on or off our object"
				},
				Suggestions: func(scene *EditorScene, partial []string) []string {
					if len(partial) != 1 {
						return nil
					}

					return filter(partial[0], []string{"true", "false"}, func(a string) string { return a })
				},
				Run: func(scene *EditorScene, args []string) (string, error) {
					if len(args) != 1 {
						return "", errors.New("need either true or false to set visibility")
					}

					if scene.activeVisual == nil {
						return "", errors.New("need active to set visability")
					}

					if args[0] != "true" && args[0] != "false" {
						return "", errors.New("need either true or false to set visibility")
					}

					scene.activeVisual.Visual.SetVisible(args[0] == "true")
					return "", nil
				},
			},
			{
				Key: "x",
				Help: func() string {
					return "change x value"
				},
				Suggestions: func(scene *EditorScene, partial []string) []string {
					return nil
				},
				Run: func(scene *EditorScene, args []string) (string, error) {
					if scene.activeVisual == nil {
						return "", errors.New("need active to set visability")
					}

					if len(args) != 2 {
						return "", errors.New("need operator and operand eg \"+ 5\"")
					}

					operand, err := strconv.ParseFloat(args[1], 64)
					if err != nil {
						return "", errors.New("operand is not a float64")
					}

					scene.activeVisual.Visual.Transform.SetX(mathOp(
						scene.activeVisual.Visual.Transform.X(),
						args[0],
						operand,
					))
					scene.activeVisual.Transform.Position.X = scene.activeVisual.Visual.Transform.X()
					return "", nil
				},
			},
		},
	}
}
*/
