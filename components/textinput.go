package components

import (
	"strings"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/miniscruff/igloo"
)

type TextEditorState string

const (
	TextEditorClosed TextEditorState = "closed"
	TextEditorOpen   TextEditorState = "open"
)

type TextEditor struct {
	State   *igloo.FSM[TextEditorState]
	Changed func(text string)
	Tab     func(text string) string
	Submit  func(text string)

	currentInput []rune
	text         []rune
}

func NewTextEditor() *TextEditor {
	return &TextEditor{
		Changed: func(text string) {},
		Tab:     func(text string) string { return text },
		Submit:  func(text string) {},
		State: igloo.NewFSM(
			TextEditorClosed,
			igloo.NewFSMTransition(TextEditorClosed, TextEditorOpen),
			igloo.NewFSMTransition(TextEditorOpen, TextEditorClosed),
		),
	}
}

func (e *TextEditor) Update() {
	switch e.State.Current() {
	case TextEditorClosed:
		if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
			e.State.Transition(TextEditorOpen)
		}

	case TextEditorOpen:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			e.State.Transition(TextEditorClosed)
			e.currentInput = nil
			e.text = nil
			return
		}

		textChanged := false
		e.currentInput = ebiten.AppendInputChars(e.currentInput[:0])

		for _, r := range e.currentInput {
			if ebiten.IsKeyPressed(ebiten.KeyShift) {
				e.text = append(e.text, unicode.ToUpper(r))
			} else {
				e.text = append(e.text, r)
			}
			textChanged = true
		}

		backspaceDur := inpututil.KeyPressDuration(ebiten.KeyBackspace)
		if len(e.text) > 0 {
			if backspaceDur == 1 || (backspaceDur-30 >= 0 && backspaceDur%5 == 0) {
				e.text = e.text[:len(e.text)-1]
				textChanged = true
			}
		}

		if inpututil.IsKeyJustReleased(ebiten.KeyTab) {
			text := string(e.text)
			newText := e.Tab(text)
			if newText != text {
				textChanged = true
				e.text = []rune(newText)
			}
		}

		if inpututil.IsKeyJustReleased(ebiten.KeyEnter) && len(e.text) > 0 {
			cleanedText := strings.Trim(string(e.text), "\n ")
			e.Submit(cleanedText)
			e.text = nil
			textChanged = true
		}

		if textChanged {
			e.Changed(string(e.text))
		}
	}
}
