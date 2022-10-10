package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/opentype"

	"github.com/miniscruff/igloo"
	"github.com/miniscruff/igloo/content"
	"github.com/miniscruff/igloo/graphics"
	"github.com/miniscruff/igloo/mathf"
)

type EditorAssets struct {
	NormalBackground  *ebiten.Image
	OverBackground    *ebiten.Image
	ClickedBackground *ebiten.Image
	SonoRegular       *opentype.Font
}

func (a *EditorAssets) Dispose() {
	a.NormalBackground.Dispose()
	a.ClickedBackground.Dispose()
	a.OverBackground.Dispose()
	a.SonoRegular = nil
}

func NewEditorAssets(assetLoader *igloo.AssetLoader) (*EditorAssets, error) {
	NormalBackground, err := assetLoader.LoadImage("normal_background.png")
	if err != nil {
		return nil, err
	}
	OverBackground, err := assetLoader.LoadImage("over_background.png")
	if err != nil {
		return nil, err
	}
	ClickedBackground, err := assetLoader.LoadImage("clicked_background.png")
	if err != nil {
		return nil, err
	}

	SonoRegular, err := assetLoader.LoadOpenType("Sono-Regular.ttf")
	if err != nil {
		return nil, err
	}

	return &EditorAssets{
		NormalBackground:  NormalBackground,
		ClickedBackground: ClickedBackground,
		OverBackground:    OverBackground,
		SonoRegular:       SonoRegular,
	}, nil
}

type EditorContent struct {
	NormalBackground  *content.Sprite // 9 slice
	OverBackground    *content.Sprite // 9 slice
	ClickedBackground *content.Sprite // 9 slice

	SonoRegular18 *content.Label
	SonoRegular24 *content.Label
	SonoRegular36 *content.Label
}

func (c *EditorContent) Dispose() {
	c.SonoRegular18.Close()
	c.SonoRegular24.Close()
	c.SonoRegular36.Close()
}

func NewEditorContent(assets *EditorAssets) (*EditorContent, error) {
	NormalBackground := &content.Sprite{
		Image: assets.NormalBackground,
	}

	OverBackground := &content.Sprite{
		Image: assets.OverBackground,
	}

	ClickedBackground := &content.Sprite{
		Image: assets.ClickedBackground,
	}

	SonoRegular18Font, err := opentype.NewFace(assets.SonoRegular, &opentype.FaceOptions{
		Size: 18,
		DPI:  72,
	})
	if err != nil {
		return nil, err
	}

	SonoRegular18 := &content.Label{
		Face: SonoRegular18Font,
	}

	SonoRegular24Font, err := opentype.NewFace(assets.SonoRegular, &opentype.FaceOptions{
		Size: 24,
		DPI:  72,
	})
	if err != nil {
		return nil, err
	}

	SonoRegular24 := &content.Label{
		Face: SonoRegular24Font,
	}

	SonoRegular36Font, err := opentype.NewFace(assets.SonoRegular, &opentype.FaceOptions{
		Size: 36,
		DPI:  72,
	})
	if err != nil {
		return nil, err
	}

	SonoRegular36 := &content.Label{
		Face: SonoRegular36Font,
	}

	return &EditorContent{
		NormalBackground:  NormalBackground,
		OverBackground:    OverBackground,
		ClickedBackground: ClickedBackground,
		SonoRegular18:     SonoRegular18,
		SonoRegular24:     SonoRegular24,
		SonoRegular36:     SonoRegular36,
	}, nil
}

type EditorTree struct {
	Editor   *graphics.EmptyVisual
	Assets   *graphics.EmptyVisual
	Content  *graphics.EmptyVisual
	Metadata *graphics.EmptyVisual
}

func NewEditorTree(content *EditorContent) (*EditorTree, error) {
	ww, wh := igloo.GetWindowSize()
	windowWidth := float64(ww)
	windowHeight := float64(wh)

	Editor := graphics.NewEmptyVisual()
	Editor.Transform.SetWidth(windowWidth)
	Editor.Transform.SetHeight(windowHeight)
	Editor.SetVisible(true)

	Assets := graphics.NewEmptyVisual()
	Assets.Transform.SetWidth(windowWidth)
	Assets.Transform.SetHeight(windowHeight)

	Content := graphics.NewEmptyVisual()
	Content.Transform.SetWidth(windowWidth)
	Content.Transform.SetHeight(windowHeight)

	Metadata := graphics.NewEmptyVisual()
	Metadata.Transform.SetWidth(windowWidth)
	Metadata.Transform.SetHeight(windowHeight)

	return &EditorTree{
		Editor:   Editor,
		Assets:   Assets,
		Content:  Content,
		Metadata: Metadata,
	}, nil
}

func (s *EditorScene) Setup(assetLoader *igloo.AssetLoader) error {
	var err error

	s.assets, err = NewEditorAssets(assetLoader)
	if err != nil {
		return err
	}

	s.content, err = NewEditorContent(s.assets)
	if err != nil {
		return err
	}

	s.tree, err = NewEditorTree(s.content)
	if err != nil {
		return err
	}

	return nil
}

func (s *EditorScene) Draw(dest *ebiten.Image) {
	offset := mathf.NewTransform()
	s.tree.Editor.Visualer.Draw(dest, offset, s.tree.Editor.Transform)
	s.tree.Assets.Visualer.Draw(dest, offset, s.tree.Assets.Transform)
	s.tree.Content.Visualer.Draw(dest, offset, s.tree.Content.Transform)
	s.tree.Metadata.Visualer.Draw(dest, offset, s.tree.Metadata.Transform)
}

func (s *EditorScene) Dispose() {
	s.assets.Dispose()
	s.content.Dispose()
	s.tree = nil
}
