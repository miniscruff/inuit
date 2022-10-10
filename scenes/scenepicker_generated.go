package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/opentype"

	"github.com/miniscruff/igloo"
	"github.com/miniscruff/igloo/content"
	"github.com/miniscruff/igloo/graphics"
	"github.com/miniscruff/igloo/mathf"
)

type ScenePickerAssets struct {
	NormalBackground  *ebiten.Image
	OverBackground    *ebiten.Image
	ClickedBackground *ebiten.Image
	SonoRegular       *opentype.Font
}

func (a *ScenePickerAssets) Dispose() {
	a.NormalBackground.Dispose()
	a.ClickedBackground.Dispose()
	a.OverBackground.Dispose()
	a.SonoRegular = nil
}

func NewScenePickerAssets(assetLoader *igloo.AssetLoader) (*ScenePickerAssets, error) {
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

	return &ScenePickerAssets{
		NormalBackground:  NormalBackground,
		ClickedBackground: ClickedBackground,
		OverBackground:    OverBackground,
		SonoRegular:       SonoRegular,
	}, nil
}

type ScenePickerContent struct {
	NormalBackground  *content.Sprite // 9 slice
	OverBackground    *content.Sprite // 9 slice
	ClickedBackground *content.Sprite // 9 slice

	SonoRegular18 *content.Label
	SonoRegular24 *content.Label
	SonoRegular36 *content.Label
}

func (c *ScenePickerContent) Dispose() {
	c.SonoRegular18.Close()
	c.SonoRegular24.Close()
	c.SonoRegular36.Close()
}

func NewScenePickerContent(assets *ScenePickerAssets) (*ScenePickerContent, error) {
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

	return &ScenePickerContent{
		NormalBackground:  NormalBackground,
		OverBackground:    OverBackground,
		ClickedBackground: ClickedBackground,
		SonoRegular18:     SonoRegular18,
		SonoRegular24:     SonoRegular24,
		SonoRegular36:     SonoRegular36,
	}, nil
}

type ScenePickerTree struct {
	UI *graphics.EmptyVisual
}

func NewScenePickerTree(content *ScenePickerContent) (*ScenePickerTree, error) {
	ww, wh := igloo.GetWindowSize()
	windowWidth := float64(ww)
	windowHeight := float64(wh)

	ui := graphics.NewEmptyVisual()

	ui.Transform.SetWidth(windowWidth)
	ui.Transform.SetHeight(windowHeight)

	ui.SetVisible(true)

	return &ScenePickerTree{
		UI: ui,
	}, nil
}

func (s *ScenePickerScene) Setup(assetLoader *igloo.AssetLoader) error {
	var err error

	s.assets, err = NewScenePickerAssets(assetLoader)
	if err != nil {
		return err
	}

	s.content, err = NewScenePickerContent(s.assets)
	if err != nil {
		return err
	}

	s.tree, err = NewScenePickerTree(s.content)
	if err != nil {
		return err
	}

	return nil
}

func (s *ScenePickerScene) Draw(dest *ebiten.Image) {
	offset := mathf.NewTransform()
	s.tree.UI.Visualer.Draw(dest, offset, s.tree.UI.Transform)
}

func (s *ScenePickerScene) Dispose() {
	s.assets.Dispose()
	s.content.Dispose()
	s.tree = nil
}
