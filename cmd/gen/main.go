package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"

	"github.com/miniscruff/igloo/mathf"
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
	genSceneTmpl = template.Must(template.New("genScene").Parse(`// Code generated by inuit DO NOT EDIT.

package scenes

import (
	{{- range .Imports}}
	"{{.}}"
	{{- end}}
)

type {{.Name}}Assets struct {
	{{- range .Assets}}
	{{.Name}} {{.GoType}}
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
	{{.Name}} {{.GoType}}
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
	{{- define "treeStruct" }}
	{{ .Name }} {{ .GoType }}
	{{- range .Children }}
	{{- template "treeStruct" . }}
	{{- end }}
	{{- end }}
	{{- range .Tree }}
	{{- template "treeStruct" . }}
	{{- end }}
}

func New{{.Name}}Tree(content *{{.Name}}Content) (*{{.Name}}Tree, error) {
	ww, wh := igloo.GetWindowSize()
	windowWidth := float64(ww)
	windowHeight := float64(wh)

	{{range .Tree }}
	{{ .Build }}
	{{- end }}

	return &{{.Name}}Tree{
		{{- define "retTree" }}
		{{ .Name }}: {{ .Name }},
		{{- range .Children }}
		{{- template "retTree" . }}
		{{- end }}
		{{- end }}
		{{- range .Tree }}
		{{- template "retTree" . }}
		{{- end }}
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

	{{- range .Tree}}
	s.tree.{{ .Name }}.Visualer.Layout(offset, s.tree.{{ .Name }}.Transform)
	s.tree.{{ .Name }}.Visualer.Draw(dest)
	{{- end}}
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
	GoType     string
	File       string
	Dispose    string
	LoadMethod string
}

type GenContent struct {
	Name    string
	GoType  string
	Dispose string
	Create  func() string
}

type GenTree struct {
	Name     string
	GoType   string
	Build    string
	Children []GenTree
}

type GeneratedSceneContext struct {
	Name     string
	Imports  []string
	Assets   []GenAsset
	Contents []GenContent
	Tree     []GenTree
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
	var tree []GenTree
	for _, v := range scene.Visuals {
		tree = append(tree, buildTree(v))
	}

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
		Tree:     tree,
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
		case internal.ContentFont:
			a.Name = c.Font.Asset
			a.LoadMethod = "LoadOpenType"
			a.Dispose = " = nil"
			a.GoType = "*opentype.Font"
		case internal.ContentSprite:
			a.Name = c.Sprite.Asset
			a.LoadMethod = "LoadImage"
			a.Dispose = ".Dispose()"
			a.GoType = "*ebiten.Image"
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
		case internal.ContentFont:
			a.Dispose = ".Close()"
			a.GoType = "*content.Font"
			a.Create = func() string {
				return fmt.Sprintf(`%vFont, err := opentype.NewFace(assets.%v, &opentype.FaceOptions{
		Size: %v,
		DPI:  %v,
	})
	if err != nil {
		return nil, err
	}

	%v := &content.Font{
		Face: %vFont,
	}
`, a.Name, c.Font.Asset, c.Font.Size, c.Font.DPI, a.Name, a.Name)
			}
		case internal.ContentSprite:
			a.GoType = "*content.Sprite"
			a.Create = func() string {
				return fmt.Sprintf(`%v := &content.Sprite{
	Image: assets.%v,
}`, a.Name, c.Sprite.Asset)
			}
		}

		genContent = append(genContent, a)
	}

	return genContent
}

func buildTree(visual internal.SceneVisual) GenTree {
	t := GenTree{
		Name: visual.Name,
	}

	switch visual.Type {
	case internal.EmptyVisualType:
		t.GoType = "*graphics.EmptyVisual"
	case internal.SpriteVisualType:
		t.GoType = "*graphics.SpriteVisual"
	case internal.LabelVisualType:
		t.GoType = "*graphics.LabelVisual"
	}

	for _, c := range visual.Children {
		t.Children = append(t.Children, buildTree(c))
	}

	t.Build = genVisualBuild(t, visual, t.Children)

	return t
}

func genVisualBuild(t GenTree, visual internal.SceneVisual, children []GenTree) string {
	var b strings.Builder

	switch visual.Type {
	case internal.EmptyVisualType:
		writeFormat(&b, "%v := graphics.NewEmptyVisual()", t.Name)
	case internal.SpriteVisualType:
		writeFormat(&b, "%v := graphics.NewSpriteVisual()", t.Name)
		writeFormat(&b, "%v.SetSprite(content.%v)", t.Name, visual.Sprite.Content)
	case internal.LabelVisualType:
		writeFormat(&b, "%v := graphics.NewLabelVisual()", t.Name)
		writeFormat(&b, "%v.SetFont(content.%v)", t.Name, visual.Label.Content)
	}

	if visual.Visible {
		writeFormat(&b, "%v.SetVisible(true)", t.Name)
	}

	if visual.UseWindowSize {
		writeFormat(&b, "%v.Transform.SetWidth(windowWidth)", t.Name)
		writeFormat(&b, "%v.Transform.SetHeight(windowHeight)", t.Name)
	}

	transform := visual.Transform
	condWrite(&b,
		transform.Position.X != 0,
		"%v.SetX(%v)",
		t.Name, transform.Position.X,
	)
	condWrite(&b,
		transform.Position.Y != 0,
		"%v.SetY(%v)",
		t.Name, transform.Position.Y,
	)
	condWrite(&b,
		transform.Rotation != 0,
		"%v.SetRotation(%v)",
		t.Name, transform.Rotation,
	)
	condWrite(&b,
		transform.Anchor != mathf.Vec2Zero,
		"%v.SetAnchor(mathf.Vec2{X: %v, Y: %v})",
		t.Name, transform.Anchor.X, transform.Anchor.Y,
	)
	condWrite(&b,
		transform.Width != 0, // temp?
		"%v.SetWidth(%v)",
		t.Name, transform.Width,
	)
	condWrite(&b,
		transform.Height != 0, // temp?
		"%v.SetHeight(%v)",
		t.Name, transform.Height,
	)

	b.WriteString("\n")
	for _, c := range children {
		writeFormat(&b, c.Build)
	}

	for _, c := range children {
		writeFormat(&b, "%v.InsertChild(%v.Visualer)", t.Name, c.Name)
	}

	return b.String()
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

		var buffer bytes.Buffer

		err = generateGeneratedScene(&buffer, scene, assets, content, metadata)
		if err != nil {
			log.Fatal(err)
		}

		formattedBytes, err := imports.Process(genScenePath, buffer.Bytes(), &imports.Options{
			Fragment: true,
			Comments: true,
		})
		if err != nil {
			log.Fatal(err)
		}

		genSceneWriter, err := os.Create(genScenePath)
		if err != nil {
			log.Fatal(err)
		}
		defer genSceneWriter.Close()

		_, err = genSceneWriter.Write(formattedBytes)
		if err != nil {
			log.Fatal(err)
		}

	}
}

func condWrite(w io.StringWriter, cond bool, format string, args ...any) {
	if !cond {
		return
	}

	writeFormat(w, format, args...)
}

func writeFormat(w io.StringWriter, format string, args ...any) {
	w.WriteString(fmt.Sprintf(format+"\n", args...))
}
