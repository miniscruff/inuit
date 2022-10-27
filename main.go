package main

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	math_rand "math/rand"

	"embed"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/miniscruff/igloo"

	"github.com/miniscruff/inuit/scenes"
	"github.com/miniscruff/inuit/commands"
)

var (
	//go:embed assets
	assetsFS embed.FS
)

func main() {
	// slightly better random seed
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	// initialize our game and window state
	igloo.InitGame(igloo.GameConfig{
		Fsys: assetsFS,
		AssetsPath: "assets",
	})
	igloo.SetWindowSize(1024, 768)
	igloo.SetScreenSize(1024, 768)
	ebiten.SetWindowTitle("inuit")

	sceneNames, err := commands.ExistingScenes()
	if err != nil {
		fmt.Printf("failure to find scenes: %v\n", err)
	}

	scene := scenes.NewEditorScene(sceneNames[0]+".json")

	// push our starting scene and run
	igloo.Push(scene)
	if err := igloo.Run(); err != nil {
		fmt.Printf("run complete: %v\n", err)
	}
}
