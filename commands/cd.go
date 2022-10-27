package commands

func cdCommand() *Command {
	return &Command{
		Key: "cd",
		Help: func() string {
			return "change active object up or down"
		},
		Suggestions: func(editor Editor, partial []string) []string {
			if len(partial) != 1 {
				return nil
			}

			var search []*SceneVisual
			if editor.Visual() == nil {
				search = editor.SceneData().Visuals
			} else {
				search = editor.Visual().Children
			}

			return Filter(partial[0], search, func(child *SceneVisual) string {
				return child.Name
			})
		},
		Run: cdAction,
	}
}

func cdAction(editor Editor, args []string) (string, error) {
	if len(args) == 0 {
		editor.SetVisual(nil)
		return "", nil
	}
	if args[0] == ".." {
		if editor.Visual() != nil {
			editor.SetVisual(editor.Visual().Parent)
		}
		return "", nil
	}

	var search []*SceneVisual
	if editor.Visual() == nil {
		search = editor.SceneData().Visuals
	} else {
		search = editor.Visual().Children
	}

	for _, child := range search {
		if child.Name == args[0] {
			editor.SetVisual(child)
			return "", nil
		}
	}

	return "", nil
}
