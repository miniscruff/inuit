package commands

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func FindCommand(base *Command, commands []*Command, split []string) (*Command, []string) {
	for _, com := range commands {
		if split[0] != com.Key {
			continue
		}

		if len(com.Subcommands) == 0 {
			return com, split[1:]
		}

		if len(split) > 1 {
			return FindCommand(com, com.Subcommands, split[1:])
		}

		return com, nil
	}

	return base, split
}

func ExistingScenes() ([]string, error) {
	var dirs []string

	entries, err := os.ReadDir(InternalDir)
	if err != nil {
		return dirs, fmt.Errorf("failure to read scenes: %w", err)
	}

	for _, dir := range entries {
		n := dir.Name()
		if n == AssetsFile || n == ContentsFile || n == MetadataFile {
			continue
		}

		dirs = append(dirs, strings.TrimSuffix(n, ".json"))
	}

	return dirs, nil
}

func SlicesEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func WriteFormat(w io.StringWriter, format string, args ...any) {
	w.WriteString(fmt.Sprintf(format+"\n", args...))
}

func Filter[T any](prefix string, search []T, keyFunc func(T) string) []string {
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

func MathOp(value float64, operator string, operand float64) float64 {
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

func RequiredArgs(count int) Validation {
	return func(editor Editor, args []string) error {
		if len(args) != count {
			return fmt.Errorf("%v != %v: %w", len(args), count, errIncorrectNumberOfArgs)
		}

		return nil
	}
}

func ArgsIn(index int, options []string) Validation {
	return func(editor Editor, args []string) error {
		if len(args) <= index {
			return errIncorrectNumberOfArgs
		}

		for _, o := range options {
			if args[index] == o {
				return nil
			}
		}

		return fmt.Errorf("%w: args[%v] not in %v", errInvalidArg, index, options)
	}
}

func ArgFloat(index int) Validation {
	return func(editor Editor, args []string) error {
		if len(args) <= index {
			return errIncorrectNumberOfArgs
		}

		_, err := strconv.ParseFloat(args[index], 64)
		if err != nil {
			return fmt.Errorf("%w: args[%v] not a float64", errInvalidArg, index)
		}

		return nil
	}
}

func RequiresVisual() Validation {
	return func(editor Editor, args []string) error {
		if editor.Visual() == nil {
			return errNoActiveVisual
		}

		return nil
	}
}

func StringUnchanged(value string) string {
	return value
}
