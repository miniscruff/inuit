package scenes

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/miniscruff/igloo"
	"github.com/miniscruff/igloo/content"
	"github.com/miniscruff/igloo/graphics"
	"github.com/miniscruff/igloo/mathf"
	"github.com/miniscruff/inuit/commands"
	"github.com/miniscruff/inuit/components"
)

const wasdSpeed = 5

type EditorScene struct {
	assets          *EditorAssets
	content         *EditorContent
	disposeHandlers []func()
	sceneData       commands.SceneData

	sceneAssets  map[string]any
	sceneContent map[string]any

	path     string
	commands *commands.Commands

	activeVisual *commands.SceneVisual
	offset       *mathf.Transform

	commandInput     *components.TextEditor
	inputRoot        *graphics.EmptyVisual
	textInputLabel   *graphics.LabelVisual
	inputResponse    *graphics.LabelVisual
	suggestionsLabel *graphics.LabelVisual
}

func NewEditorScene(path string) *EditorScene {
	s := &EditorScene{
		path:   path,
		offset: mathf.NewTransform(),
	}
	s.commands = commands.NewCommands(s)
	return s
}

func (s *EditorScene) Setup(assetLoader *igloo.AssetLoader) (err error) {
	s.assets, err = NewEditorAssets(assetLoader)
	if err != nil {
		return err
	}

	s.content, err = NewEditorContent(s.assets)
	if err != nil {
		return err
	}

	var (
		assetData    map[string]commands.Asset
		contentData  map[string]commands.Content
		metadataData commands.Metadata
	)

	if err := commands.LoadAssets(&assetData); err != nil {
		return err
	}
	if err := commands.LoadContent(&contentData); err != nil {
		return err
	}
	if err := commands.LoadMetadata(&metadataData); err != nil {
		return err
	}

	if err := commands.LoadSceneData(&s.sceneData, s.path); err != nil {
		return fmt.Errorf("unable to load scene data: %w", err)
	}

	s.sceneAssets = make(map[string]any)
	s.sceneContent = make(map[string]any)

	for k, a := range assetData {
		switch a.Type {
		case commands.AssetImage:
			// TODO: can not use the asset loader when loading an external asset
			img, err := assetLoader.LoadImage(a.File)
			if err != nil {
				return err
			}

			s.sceneAssets[k] = img
			s.disposeHandlers = append(s.disposeHandlers, func() {
				img.Dispose()
			})
		}
	}

	for k, c := range contentData {
		switch c.Type {
		case commands.ContentSprite:
			sprite := &content.Sprite{
				Image: s.sceneAssets[c.Sprite.Asset].(*ebiten.Image),
				// TODO: other sprite attributes
			}

			s.sceneContent[k] = sprite
		}
	}

	for _, t := range s.sceneData.Visuals {
		loadVisual(t, s.sceneContent, nil)
	}

	ww, wh := igloo.GetWindowSize()
	windowWidth := float64(ww)
	windowHeight := float64(wh)

	s.textInputLabel = graphics.NewLabelVisual()
	s.textInputLabel.SetFont(s.content.SonoRegular18)
	s.textInputLabel.ColorM.ScaleWithColor(color.White)
	s.textInputLabel.Transform.SetPivot(mathf.Vec2BottomLeft)
	s.textInputLabel.Transform.SetAnchors(mathf.SidesBottomLeft)
	s.textInputLabel.Transform.SetPosition(mathf.Vec2{X: 5, Y: -5})
	s.textInputLabel.SetVisible(true)

	s.inputResponse = graphics.NewLabelVisual()
	s.inputResponse.SetFont(s.content.SonoRegular18)
	s.inputResponse.ColorM.ScaleWithColor(color.White)
	s.inputResponse.Transform.SetPosition(mathf.Vec2{X: 5, Y: 5})
	s.inputResponse.SetVisible(true)

	s.suggestionsLabel = graphics.NewLabelVisual()
	s.suggestionsLabel.SetFont(s.content.SonoRegular18)
	s.suggestionsLabel.Transform.SetPivot(mathf.Vec2BottomLeft)
	s.suggestionsLabel.Transform.SetAnchors(mathf.SidesBottomLeft)
	s.suggestionsLabel.Transform.SetPosition(mathf.Vec2{X: 5, Y: -35})
	s.suggestionsLabel.SetVisible(true)

	s.inputRoot = graphics.NewEmptyVisual()
	s.inputRoot.Transform.SetWidth(windowWidth)
	s.inputRoot.Transform.SetHeight(windowHeight)
	s.inputRoot.SetVisible(false)

	s.inputRoot.InsertChild(s.suggestionsLabel.Visualer)
	s.inputRoot.InsertChild(s.textInputLabel.Visualer)
	s.inputRoot.InsertChild(s.inputResponse.Visualer)

	inputBackground := graphics.NewSpriteVisual()
	inputBackground.SetSprite(s.content.NormalBackground)
	inputBackground.SetAnchors(mathf.SidesStretchBoth)
	inputBackground.SetOffsets(mathf.Sides{Left: -5, Right: -5, Top: -5, Bottom: -5})
	inputBackground.SetVisible(true)
	s.textInputLabel.InsertChild(inputBackground.Visualer)

	suggestionBackground := graphics.NewSpriteVisual()
	suggestionBackground.SetSprite(s.content.NormalBackground)
	suggestionBackground.SetAnchors(mathf.SidesStretchBoth)
	suggestionBackground.SetOffsets(mathf.Sides{Left: -5, Right: -5, Top: -5, Bottom: -5})
	suggestionBackground.SetVisible(true)
	s.suggestionsLabel.InsertChild(suggestionBackground.Visualer)

	responseBackground := graphics.NewSpriteVisual()
	responseBackground.SetSprite(s.content.NormalBackground)
	responseBackground.SetAnchors(mathf.SidesStretchBoth)
	responseBackground.SetOffsets(mathf.Sides{Left: -5, Right: -5, Top: -5, Bottom: -5})
	responseBackground.SetVisible(true)
	s.inputResponse.InsertChild(responseBackground.Visualer)

	s.commandInput = components.NewTextEditor()
	s.commandInput.State.OnTransitionTo(components.TextEditorClosed, func() {
		s.inputRoot.SetVisible(false)
	})
	s.commandInput.State.OnTransitionTo(components.TextEditorOpen, func() {
		s.inputRoot.SetVisible(true)
		s.textInputLabel.SetText("_")
	})
	s.commandInput.Changed = func(text string) {
		suggestions := s.commands.BuildSuggestions(text)
		if len(suggestions) > 0 {
			s.suggestionsLabel.SetText(strings.Join(suggestions, "\n"))
		} else {
			s.suggestionsLabel.SetText("<none>")
		}

		s.textInputLabel.SetText(text + "_")
	}
	s.commandInput.Tab = func(text string) string {
		suggestions := s.commands.BuildSuggestions(text)
		if len(suggestions) != 1 {
			return text
		}

		split := strings.Split(text, " ")
		if len(split) == 1 {
			return suggestions[0] + " "
		}

		return strings.Join(split[:len(split)-1], " ") + " " + suggestions[0] + " "
	}
	s.commandInput.Submit = func(text string) {
		output, err := s.commands.Run(text)
		if err != nil {
			output = fmt.Sprintf("unable to run command: %v\n%v", text, err.Error())
		}

		s.inputResponse.SetText(output)
	}

	return nil
}

func loadVisual(visual *commands.SceneVisual, contentMap map[string]any, parent *commands.SceneVisual) {
	ww, wh := igloo.GetWindowSize()
	windowWidth := float64(ww)
	windowHeight := float64(wh)

	var newVis *igloo.Visualer

	switch visual.Type {
	case commands.EmptyVisualType:
		newVis = graphics.NewEmptyVisual().Visualer
	case commands.SpriteVisualType:
		spriteContent := contentMap[visual.Sprite.Content].(*content.Sprite)
		spriteVis := graphics.NewSpriteVisual()
		spriteVis.SetSprite(spriteContent)
		newVis = spriteVis.Visualer
	}

	visual.Visual = newVis
	newVis.SetVisible(visual.Visible)
	newVis.SetPosition(visual.Transform.Position)
	newVis.SetAnchors(visual.Transform.Anchors)
	newVis.SetOffsets(visual.Transform.Offsets)
	newVis.SetPivot(visual.Transform.Pivot)
	newVis.SetRotation(visual.Transform.Rotation)

	if visual.UseWindowSize {
		newVis.SetWidth(windowWidth)
		newVis.SetHeight(windowHeight)
	} else {
		newVis.SetWidth(visual.Transform.Width)
		newVis.SetHeight(visual.Transform.Height)
	}

	if parent != nil {
		parent.Visual.InsertChild(newVis)
		visual.Parent = parent
	}

	for _, child := range visual.Children {
		loadVisual(child, contentMap, visual)
	}
}

func (s *EditorScene) Update() {
	s.commandInput.Update()

	if s.commandInput.State.Current() == components.TextEditorClosed {
		dir := mathf.Vec2Zero
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			dir.Y = -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			dir.Y = 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			dir.X = 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			dir.X = -1
		}
		if dir != mathf.Vec2Zero {
			s.offset.Translate(dir.Unit().MulScalar(wasdSpeed))
		}
	}
}

func (s *EditorScene) Draw(dest *ebiten.Image) {
	for _, v := range s.sceneData.Visuals {
		v.Visual.Layout(v.Visual.Transform, nil)
		v.Visual.Draw(dest)
	}

	s.inputRoot.Visualer.Layout(s.inputRoot.Transform, nil)
	s.inputRoot.Visualer.Draw(dest)
}

func (s *EditorScene) Dispose() {
	s.assets.Dispose()
	s.content.Dispose()
	for _, h := range s.disposeHandlers {
		h()
	}
}

// commands.Editor implementations

func (s *EditorScene) Visual() *commands.SceneVisual {
	return s.activeVisual
}

func (s *EditorScene) SetVisual(visual *commands.SceneVisual) {
	s.activeVisual = visual
}

func (s *EditorScene) Path() string {
	return s.path
}

func (s *EditorScene) SceneData() *commands.SceneData {
	return &s.sceneData
}

func (s *EditorScene) Commands() *commands.Commands {
	return s.commands
}
