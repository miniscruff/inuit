package scenes

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/miniscruff/inuit/internal"
)

type EditorScene struct {
	assets  *ScenePickerAssets
	content *ScenePickerContent
	tree    *ScenePickerTree

	assetData    map[string]internal.Asset
	contentData  map[string]internal.Content
	metadataData internal.Metadata
	sceneData    internal.SceneData
}

func NewEditorScene(path string) (*EditorScene, error) {
	assetFileBytes, err := os.ReadFile(filepath.Join(".inuit", "_assets.json"))
	if err != nil {
		return nil, err
	}

	contentFileBytes, err := os.ReadFile(filepath.Join(".inuit", "_content.json"))
	if err != nil {
		return nil, err
	}

	metadataFileBytes, err := os.ReadFile(filepath.Join(".inuit", "_metadata.json"))
	if err != nil {
		return nil, err
	}

	sceneFileBytes, err := os.ReadFile(filepath.Join(".inuit", path))
	if err != nil {
		return nil, err
	}

	e := &EditorScene{}

	json.Unmarshal(assetFileBytes, &e.assetData)
	json.Unmarshal(contentFileBytes, &e.contentData)
	json.Unmarshal(metadataFileBytes, &e.metadataData)
	json.Unmarshal(sceneFileBytes, &e.sceneData)

	return e, nil
}

func (s *EditorScene) PostSetup() {
}

func (s *EditorScene) Update() {
}
