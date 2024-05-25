package lasagne

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Camera struct {
	Target   rl.Vector2
	Rotation rl.Vector2
	Zoom     float32
}

func NewCamera() Camera {
	return Camera{
		Target:   rl.Vector2{X: 0, Y: 0},
		Rotation: rl.Vector2{X: rl.Pi / 4, Y: 2 * rl.Pi / 3},
		Zoom:     2,
	}
}

func (c Camera) ApplyMatrix() {
	// Translate to put target position in center of screen
	rl.Translatef(
		float32(rl.GetScreenWidth()/2)+c.Target.X,
		float32(rl.GetScreenHeight()/2)+c.Target.Y,
		0,
	)
}
