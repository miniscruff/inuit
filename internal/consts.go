package internal

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	InternalDir = ".inuit"
	AssetsFile = "_assets.json"
	ContentsFile = "_content.json"
	MetadataFile = "_metadata.json"
)

var (
	EditAssetsKeys = []ebiten.Key{ebiten.KeyE, ebiten.KeyA}
	EditContentsKeys = []ebiten.Key{ebiten.KeyE, ebiten.KeyC}
	EditMetadataKeys = []ebiten.Key{ebiten.KeyE, ebiten.KeyM}
)
