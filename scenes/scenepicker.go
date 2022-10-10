package scenes

import (
	"image/color"

	"github.com/miniscruff/inuit/components"
	"github.com/miniscruff/igloo"
	"github.com/miniscruff/igloo/graphics"
	"github.com/miniscruff/igloo/mathf"
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

func (s *ScenePickerScene) PostSetup() {
	sw, _ := igloo.GetScreenSize()

	sceneButton := graphics.NewSpriteVisual(s.content.NormalBackground)
	sceneButton.Transform.SetAnchor(mathf.Vec2MiddleCenter)
	sceneButton.Transform.SetX(float64(sw)/2)
	sceneButton.Transform.SetY(80)
	sceneButton.Transform.SetSize(250, 40)
	sceneButton.SetVisible(true)

	sceneButtonText := graphics.NewLabelVisual(s.content.SonoRegular18)
	sceneButtonText.SetText("New Scene")
	sceneButtonText.ColorM.ScaleWithColor(color.White)
	sceneButtonText.Transform.SetAnchor(mathf.Vec2MiddleCenter)
	sceneButtonText.SetVisible(true)
	sceneButton.InsertChild(sceneButtonText.Visualer)

	newButton := components.NewButton(
		sceneButton,
		s.content.NormalBackground,
		s.content.OverBackground,
		s.content.ClickedBackground,
	)
	newButton.State.OnTransition(components.ReleasedButtonState, func() {
		// create a new scene...
	})

	s.buttons = []*components.Button{
		newButton,
	}

	s.tree.UI.InsertChild(sceneButton.Visualer)
}

func (s *ScenePickerScene) Update() {
	for _, b := range s.buttons {
		b.Update()
	}
}
