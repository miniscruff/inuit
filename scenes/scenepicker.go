package scenes

import (
	"image/color"

	"github.com/miniscruff/igloo"
	"github.com/miniscruff/igloo/graphics"
	"github.com/miniscruff/igloo/mathf"
	"github.com/miniscruff/inuit/components"
	"github.com/miniscruff/inuit/internal"
)

type ScenePickerScene struct {
	assets  *ScenePickerAssets
	content *ScenePickerContent
	tree    *ScenePickerTree

	buttons []*components.Button
}

func NewScenePickerScene() (*ScenePickerScene, error) {
	return &ScenePickerScene{}, nil
}

func (s *ScenePickerScene) PostSetup() error {
	sw, _ := igloo.GetScreenSize()

	var y float64 = 80

	scenes, err := internal.ExistingScenes()
	if err != nil {
		return err
	}

	for _, sceneName := range scenes {
		sceneName := sceneName

		button := graphics.NewSpriteVisual()
		button.SetSprite(s.content.NormalBackground)
		button.Transform.SetAnchor(mathf.Vec2MiddleCenter)
		button.Transform.SetX(float64(sw) / 2)
		button.Transform.SetY(y)
		button.Transform.SetSize(250, 40)
		button.SetVisible(true)
		s.tree.UI.InsertChild(button.Visualer)
		y += 80

		text := graphics.NewLabelVisual()
		text.SetFont(s.content.SonoRegular18)
		text.SetText(sceneName)
		text.ColorM.ScaleWithColor(color.White)
		text.Transform.SetAnchor(mathf.Vec2MiddleCenter)
		text.SetVisible(true)
		button.InsertChild(text.Visualer)

		b := components.NewButton(
			button,
			s.content.NormalBackground,
			s.content.OverBackground,
			s.content.ClickedBackground,
		)
		b.State.OnTransitionTo(components.ReleasedButtonState, func() {
			igloo.Pop()
			igloo.Push(NewEditorScene(sceneName + ".json"))
		})

		s.buttons = append(s.buttons, b)
	}

	return nil
}

func (s *ScenePickerScene) Update() {
	for _, b := range s.buttons {
		b.Update()
	}
}
