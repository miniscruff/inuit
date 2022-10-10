package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AssetType string
type ContentType string

const (
	AssetImage    AssetType = "Image"
	AssetOpenType AssetType = "OpenType"

	ContentSprite ContentType = "Sprite"
	ContentSliced ContentType = "SlicedSprite"
	ContentLabel  ContentType = "Label"
)

type Asset struct {
	Type AssetType `json:"type"`
	File string    `json:"file"`
}

type SpriteContent struct {
	Image string `json:"image"`
}

type LabelContent struct {
	Font string `json:"font"`
	Size int    `json:"size"`
	DPI  int    `json:"dpi"`
}

type Content struct {
	Type   ContentType   `json:"type"`
	Sprite SpriteContent `json:"sprite,omitempty"`
	Label  LabelContent  `json:"label,omitempty"`
}

type Metadata struct {
	AssetsPath string `json:"assetsPath"`
	ScenesPath string `json:"scenesPath"`
}

type SceneMetadata struct {
	Name string `json:"name"`
}

type SceneVisual struct {
	Name          string `json:"name"`
	UseWindowSize bool   `json:"useWindowSize"`
}

type SceneData struct {
	Metadata SceneMetadata `json:"metadata"`
	Content  []string      `json:"content"`
	Visuals  []SceneVisual `json:"visuals"`
}

func LoadAssets(output *map[string]Asset) error {
	assetFileBytes, err := os.ReadFile(filepath.Join(InternalDir, AssetsFile))
	if err != nil {
		return err
	}

	json.Unmarshal(assetFileBytes, output)
	return nil
}

func LoadContent(output *map[string]Content) error {
	contentFileBytes, err := os.ReadFile(filepath.Join(InternalDir, ContentsFile))
	if err != nil {
		return err
	}

	json.Unmarshal(contentFileBytes, output)
	return nil
}

func LoadMetadata(output *Metadata) error {
	metadataFileBytes, err := os.ReadFile(filepath.Join(InternalDir, MetadataFile))
	if err != nil {
		return err
	}

	json.Unmarshal(metadataFileBytes, output)
	return nil
}

func LoadSceneData(output *SceneData, path string) error {
	sceneFileBytes, err := os.ReadFile(filepath.Join(InternalDir, path))
	if err != nil {
		return err
	}

	json.Unmarshal(sceneFileBytes, output)
	return nil
}
