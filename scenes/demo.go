package scenes

import "log"

type DemoScene struct {
	assets  *DemoAssets
	content *DemoContent
	tree    *DemoTree
}

func NewDemoScene() *DemoScene {
	return &DemoScene{}
}

func (s *DemoScene) PostSetup() error {
	log.Println("post setup")
	return nil
}

func (s *DemoScene) Update() {
}
