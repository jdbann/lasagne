package lasagne

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Camera struct {
	Target   rl.Vector3
	Rotation rl.Vector2
	Zoom     float32
}

func NewCamera() Camera {
	return Camera{
		Target:   rl.Vector3{X: 0, Y: 0, Z: 0},
		Rotation: rl.Vector2{X: rl.Pi / 4, Y: rl.Pi / 3},
		Zoom:     2,
	}
}
