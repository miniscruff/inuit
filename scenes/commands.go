package scenes

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/miniscruff/inuit/internal"
)

/*
in "normal" mode:
WASD = move camera for current view
R = reset position
Space + mouse drag = move camera
Enter = eneter a command

commands:
help = get help
	help <cmd ...> help for a command
open = open a new scene
write = save scene to disk
reload = reload scene from disk undoing changes
history = output last 50 messages
	history <contains> = some sort of search would be good
cd = change directory, moves up and down the scene tree
	cd = go to root
	cd .. = go up
	cd <child> = go into child
ls = view current scene tree item and children
pwd = print working directory, our current object and the parents
hide/show = turn on or off for editor only
get = same as set but will just output the value
add <type> <name> = add a new visual under our current object

for anything that is a number you can use =, *=, +=, -=, /=, %=, ++, --
the "dot" command of "." will rerun the last command
a command stack ( say 50? ) will be saved and you can use the arrow keys to view them
*/

type Command struct {
	Key         string
	Help        func() string
	Suggestions func(scene *EditorScene, partial []string) []string
	Run         func(scene *EditorScene, args []string) (string, error)
	Subcommands []*Command
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

func runCommand(scene *EditorScene, commands []*Command, text string) (string, error) {
	split := strings.Split(text, " ")
	cmd, args := findCommand(nil, commands, split)
	if cmd == nil {
		return "", errors.New("command not found")
	}

	return cmd.Run(scene, args)
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

func buildSuggestions(scene *EditorScene, commands []*Command, text string) []string {
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

	if cmd != nil {
		return cmd.Suggestions(scene, partial)
	}

	return nil
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
*/

func helpCommand() *Command {
	return &Command{
		Key: "help",
		Help: func() string {
			return "print help on running a command"
		},
		Suggestions: func(scene *EditorScene, partial []string) []string {
			return nil
		},
		Run: func(scene *EditorScene, args []string) (string, error) {
			var builder strings.Builder
			if len(args) == 0 {
				longestKey := 0
				for _, cmd := range scene.commands {
					lKey := len(cmd.Key)
					if lKey > longestKey {
						longestKey = lKey
					}
				}

				for _, cmd := range scene.commands {
					writeFormat(
						&builder,
						"%v%v %v",
						cmd.Key,
						strings.Repeat(" ", longestKey-len(cmd.Key)),
						cmd.Help(),
					)
				}
				return builder.String(), nil
			} else {
				cmd, _ := findCommand(nil, scene.commands, args)
				if cmd == nil {
					return "", errors.New("command not found:" + strings.Join(args, " "))
				}
				return cmd.Help(), nil
			}
		},
	}
}

func writeCommand() *Command {
	return &Command{
		Key: "write",
		Help: func() string {
			return "save scene to disk"
		},
		Suggestions: func(scene *EditorScene, partial []string) []string {
			return nil
		},
		Run: func(scene *EditorScene, args []string) (string, error) {
			if err := internal.SaveSceneData(&scene.sceneData, "testing.json"); err != nil {
				return "scene unable to save", err
			}
			return "scene saved", nil
		},
	}
}

func lsCommand() *Command {
	return &Command{
		Key: "ls",
		Help: func() string {
			return "list the current object and children names"
		},
		Suggestions: func(scene *EditorScene, partial []string) []string {
			return nil
		},
		Run: func(scene *EditorScene, args []string) (string, error) {
			var builder strings.Builder

			if scene.activeVisual != nil {
				writeFormat(&builder, ".%v", scene.activeVisual.Name)
				for _, child := range scene.activeVisual.Children {
					writeFormat(
						&builder,
						"./%v",
						child.Name,
					)
				}
			} else {
				writeFormat(&builder, "No selection")
				for _, child := range scene.sceneData.Visuals {
					writeFormat(
						&builder,
						"./%v",
						child.Name,
					)
				}
			}
			return builder.String(), nil
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

/*
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
*/

func setCommand() *Command {
	return &Command{
		Key: "set",
		Help: func() string {
			return "modify a value of our visual"
		},
		Suggestions: func(scene *EditorScene, partial []string) []string {
			return nil
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
func writeFormat(w io.StringWriter, format string, args ...any) {
	w.WriteString(fmt.Sprintf(format+"\n", args...))
}

func filter[T any](prefix string, search []T, keyFunc func(T) string) []string {
	var res []string
	prefixLower := strings.ToLower(prefix)
	for _, item := range search {
		key := keyFunc(item)
		if strings.HasPrefix(strings.ToLower(key), prefixLower) {
			res = append(res, key)
		}
	}
	return res
}

func mathOp(value float64, operator string, operand float64) float64 {
	switch operator {
	case "=":
		return operand
	case "*":
		return value * operand
	case "+":
		return value + operand
	case "-":
		return value - operand
	case "/":
		return value / operand
	default:
		panic("unknown operator:" + operator)
	}
}
