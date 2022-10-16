package scenes

import (
	"fmt"
	"image/color"
	"strings"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/miniscruff/igloo"
	"github.com/miniscruff/igloo/content"
	"github.com/miniscruff/igloo/graphics"
	"github.com/miniscruff/igloo/mathf"
	"github.com/miniscruff/inuit/internal"
)

const wasdSpeed = 5

type TextEditorState string

const (
	TextEditorClosed TextEditorState = "closed"
	TextEditorOpen   TextEditorState = "open"
)

type EditorScene struct {
	assets          *EditorAssets
	content         *EditorContent
	disposeHandlers []func()
	sceneData       internal.SceneData

	sceneAssets  map[string]any
	sceneContent map[string]any

	path     string
	commands []*Command

	activeVisual *internal.SceneVisual
	offset       *mathf.Transform

	inputState *igloo.FSM[TextEditorState]

	textInput        []rune
	textBuffer       []rune
	inputRoot        *graphics.EmptyVisual
	textInputLabel   *graphics.LabelVisual
	inputResponse    *graphics.LabelVisual
	inputBackground  *graphics.SpriteVisual
	suggestionsLabel *graphics.LabelVisual
}

func NewEditorScene(path string) *EditorScene {
	return &EditorScene{
		path:     path,
		offset:   mathf.NewTransform(),
		commands: buildCommands(),
	}
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
		assetData    map[string]internal.Asset
		contentData  map[string]internal.Content
		metadataData internal.Metadata
	)

	if err := internal.LoadAssets(&assetData); err != nil {
		return err
	}
	if err := internal.LoadContent(&contentData); err != nil {
		return err
	}
	if err := internal.LoadMetadata(&metadataData); err != nil {
		return err
	}
	if err := internal.LoadSceneData(&s.sceneData, s.path); err != nil {
		return err
	}

	s.sceneAssets = make(map[string]any)
	s.sceneContent = make(map[string]any)

	for k, a := range assetData {
		switch a.Type {
		case internal.AssetImage:
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
		case internal.ContentSprite:
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
	s.textInputLabel.Transform.SetX(5)
	s.textInputLabel.Transform.SetY(windowHeight - 20)
	s.textInputLabel.Transform.SetAnchor(mathf.Vec2TopLeft)
	s.textInputLabel.SetVisible(true)

	s.inputResponse = graphics.NewLabelVisual()
	s.inputResponse.SetFont(s.content.SonoRegular18)
	s.inputResponse.ColorM.ScaleWithColor(color.White)
	s.inputResponse.Transform.SetAnchor(mathf.Vec2TopLeft)
	s.inputResponse.SetVisible(true)

	s.inputBackground = graphics.NewSpriteVisual()
	s.inputBackground.SetSprite(s.content.NormalBackground)
	s.inputBackground.SetWidth(windowWidth)
	s.inputBackground.SetAnchor(mathf.Vec2TopLeft)
	s.inputBackground.SetVisible(true)

	s.suggestionsLabel = graphics.NewLabelVisual()
	s.suggestionsLabel.SetFont(s.content.SonoRegular18)
	s.suggestionsLabel.Transform.SetY(windowHeight - 24)
	s.suggestionsLabel.Transform.SetAnchor(mathf.Vec2BottomLeft)
	s.suggestionsLabel.SetVisible(true)

	s.inputRoot = graphics.NewEmptyVisual()
	s.inputRoot.Transform.SetWidth(windowWidth)
	s.inputRoot.Transform.SetHeight(windowHeight)
	s.inputRoot.SetVisible(false)

	s.inputRoot.InsertChild(s.inputBackground.Visualer)
	s.inputRoot.InsertChild(s.textInputLabel.Visualer)
	s.inputRoot.InsertChild(s.inputResponse.Visualer)
	s.inputRoot.InsertChild(s.suggestionsLabel.Visualer)

	s.inputState = igloo.NewFSM(
		TextEditorClosed,
		igloo.NewFSMTransition(TextEditorClosed, TextEditorOpen),
		igloo.NewFSMTransition(TextEditorOpen, TextEditorClosed),
	)

	s.inputState.OnTransitionTo(TextEditorClosed, func() {
		s.inputRoot.SetVisible(false)
	})

	s.inputState.OnTransitionTo(TextEditorOpen, func() {
		s.inputRoot.SetVisible(true)
		s.textBuffer = nil
		s.textInputLabel.SetText("_")
		s.inputBackground.SetHeight(20)
	})

	return nil
}

func loadVisual(visual *internal.SceneVisual, contentMap map[string]any, parent *internal.SceneVisual) {
	ww, wh := igloo.GetWindowSize()
	windowWidth := float64(ww)
	windowHeight := float64(wh)

	var newVis *igloo.Visualer

	switch visual.Type {
	case internal.EmptyVisualType:
		newVis = graphics.NewEmptyVisual().Visualer
	case internal.SpriteVisualType:
		spriteContent := contentMap[visual.Sprite.Content].(*content.Sprite)
		spriteVis := graphics.NewSpriteVisual()
		spriteVis.SetSprite(spriteContent)
		newVis = spriteVis.Visualer
	}

	visual.Visual = newVis
	newVis.SetVisible(visual.Visible)
	newVis.SetPosition(visual.Transform.Position)
	newVis.SetAnchor(visual.Transform.Anchor)
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
	switch s.inputState.Current() {
	case TextEditorClosed:
		if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
			s.inputState.Transition(TextEditorOpen)
			break
		}

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
	case TextEditorOpen:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			s.inputState.Transition(TextEditorClosed)
			break
		}

		s.textInput = ebiten.AppendInputChars(s.textInput[:0])
		textChanged := false
		for _, r := range s.textInput {
			if ebiten.IsKeyPressed(ebiten.KeyShift) {
				s.textBuffer = append(s.textBuffer, unicode.ToUpper(r))
			} else {
				s.textBuffer = append(s.textBuffer, r)
			}
			textChanged = true
		}

		backspaceDur := inpututil.KeyPressDuration(ebiten.KeyBackspace)
		if len(s.textBuffer) > 0 {
			if backspaceDur == 1 || (backspaceDur-30 >= 0 && backspaceDur%5 == 0) {
				s.textBuffer = s.textBuffer[:len(s.textBuffer)-1]
				textChanged = true
			}
		}

		if inpututil.IsKeyJustReleased(ebiten.KeyTab) {
			cmd := string(s.textBuffer)
			suggestions := buildSuggestions(s, s.commands, cmd)
			if len(suggestions) == 1 {
				split := strings.Split(cmd, " ")
				if len(split) == 1 {
					s.textBuffer = []rune(suggestions[0])
				} else {
					s.textBuffer = []rune(strings.Join(split[:len(split)-1], " "))
					s.textBuffer = append(s.textBuffer, []rune(" ")...)
					s.textBuffer = append(s.textBuffer, []rune(suggestions[0])...)
				}
				s.textBuffer = append(s.textBuffer, []rune(" ")...)
				textChanged = true
			}
		}

		if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
			cmd := strings.Trim(string(s.textBuffer), "\n ")
			if len(cmd) > 0 {
				output, err := runCommand(s, s.commands, cmd)
				if err != nil {
					output = fmt.Sprintf("unable to run command: %v\n%v", cmd, err.Error())
				}

				s.inputResponse.SetText(output)
				s.inputBackground.SetHeight(s.inputResponse.NaturalHeight())
				s.textBuffer = nil
				textChanged = true
			}
		}

		if textChanged {
			cmd := string(s.textBuffer)
			suggestions := buildSuggestions(s, s.commands, cmd)
			if len(suggestions) > 0 {
				s.suggestionsLabel.SetText(strings.Join(suggestions, "\n"))
			} else {
				s.suggestionsLabel.SetText("")
			}
			s.textInputLabel.SetText(cmd + "_")
		}
	}
}

func (s *EditorScene) Draw(dest *ebiten.Image) {
	for _, v := range s.sceneData.Visuals {
		var baseOffset *mathf.Transform
		if true {
			baseOffset = s.offset
		} else {
			baseOffset = mathf.NewTransform()
		}

		v.Visual.Layout(baseOffset, v.Visual.Transform)
		v.Visual.Draw(dest)
	}

	fixedOffset := mathf.NewTransform()
	s.inputRoot.Visualer.Layout(fixedOffset, s.inputRoot.Transform)
	s.inputRoot.Visualer.Draw(dest)
}

func (s *EditorScene) Dispose() {
	s.assets.Dispose()
	s.content.Dispose()
	for _, h := range s.disposeHandlers {
		h()
	}
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
