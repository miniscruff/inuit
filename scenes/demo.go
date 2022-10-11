package scenes

type DemoScene struct {
	assets *DemoAssets
	content *DemoContent
	tree *DemoTree
}

func NewDemoScene() *DemoScene {
	return &DemoScene{}
}

func (s *DemoScene) PostSetup() error {
	return nil
}

func (s *DemoScene) Update() {
}
