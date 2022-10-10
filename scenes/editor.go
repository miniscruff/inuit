package scenes

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/miniscruff/igloo"
	"github.com/miniscruff/igloo/graphics"
	"github.com/miniscruff/inuit/internal"
)

type EditorState string

const (
	StateScene    EditorState = "scene"
	StateAssets   EditorState = "assets"
	StateContent  EditorState = "content"
	StateMetadata EditorState = "metadata"
)

type EditorScene struct {
	assets  *EditorAssets
	content *EditorContent
	tree    *EditorTree

	path         string
	assetData    map[string]internal.Asset
	contentData  map[string]internal.Content
	metadataData internal.Metadata
	sceneData    internal.SceneData

	state *igloo.FSM[EditorState]

	keys      []ebiten.Key
	lastInput int
}

func NewEditorScene(path string) *EditorScene {
	state := igloo.NewFSM(
		StateScene,
		igloo.NewFSMTransition(StateScene, StateAssets, StateContent, StateMetadata),
		igloo.NewFSMTransition(StateAssets, StateScene),
		igloo.NewFSMTransition(StateContent, StateScene),
		igloo.NewFSMTransition(StateMetadata, StateScene),
	)

	return &EditorScene{
		path:  path,
		state: state,
	}
}

func (s *EditorScene) PostSetup() (err error) {
	if err := internal.LoadAssets(&s.assetData); err != nil {
		return err
	}
	if err := internal.LoadContent(&s.contentData); err != nil {
		return err
	}
	if err := internal.LoadMetadata(&s.metadataData); err != nil {
		return err
	}
	if err := internal.LoadSceneData(&s.sceneData, s.path); err != nil {
		return err
	}

	s.state.OnTransitionTo(StateScene, func() {
	})
	s.state.OnTransitionFrom(StateScene, func() {
	})
	s.state.OnTransitionTo(StateAssets, func() {
		s.tree.Assets.SetVisible(true)
	})
	s.state.OnTransitionFrom(StateAssets, func() {
		s.tree.Assets.SetVisible(false)
	})
	s.state.OnTransitionTo(StateContent, func() {
		s.tree.Content.SetVisible(true)
	})
	s.state.OnTransitionFrom(StateContent, func() {
		s.tree.Content.SetVisible(false)
	})
	s.state.OnTransitionTo(StateMetadata, func() {
		s.tree.Metadata.SetVisible(true)
	})
	s.state.OnTransitionFrom(StateMetadata, func() {
		s.tree.Metadata.SetVisible(false)
	})

	var y float64 = 10
	for name, asset := range s.assetData {
		n := graphics.NewLabelVisual(s.content.SonoRegular18)
		n.SetText(name)
		n.SetY(y)
		y += 20
		s.tree.Assets.InsertChild(n.Visualer)

		t := graphics.NewLabelVisual(s.content.SonoRegular18)
		t.SetText(string(asset.Type))
		t.SetY(y)
		y += 20
		s.tree.Assets.InsertChild(t.Visualer)

		f := graphics.NewLabelVisual(s.content.SonoRegular18)
		f.SetText(asset.File)
		f.SetY(y)
		y += 20
		s.tree.Assets.InsertChild(f.Visualer)
	}

	return nil
}

func (s *EditorScene) Update() {
	hasInput := false
	for i := ebiten.KeyA; i <= ebiten.KeyZ; i++ {
		k := ebiten.Key(i)
		if inpututil.IsKeyJustReleased(k) {
			s.keys = append(s.keys, k)
			hasInput = true
		}
	}

	if !hasInput {
		s.lastInput++
	} else {
		s.lastInput = 0
	}

	if len(s.keys) > 2 {
		s.keys = s.keys[len(s.keys)-2:]
	}

	if len(s.keys) > 0 {
		if s.lastInput > 30 {
			s.clearInput()
			log.Println("cleared")
		}
	}

	if SlicesEqual(s.keys, internal.EditAssetsKeys) {
		s.state.Transition(StateAssets)
		s.clearInput()
	} else if SlicesEqual(s.keys, internal.EditContentsKeys) {
		s.state.Transition(StateContent)
		s.clearInput()
	} else if SlicesEqual(s.keys, internal.EditMetadataKeys) {
		s.state.Transition(StateMetadata)
		s.clearInput()
	} else if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		s.state.Transition(StateScene)
		s.clearInput()
	}
}

func (s *EditorScene) clearInput() {
	s.keys = nil
	s.lastInput = 0
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