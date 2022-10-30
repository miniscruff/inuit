package commands

func writeCommand() *Command {
	return &Command{
		Key: "write",
		Help: func() string {
			return "save scene to disk"
		},
		Run: func(editor Editor, args []string) (string, error) {
			err := SaveSceneData(editor.SceneData(), editor.Path())
			return "scene saved", err
		},
	}
}
