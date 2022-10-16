package scenes

import (
	"errors"
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
? = help shortcut eg "?" = all commands, "? set" = help for set
open = open a new scene
help = same as ?
write = save scene to disk
reload = reload scene from disk undoing changes
history = output last 50 messages
	history <contains> = some sort of search would be good
cd = change directory, moves up and down the scene tree
	cd = go to root
	cd .. = go up
	cd <child> = go into child
ls = view current scene tree item and children
hide/show = turn on or off for editor only
set = set some value
	name = our objects name
	visible = true/false
	windowSize = true/false
	parent = name of our new parent
	sprite/font/text depending on type
	transform
		x = x pos
		y = y pos
		width
		height
		anchor <x y>
		anchor x value
		anchor y yvalue
		etc
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
		writeCommand(),
		helpCommand(),
	}
}

func runCommand(scene *EditorScene, commands []*Command, text string) (string, error) {
	split := strings.Split(text, " ")
	cmd, args := findCommand(commands, split)
	if cmd == nil {
		return "", errors.New("command not found")
	}

	return cmd.Run(scene, args)
}

func findCommand(commands []*Command, split []string) (*Command, []string) {
	for _, c := range commands {
		if split[0] == c.Key {
			if len(c.Subcommands) == 0 {
				return c, split[1:]
			} else {
				return findCommand(c.Subcommands, split[1:])
			}
		}
	}
	return nil, split
}

func buildSuggestions(scene *EditorScene, commands []*Command, text string) []string {
	if len(text) < 2 {
		return nil
	}

	// TODO: partial has to be all strings after the last subcommand
	// as its possible to have many levels of arguments with suggestions
	// such as for help or math: eg set ...x [operation] [const]
	split := strings.Split(text, " ")

	var search []*Command
	cmd, partial := findCommand(commands, split)

	if cmd == nil {
		search = commands
	} else {
		search = cmd.Subcommands
	}

	if len(search) > 0 {
		var filtered []string
		lastWord := partial[len(partial)-1]
		for _, sub := range search {
			if strings.HasPrefix(strings.ToLower(sub.Key), lastWord) {
				filtered = append(filtered, sub.Key)
			}
		}

		return filtered
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
		Keys: []string{""},
		Help: func() string {
			return ""
		},
		Suggestions: func(scene *EditorScene, partial string) []string {
			return nil
		},
		Run: func(scene *EditorScene, args []string) string {
			return ""
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
			return "jaskd fjalskdf\naaklsdjf\najsdflasd\nalksdjfa\nasdfasdf", nil
		},
		Subcommands: []*Command{},
	}
}

func writeCommand() *Command {
	return &Command{
		Key: "write",
		Help: func() string {
			return "write the scene to disk"
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
		Subcommands: nil,
	}
}
