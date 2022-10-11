package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/miniscruff/inuit/internal"
)

var (
	baseSceneTmpl = template.Must(template.New("baseScene").Parse(`package scenes

type {{.Name}}Scene struct {
	assets *{{.Name}}Assets
	content *{{.Name}}Content
	tree *{{.Name}}Tree
}

func New{{.Name}}Scene() *{{.Name}}Scene {
	return &{{.Name}}Scene{}
}

func (s *{{.Name}}Scene) PostSetup() error {
	return nil
}

func (s *{{.Name}}Scene) Update() {
}
`))
	genSceneTmpl = template.Must(template.New("genScene").Parse(`package scenes

import (
	{{- range .Imports}}
	"{{.}}"
	{{- end}}
)

type {{.Name}}Assets struct {
	{{- range .Assets}}
	{{.Name}} {{.Type}}
	{{- end}}
}

func (a *{{.Name}}Assets) Dispose() {
	{{- range .Assets}}
	a.{{.Name}}{{.Dispose}}
	{{- end}}
}

func New{{.Name}}Assets(assetLoader *igloo.AssetLoader) (*{{.Name}}Assets, error) {
	{{- range .Assets}}
	{{.Name}}, err := assetLoader.{{.LoadMethod}}("{{.File}}")
	if err != nil {
		return nil, err
	}
	{{end}}

	return &{{.Name}}Assets{
		{{- range .Assets}}
		{{.Name}}: {{.Name}},
		{{- end}}
	}, nil
}

type {{.Name}}Content struct {
	{{- range .Contents}}
	{{.Name}} {{.Type}}
	{{- end}}
}

func (c *{{.Name}}Content) Dispose() {
	{{- range .Contents}}
	{{- if .Dispose}}
	c.{{.Name}}{{.Dispose}}
	{{- end }}{{- end}}
}

func New{{.Name}}Content(assets *{{.Name}}Assets) (*{{.Name}}Content, error) {
	{{- range .Contents}}
		{{ call .Create }}
	{{end}}

	return &{{.Name}}Content{
		{{- range .Contents}}
		{{ .Name }}: {{.Name}},
		{{- end}}
	}, nil
}

type {{.Name}}Tree struct {
	UI *graphics.EmptyVisual
}

func New{{.Name}}Tree(content *{{.Name}}Content) (*{{.Name}}Tree, error) {
	ww, wh := igloo.GetWindowSize()
	windowWidth := float64(ww)
	windowHeight := float64(wh)

	ui := graphics.NewEmptyVisual()

	ui.Transform.SetWidth(windowWidth)
	ui.Transform.SetHeight(windowHeight)

	ui.SetVisible(true)

	return &{{.Name}}Tree{
		UI: ui,
	}, nil
}

func (s *{{.Name}}Scene) Setup(assetLoader *igloo.AssetLoader) error {
	var err error

	s.assets, err = New{{.Name}}Assets(assetLoader)
	if err != nil {
		return err
	}

	s.content, err = New{{.Name}}Content(s.assets)
	if err != nil {
		return err
	}

	s.tree, err = New{{.Name}}Tree(s.content)
	if err != nil {
		return err
	}

	return nil
}

func (s *{{.Name}}Scene) Draw(dest *ebiten.Image) {
	offset := mathf.NewTransform()
	s.tree.UI.Visualer.Layout(offset, s.tree.UI.Transform)

	s.tree.UI.Visualer.Draw(dest)
}

func (s *{{.Name}}Scene) Dispose() {
	s.assets.Dispose()
	s.content.Dispose()
	s.tree = nil
}
`))
)

type BaseSceneContext struct {
	Name string
}

type GenAsset struct {
	Name       string
	Type       string
	File       string
	Dispose    string
	LoadMethod string
}

type GenContent struct {
	Name    string
	Type    string
	Dispose string
	Create  func() string
}

type GeneratedSceneContext struct {
	Name     string
	Imports  []string
	Assets   []GenAsset
	Contents []GenContent
}

func generateBaseScene(w io.Writer, scene internal.SceneData) error {
	ctx := BaseSceneContext{
		Name: scene.Metadata.Name,
	}
	return baseSceneTmpl.Execute(w, ctx)
}

func generateGeneratedScene(
	w io.Writer,
	scene internal.SceneData,
	assets map[string]internal.Asset,
	content map[string]internal.Content,
	metadata internal.Metadata,
) error {
	ctx := GeneratedSceneContext{
		Name: scene.Metadata.Name,
		Imports: []string{
			"github.com/hajimehoshi/ebiten/v2",
			"golang.org/x/image/font/opentype",
			"github.com/miniscruff/igloo",
			"github.com/miniscruff/igloo/mathf",
			"github.com/miniscruff/igloo/graphics",
			"github.com/miniscruff/igloo/content",
		},
		Assets:   findAllAssets(assets, content, scene.Content),
		Contents: findAllContent(content, scene.Content),
	}
	return genSceneTmpl.Execute(w, ctx)
}

func findAllAssets(assets map[string]internal.Asset, content map[string]internal.Content, contentKeys []string) []GenAsset {
	var genAssets []GenAsset
	seen := make(map[string]struct{})

	for _, key := range contentKeys {
		c := content[key]
		var a GenAsset
		switch c.Type {
		case internal.ContentLabel:
			a.Name = c.Label.Font
			a.LoadMethod = "LoadOpenType"
			a.Dispose = " = nil"
			a.Type = "*opentype.Font"
		case internal.ContentSprite:
			a.Name = c.Sprite.Image
			a.LoadMethod = "LoadImage"
			a.Dispose = ".Dispose()"
			a.Type = "*ebiten.Image"
		}

		if _, found := seen[a.Name]; found {
			continue
		}

		seen[a.Name] = struct{}{}

		a.File = assets[a.Name].File
		genAssets = append(genAssets, a)
	}

	return genAssets
}

func findAllContent(content map[string]internal.Content, contentKeys []string) []GenContent {
	var genContent []GenContent
	seen := make(map[string]struct{})

	for _, key := range contentKeys {
		if _, found := seen[key]; found {
			continue
		}
		seen[key] = struct{}{}

		c := content[key]
		a := GenContent{
			Name: key,
		}
		switch c.Type {
		case internal.ContentLabel:
			a.Dispose = ".Close()"
			a.Type = "*content.Label"
			a.Create = func() string {
				return fmt.Sprintf(`%vFont, err := opentype.NewFace(assets.%v, &opentype.FaceOptions{
		Size: %v,
		DPI:  %v,
	})
	if err != nil {
		return nil, err
	}

	%v := &content.Label{
		Face: %vFont,
	}
`, a.Name, c.Label.Font, c.Label.Size, c.Label.DPI, a.Name, a.Name)
			}
		case internal.ContentSprite:
			a.Type = "*content.Sprite"
			a.Create = func() string {
				return fmt.Sprintf(`%v := &content.Sprite{
	Image: assets.%v,
}`, a.Name, c.Sprite.Image)
			}
		}

		genContent = append(genContent, a)
	}

	return genContent
}

func main() {
	var assets map[string]internal.Asset
	var content map[string]internal.Content
	var metadata internal.Metadata
	var scenes []string
	var err error

	if err = internal.LoadAssets(&assets); err != nil {
		log.Fatal(err)
	}

	if err = internal.LoadContent(&content); err != nil {
		log.Fatal(err)
	}

	if err = internal.LoadMetadata(&metadata); err != nil {
		log.Fatal(err)
	}

	if scenes, err = internal.ExistingScenes(); err != nil {
		log.Fatal(err)
	}

	for _, fileName := range scenes {
		var scene internal.SceneData
		if err = internal.LoadSceneData(&scene, fileName+".json"); err != nil {
			log.Fatal(err)
		}

		rootPath := strings.ToLower(filepath.Join(metadata.ScenesPath, scene.Metadata.Name))
		baseScenePath := rootPath + ".go"
		genScenePath := rootPath + "_generated.go"

		// only generate a base scene if it does not already exist
		if _, err := os.Stat(baseScenePath); os.IsNotExist(err) {
			baseSceneWriter, err := os.Create(baseScenePath)
			if err != nil {
				log.Fatal(err)
			}

			err = generateBaseScene(baseSceneWriter, scene)
			if err != nil {
				log.Fatal(err)
			}
		}

		genSceneWriter, err := os.Create(genScenePath)
		if err != nil {
			log.Fatal(err)
		}

		err = generateGeneratedScene(genSceneWriter, scene, assets, content, metadata)
		if err != nil {
			log.Fatal(err)
		}
	}
}
