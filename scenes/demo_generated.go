package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/miniscruff/igloo"
	"github.com/miniscruff/igloo/content"
	"github.com/miniscruff/igloo/graphics"
	"github.com/miniscruff/igloo/mathf"
	"golang.org/x/image/font/opentype"
)

type DemoAssets struct {
	NormalBackground  *ebiten.Image
	OverBackground    *ebiten.Image
	ClickedBackground *ebiten.Image
	SonoRegular       *opentype.Font
}

func (a *DemoAssets) Dispose() {
	a.NormalBackground.Dispose()
	a.OverBackground.Dispose()
	a.ClickedBackground.Dispose()
	a.SonoRegular = nil
}

func NewDemoAssets(assetLoader *igloo.AssetLoader) (*DemoAssets, error) {
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

	return &DemoAssets{
		NormalBackground:  NormalBackground,
		OverBackground:    OverBackground,
		ClickedBackground: ClickedBackground,
		SonoRegular:       SonoRegular,
	}, nil
}

type DemoContent struct {
	NormalBackground  *content.Sprite
	OverBackground    *content.Sprite
	ClickedBackground *content.Sprite
	SonoRegular18     *content.Label
	SonoRegular24     *content.Label
	SonoRegular36     *content.Label
}

func (c *DemoContent) Dispose() {
	c.SonoRegular18.Close()
	c.SonoRegular24.Close()
	c.SonoRegular36.Close()
}

func NewDemoContent(assets *DemoAssets) (*DemoContent, error) {
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

	return &DemoContent{
		NormalBackground:  NormalBackground,
		OverBackground:    OverBackground,
		ClickedBackground: ClickedBackground,
		SonoRegular18:     SonoRegular18,
		SonoRegular24:     SonoRegular24,
		SonoRegular36:     SonoRegular36,
	}, nil
}

type DemoTree struct {
	UI *graphics.EmptyVisual
}

func NewDemoTree(content *DemoContent) (*DemoTree, error) {
	ww, wh := igloo.GetWindowSize()
	windowWidth := float64(ww)
	windowHeight := float64(wh)

	ui := graphics.NewEmptyVisual()

	ui.Transform.SetWidth(windowWidth)
	ui.Transform.SetHeight(windowHeight)

	ui.SetVisible(true)

	return &DemoTree{
		UI: ui,
	}, nil
}

func (s *DemoScene) Setup(assetLoader *igloo.AssetLoader) error {
	var err error

	s.assets, err = NewDemoAssets(assetLoader)
	if err != nil {
		return err
	}

	s.content, err = NewDemoContent(s.assets)
	if err != nil {
		return err
	}

	s.tree, err = NewDemoTree(s.content)
	if err != nil {
		return err
	}

	return nil
}

func (s *DemoScene) Draw(dest *ebiten.Image) {
	offset := mathf.NewTransform()
	s.tree.UI.Visualer.Layout(offset, s.tree.UI.Transform)

	s.tree.UI.Visualer.Draw(dest)
}

func (s *DemoScene) Dispose() {
	s.assets.Dispose()
	s.content.Dispose()
	s.tree = nil
}