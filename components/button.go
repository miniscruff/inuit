package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/miniscruff/igloo"
	"github.com/miniscruff/igloo/content"
	"github.com/miniscruff/igloo/graphics"
	"github.com/miniscruff/igloo/mathf"
)

type ButtonState string

const (
	NormalButtonState   ButtonState = "normal"
	OverButtonState     ButtonState = "over"
	ClickedButtonState  ButtonState = "clicked"
	ReleasedButtonState ButtonState = "released"
)

type Button struct {
	*mathf.Transform
	State *igloo.FSM[ButtonState]
}

func (b *Button) Update() {
	mx, my := ebiten.CursorPosition()
	mousePoint := mathf.Vec2FromInts(mx, my)
	inButton := b.Transform.Bounds().Contains(mousePoint)
	leftDown := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	if !inButton {
		b.State.Transition(NormalButtonState)
		return
	}

	switch b.State.Current() {
	case NormalButtonState:
		b.State.Transition(OverButtonState)
	case OverButtonState:
		if leftDown {
			b.State.Transition(ClickedButtonState)
		}
	case ClickedButtonState:
		if !leftDown {
			b.State.Transition(ReleasedButtonState)
		}
	case ReleasedButtonState:
		b.State.Transition(OverButtonState)
	}
}

func NewButton(sprite *graphics.SpriteVisual, normal, over, clicked *content.Sprite) *Button {
	state := igloo.NewFSM(
		NormalButtonState,
		igloo.NewFSMTransition(NormalButtonState, OverButtonState),
		igloo.NewFSMTransition(OverButtonState, NormalButtonState, ClickedButtonState),
		igloo.NewFSMTransition(ClickedButtonState, NormalButtonState, ReleasedButtonState),
		igloo.NewFSMTransition(ReleasedButtonState, OverButtonState),
	)

	state.OnTransition(NormalButtonState, func() {
		sprite.SetSprite(normal)
	})

	state.OnTransition(OverButtonState, func() {
		sprite.SetSprite(over)
	})

	state.OnTransition(ClickedButtonState, func() {
		sprite.SetSprite(clicked)
	})

	state.OnTransition(ReleasedButtonState, func() {
		sprite.SetSprite(normal)
	})

	return &Button{
		Transform: sprite.Transform,
		State:     state,
	}
}
