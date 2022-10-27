package commands

import (
	"errors"
	"strconv"
)

var (
	errNoActiveVisual        = errors.New("need active visual")
	errIncorrectNumberOfArgs = errors.New("incorrect number of arguments")
	errInvalidArg            = errors.New("argument invalid")

	trueFalseOptions = []string{"true", "false"}
	sideKeys         = []string{"left", "right", "top", "bottom"}
	vec2Keys         = []string{"x", "y"}
	ops              = []string{"+", "-", "*", "/", "="}
)

func setCommand() *Command {
	return &Command{
		Key: "set",
		Help: func() string {
			return "modify a value of our visual"
		},
		Subcommands: []*Command{
			nameCommand(),
			visibleCommand(),
			/*
				{
					Key: "x",
					Help: func() string {
						return "change x value"
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
			*/
		},
	}
}

func nameCommand() *Command {
	return &Command{
		Key: "name",
		Help: func() string {
			return "change the name of our visual"
		},
		Validations: []Validation{
			RequiredArgs(1),
			RequiresVisual(),
		},
		Run: func(editor Editor, args []string) (string, error) {
			editor.Visual().Name = args[0]
			return "", nil
		},
	}
}

func visibleCommand() *Command {
	return &Command{
		Key: "visible",
		Help: func() string {
			return "turn on or off our object"
		},
		Suggestions: func(editor Editor, partial []string) []string {
			if len(partial) != 1 {
				return nil
			}

			return Filter(partial[0], trueFalseOptions, StringUnchanged)
		},
		Validations: []Validation{
			RequiresVisual(),
			RequiredArgs(1),
			ArgsIn(0, trueFalseOptions),
		},
		Run: func(editor Editor, args []string) (string, error) {
			editor.Visual().Visual.SetVisible(args[0] == "true")
			return "", nil
		},
	}
}

func positionCommand() *Command {
	return &Command{
		Key: "position",
		Help: func() string {
			return "set the position of our object"
		},
		Suggestions: func(editor Editor, partial []string) []string {
			switch len(partial) {
			case 0:
				return Filter(partial[0], vec2Keys, StringUnchanged)
			case 1:
				return Filter(partial[1], ops, StringUnchanged)
			default:
				return nil
			}
		},
		Validations: []Validation{
			RequiresVisual(),
			RequiredArgs(3),
			ArgsIn(0, vec2Keys),
			ArgsIn(1, ops),
			ArgFloat(2),
		},
		Run: func(editor Editor, args []string) (string, error) {
			operand, _ := strconv.ParseFloat(args[2], 64)

			switch args[0] {
			case "x":
				x := MathOp(
					editor.Visual().Visual.Transform.X(),
					args[1],
					operand,
				)
				editor.Visual().Visual.Transform.SetX(x)
				editor.Visual().Transform.Position.X = x
			case "y":
				y := MathOp(
					editor.Visual().Visual.Transform.Y(),
					args[1],
					operand,
				)
				editor.Visual().Visual.Transform.SetY(y)
				editor.Visual().Transform.Position.Y = y
			}

			return "", nil
		},
	}
}
