package lasagne

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Scene struct {
	tileMap TileMap
	tileSet TileSet
}

type SceneParams struct {
	TileMap TileMap
	TileSet TileSet
}

func NewScene(params SceneParams) Scene {
	return Scene{
		tileMap: params.TileMap,
		tileSet: params.TileSet,
	}
}

func (s Scene) Draw(camera Camera) {
	// Save and restore current transformation matrix
	rl.PushMatrix()
	defer rl.PopMatrix()

	// Translate to center of screen
	rl.Translatef(
		float32(rl.GetScreenWidth()/2),
		float32(rl.GetScreenHeight()/2),
		0,
	)

	// Scale to zoom level, scaling y to simulate camera pitch
	cosCameraY := float32(math.Cos(float64(camera.Rotation.Y)))
	rl.Scalef(camera.Zoom, camera.Zoom*cosCameraY, 1)

	// Rotate around target
	rl.Rotatef(camera.Rotation.X*rl.Rad2deg, 0, 0, 1)

	// Translate to camera target
	rl.Translatef(
		-camera.Target.X*s.tileSet.size,
		-camera.Target.Y*s.tileSet.size,
		0,
	)

	origin := rl.Vector2{X: s.tileSet.size / 2, Y: s.tileSet.size / 2}

	for y := range s.tileMap.tiles {
		for x, tile := range s.tileMap.tiles[y] {
			rl.DrawTexturePro(
				s.tileSet.textures[tile],
				rl.NewRectangle(240, 0, s.tileSet.size, s.tileSet.size),
				rl.NewRectangle((float32(x)+.5)*s.tileSet.size, (float32(y)+.5)*s.tileSet.size, s.tileSet.size, s.tileSet.size),
				origin,
				0,
				rl.White,
			)
		}
	}
}
