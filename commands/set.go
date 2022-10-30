package commands

import (
	"errors"
	"strconv"
)

var (
	errNoActiveVisual        = errors.New("need active visual")
	errIncorrectNumberOfArgs = errors.New("incorrect number of arguments")
	errInvalidArg            = errors.New("argument invalid")

	// some suggestion lists
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
			anchorCommand(),
			heightCommand(),
			nameCommand(),
			pivotCommand(),
			positionCommand(),
			visibleCommand(),
			widthCommand(),
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

func anchorCommand() *Command {
	return &Command{
		Key: "anchor",
		Help: func() string {
			return "set the anchors of our object"
		},
		Suggestions: func(editor Editor, partial []string) []string {
			switch len(partial) {
			case 1:
				return Filter(partial[0], sideKeys, StringUnchanged)
			case 2:
				return Filter(partial[1], ops, StringUnchanged)
			default:
				return nil
			}
		},
		Validations: []Validation{
			RequiresVisual(),
			RequiredArgs(3),
			ArgsIn(0, sideKeys),
			ArgsIn(1, ops),
			ArgFloat(2),
		},
		Run: func(editor Editor, args []string) (string, error) {
			operand, _ := strconv.ParseFloat(args[2], 64)
			var valueFunc func(value float64)
			var valuePtr *float64

			switch args[0] {
			case "left":
				valuePtr = &editor.Visual().Transform.Anchors.Left
				valueFunc = editor.Visual().Visual.Transform.SetLeftAnchor
			case "right":
				valuePtr = &editor.Visual().Transform.Anchors.Right
				valueFunc = editor.Visual().Visual.Transform.SetRightAnchor
			case "top":
				valuePtr = &editor.Visual().Transform.Anchors.Top
				valueFunc = editor.Visual().Visual.Transform.SetRightAnchor
			default:
				valuePtr = &editor.Visual().Transform.Anchors.Bottom
				valueFunc = editor.Visual().Visual.Transform.SetBottomAnchor
			}

			newValue := MathOp(
				*valuePtr,
				args[1],
				operand,
			)
			*valuePtr = newValue
			valueFunc(newValue)

			return "", nil
		},
	}
}

func heightCommand() *Command {
	return &Command{
		Key: "height",
		Help: func() string {
			return "set the height of our object"
		},
		Suggestions: func(editor Editor, partial []string) []string {
			switch len(partial) {
			case 1:
				return Filter(partial[0], ops, StringUnchanged)
			default:
				return nil
			}
		},
		Validations: []Validation{
			RequiresVisual(),
			RequiredArgs(2),
			ArgsIn(0, ops),
			ArgFloat(1),
		},
		Run: func(editor Editor, args []string) (string, error) {
			operand, _ := strconv.ParseFloat(args[1], 64)

			h := MathOp(
				editor.Visual().Visual.Transform.Height(),
				args[0],
				operand,
			)
			editor.Visual().Visual.Transform.SetHeight(h)
			editor.Visual().Transform.Height = h

			return "", nil
		},
	}
}

func offsetsCommand() *Command {
	return &Command{
		Key: "offsets",
		Help: func() string {
			return "set the offsets of our object"
		},
		Suggestions: func(editor Editor, partial []string) []string {
			switch len(partial) {
			case 1:
				return Filter(partial[0], sideKeys, StringUnchanged)
			case 2:
				return Filter(partial[1], ops, StringUnchanged)
			default:
				return nil
			}
		},
		Validations: []Validation{
			RequiresVisual(),
			RequiredArgs(3),
			ArgsIn(0, sideKeys),
			ArgsIn(1, ops),
			ArgFloat(2),
		},
		Run: func(editor Editor, args []string) (string, error) {
			operand, _ := strconv.ParseFloat(args[2], 64)
			var valueFunc func(value float64)
			var valuePtr *float64

			switch args[0] {
			case "left":
				valuePtr = &editor.Visual().Transform.Offsets.Left
				valueFunc = editor.Visual().Visual.Transform.SetLeftOffset
			case "right":
				valuePtr = &editor.Visual().Transform.Offsets.Right
				valueFunc = editor.Visual().Visual.Transform.SetRightOffset
			case "top":
				valuePtr = &editor.Visual().Transform.Offsets.Top
				valueFunc = editor.Visual().Visual.Transform.SetRightOffset
			default:
				valuePtr = &editor.Visual().Transform.Offsets.Bottom
				valueFunc = editor.Visual().Visual.Transform.SetBottomOffset
			}

			newValue := MathOp(
				*valuePtr,
				args[1],
				operand,
			)
			*valuePtr = newValue
			valueFunc(newValue)

			return "", nil
		},
	}
}

func pivotCommand() *Command {
	return &Command{
		Key: "pivot",
		Help: func() string {
			return "set the pivot of our object"
		},
		Suggestions: func(editor Editor, partial []string) []string {
			switch len(partial) {
			case 1:
				return Filter(partial[0], vec2Keys, StringUnchanged)
			case 2:
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
					editor.Visual().Visual.Transform.Pivot().X,
					args[1],
					operand,
				)
				editor.Visual().Visual.Transform.SetPivotX(x)
				editor.Visual().Transform.Pivot.X = x
			case "y":
				y := MathOp(
					editor.Visual().Visual.Transform.Pivot().Y,
					args[1],
					operand,
				)
				editor.Visual().Visual.Transform.SetPivotY(y)
				editor.Visual().Transform.Pivot.Y = y
			}

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
			case 1:
				return Filter(partial[0], vec2Keys, StringUnchanged)
			case 2:
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

func widthCommand() *Command {
	return &Command{
		Key: "width",
		Help: func() string {
			return "set the width of our object"
		},
		Suggestions: func(editor Editor, partial []string) []string {
			switch len(partial) {
			case 1:
				return Filter(partial[0], ops, StringUnchanged)
			default:
				return nil
			}
		},
		Validations: []Validation{
			RequiresVisual(),
			RequiredArgs(2),
			ArgsIn(0, ops),
			ArgFloat(1),
		},
		Run: func(editor Editor, args []string) (string, error) {
			operand, _ := strconv.ParseFloat(args[1], 64)

			w := MathOp(
				editor.Visual().Visual.Transform.Width(),
				args[0],
				operand,
			)
			editor.Visual().Visual.Transform.SetWidth(w)
			editor.Visual().Transform.Width = w

			return "", nil
		},
	}
}
