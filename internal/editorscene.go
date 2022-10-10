package internal

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
