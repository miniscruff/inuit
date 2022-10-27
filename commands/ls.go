package commands

import "strings"

func lsCommand() *Command {
	return &Command{
		Key: "ls",
		Help: func() string {
			return "list the current object and children names"
		},
		Run: lsAction,
	}
}

func lsAction(editor Editor, args []string) (string, error) {
	var builder strings.Builder
	var search []*SceneVisual

	if editor.Visual() != nil {
		WriteFormat(&builder, ".%v", editor.Visual().Name)
		search = editor.Visual().Children
	} else {
		WriteFormat(&builder, "No selection")
		search = editor.SceneData().Visuals
	}

	for _, child := range search {
		WriteFormat(
			&builder,
			"./%v",
			child.Name,
		)
	}

	return builder.String(), nil
}
