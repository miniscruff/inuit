package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/miniscruff/igloo"
	"github.com/miniscruff/igloo/content"
	"github.com/miniscruff/igloo/graphics"
	"github.com/miniscruff/igloo/mathf"
	"github.com/miniscruff/inuit/internal"
)

type EditorScene struct {
	assets          *EditorAssets
	content         *EditorContent
	disposeHandlers []func()

	sceneAssets  map[string]any
	sceneContent map[string]any
	sceneTree    []*igloo.Visualer

	path string

	keys      []ebiten.Key
	lastInput int
}

func NewEditorScene(path string) *EditorScene {
	return &EditorScene{
		path: path,
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
		sceneData    internal.SceneData
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
	if err := internal.LoadSceneData(&sceneData, s.path); err != nil {
		return err
	}

	s.sceneAssets = make(map[string]any)
	s.sceneContent = make(map[string]any)
	s.sceneTree = make([]*igloo.Visualer, 0)

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

	for _, t := range sceneData.Visuals {
		rootVis := loadVisual(t, s.sceneContent, nil)
		s.sceneTree = append(s.sceneTree, rootVis)
	}

	return nil
}

func loadVisual(visual internal.SceneVisual, contentMap map[string]any, parent *igloo.Visualer) *igloo.Visualer {
	ww, wh := igloo.GetWindowSize()
	windowWidth := float64(ww)
	windowHeight := float64(wh)

	var newVis *igloo.Visualer
	switch visual.Type {
	case internal.EmptyVisualType:
		newVis = graphics.NewEmptyVisual().Visualer
	case internal.SpriteVisualType:
		sprite := graphics.NewSpriteVisual()
		sprite.SetSprite(contentMap[visual.Sprite.Content].(*content.Sprite))
		newVis = sprite.Visualer
	}

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
		parent.InsertChild(newVis)
	}

	for _, child := range visual.Children {
		loadVisual(child, contentMap, newVis)
	}

	return newVis
}

func (s *EditorScene) Update() {
}

func (s *EditorScene) Draw(dest *ebiten.Image) {
	offset := mathf.NewTransform()
	for _, r := range s.sceneTree {
		r.Layout(offset, r.Transform)
		r.Draw(dest)
	}
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
