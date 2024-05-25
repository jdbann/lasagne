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

	tileMatrix := combine(rl.MatrixMultiply,
		rl.MatrixTranslate(-camera.Target.X, -camera.Target.Y, 0), // Focus camera on target
		rl.MatrixTranslate(0.5, 0.5, 0),                           // Adjust for origin of tiles
		rl.MatrixScale(s.tileSet.size, s.tileSet.size, 1),         // Scale to tile size
		rl.MatrixScale(camera.Zoom, camera.Zoom, 1),               // Scale to zoom level
		rl.MatrixRotateZ(-camera.Rotation.X),                      // Horizontal rotation
	)

	tileSize := s.tileSet.size * camera.Zoom
	tileOrigin := rl.Vector2{X: tileSize / 2, Y: tileSize / 2}
	tileRotation := camera.Rotation.X * rl.Rad2deg

	// Scale y in viewport to simulate vertical rotation
	cosCameraY := float32(math.Cos(float64(camera.Rotation.Y)))
	rl.Scalef(1, cosCameraY, 1)

	for y := range s.tileMap.tiles {
		for x, tile := range s.tileMap.tiles[y] {
			tilePosition := rl.Vector3Transform(rl.Vector3{X: float32(x), Y: float32(y)}, tileMatrix)
			rl.DrawTexturePro(
				s.tileSet.textures[tile],
				rl.NewRectangle(240, 0, s.tileSet.size, s.tileSet.size),
				rl.NewRectangle(tilePosition.X, tilePosition.Y, tileSize, tileSize),
				tileOrigin,
				tileRotation,
				rl.White,
			)
		}
	}
}

func combine[T any](combineFn func(T, T) T, in ...T) T {
	out := in[0]
	for i := 1; i < len(in); i++ {
		out = combineFn(out, in[i])
	}
	return out
}
