package commands

func writeCommand() *Command {
	return &Command{
		Key: "write",
		Help: func() string {
			return "save scene to disk"
		},
		Run: writeAction,
	}
}

func writeAction(editor Editor, args []string) (string, error) {
	err := SaveSceneData(editor.SceneData(), "testing.json")
	return "scene saved", err
}
