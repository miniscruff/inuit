package internal

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/miniscruff/igloo"
	"github.com/miniscruff/igloo/mathf"
)

type AssetType string
type ContentType string
type VisualType string

const (
	AssetImage    AssetType = "Image"
	AssetOpenType AssetType = "OpenType"

	ContentSprite ContentType = "Sprite"
	ContentSliced ContentType = "SlicedSprite"
	ContentFont   ContentType = "Font"

	EmptyVisualType  VisualType = "Empty"
	SpriteVisualType VisualType = "Sprite"
	LabelVisualType  VisualType = "Label"
)

type Asset struct {
	Type AssetType `json:"type"`
	File string    `json:"file"`
}

type BaseContent struct {
	Asset string `json:"asset"`
}

type SpriteContent struct {
	BaseContent
}

type FontContent struct {
	BaseContent
	Image string `json:"image"`
	Size  int    `json:"size"`
	DPI   int    `json:"dpi"`
}

type Content struct {
	Type   ContentType   `json:"type"`
	Sprite SpriteContent `json:"sprite,omitempty"`
	Font   FontContent   `json:"font,omitempty"`
}

type Metadata struct {
	AssetsPath string `json:"assetsPath"`
	ScenesPath string `json:"scenesPath"`
}

type SceneMetadata struct {
	Name string `json:"name"`
}

type SceneTransform struct {
	Position mathf.Vec2 `json:"position,omitempty"`
	Rotation float64    `json:"rotation,omitempty"`
	Anchor   mathf.Vec2 `json:"anchor,omitempty"`
	Width    float64    `json:"width,omitempty"`
	Height   float64    `json:"height,omitempty"`
}

type BaseVisualData struct {
	Content string `json:"content,omitempty"`
}

type SpriteVisualData struct {
	BaseVisualData
}

type LabelVisualData struct {
	BaseVisualData
}

type SceneVisual struct {
	Name          string           `json:"name"`
	Type          VisualType       `json:"type"`
	UseWindowSize bool             `json:"useWindowSize"`
	Visible       bool             `json:"visible"`
	Transform     SceneTransform   `json:"transform,omitempty"`
	Sprite        SpriteVisualData `json:"sprite,omitempty"`
	Label         LabelVisualData  `json:"label,omitempty"`
	Children      []*SceneVisual   `json:"children,omitempty"`
	Visual        *igloo.Visualer  `json:"-"`
}

type SceneData struct {
	Metadata SceneMetadata  `json:"metadata"`
	Content  []string       `json:"content"`
	Visuals  []*SceneVisual `json:"visuals"`
}

func LoadAssets(output *map[string]Asset) error {
	assetFileBytes, err := os.ReadFile(filepath.Join(InternalDir, AssetsFile))
	if err != nil {
		return err
	}

	return json.Unmarshal(assetFileBytes, output)
}

func LoadContent(output *map[string]Content) error {
	contentFileBytes, err := os.ReadFile(filepath.Join(InternalDir, ContentsFile))
	if err != nil {
		return err
	}

	return json.Unmarshal(contentFileBytes, output)
}

func LoadMetadata(output *Metadata) error {
	metadataFileBytes, err := os.ReadFile(filepath.Join(InternalDir, MetadataFile))
	if err != nil {
		return err
	}

	return json.Unmarshal(metadataFileBytes, output)
}

func LoadSceneData(output *SceneData, path string) error {
	sceneFileBytes, err := os.ReadFile(filepath.Join(InternalDir, path))
	if err != nil {
		return err
	}

	return json.Unmarshal(sceneFileBytes, output)
}

func SaveSceneData(output *SceneData, path string) error {
	outputFile, err := os.Create(filepath.Join(InternalDir, path))
	if err != nil {
		return err
	}

	defer outputFile.Close()

	outputBytes, err := json.Marshal(output)
	if err != nil {
		return err
	}

	_, err = outputFile.Write(outputBytes)
	return err
}
