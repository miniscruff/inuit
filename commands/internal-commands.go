package commands

/*
import (
	"errors"
	"strconv"
	"strings"
)


type Editor interface {
}

type Command struct {
	Key         string
	Help        func() string
	Suggestions func(editor Editor, partial []string) []string
	Run         func(editor Editor, args []string) (string, error)
	Subcommands []*Command
}

func runCommand(editor Editor, commands []*Command, text string) (string, error) {
	split := strings.Split(text, " ")
	cmd, args := findCommand(nil, commands, split)
	if cmd == nil {
		return "", errCommandNotFound
	}

	return cmd.Run(editor, args)
}

func findCommand(base *Command, commands []*Command, split []string) (*Command, []string) {
	for _, c := range commands {
		if split[0] != c.Key {
			continue
		}

		if len(c.Subcommands) == 0 {
			return c, split[1:]
		}

		if len(split) > 1 {
			return findCommand(c, c.Subcommands, split[1:])
		}

		return c, nil
	}

	return base, split
}

func buildSuggestions(editor Editor, commands []*Command, text string) []string {
	if text == "" {
		return nil
	}

	split := strings.Split(strings.Trim(text, " "), " ")
	cmd, partial := findCommand(nil, commands, split)

	var search []*Command
	if cmd == nil {
		search = commands
	} else {
		search = cmd.Subcommands
	}

	if len(search) > 0 {
		var lastWord string
		if len(partial) > 0 {
			lastWord = partial[len(partial)-1]
		}
		return filter(lastWord, search, func(cmd *Command) string {
			return cmd.Key
		})
	}

	if cmd != nil && cmd.Suggestions != nil {
		return cmd.Suggestions(scene, partial)
	}

	return nil
}

func buildCommands() []*Command {
	return []*Command{
		cdCommand(),
		{
			Key: "help",
			Help: func() string {
				return "print help on running a command"
			},
			Run: helpAction,
		},
		{
			Key: "ls",
			Help: func() string {
				return "list the current object and children names"
			},
			Run: lsAction,
		},
		setCommand(),
		{
			Key: "write",
			Help: func() string {
				return "save scene to disk"
			},
			Run: func(scene *EditorScene, args []string) (string, error) {
				if err := internal.SaveSceneData(&scene.sceneData, "testing.json"); err != nil {
					return "scene unable to save", err
				}
				return "scene saved", nil
			},
		},
	}
}

/*
:::template:::

func xCommand() *Command {
	return &Command{
		Key: "",
		Help: func() string {
			return ""
		},
		Suggestions: func(scene *EditorScene, partial []string) []string {
			return nil
		},
		Run: func(scene *EditorScene, args []string) (string, error) {
			return "", nil
		},
		Subcommands: []*Command{
		},
	}
}

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
