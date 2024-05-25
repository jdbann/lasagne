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

	// Scale y in viewport to simulate vertical rotation
	cosCameraY := float32(math.Cos(float64(camera.Rotation.Y)))
	rl.Scalef(1, cosCameraY, 1)

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

	frameStep := float32(math.Sin(float64(camera.Rotation.Y))) * camera.Zoom / cosCameraY
	subframes := int(math.Ceil(float64(frameStep)))

	yFrom, yTo := 0, len(s.tileMap.tiles)-1
	if math.Cos(float64(camera.Rotation.X)) < 0 {
		yFrom, yTo = yTo, yFrom
	}

	xFrom, xTo := 0, len(s.tileMap.tiles[0])-1
	if math.Sin(float64(camera.Rotation.X)) < 0 {
		xFrom, xTo = xTo, xFrom
	}

	for yNext, yDone := iterator(yFrom, yTo); !yDone(); {
		y := yNext()
		for xNext, xDone := iterator(xFrom, xTo); !xDone(); {
			x := xNext()
			tileIdx := s.tileMap.tiles[y][x]
			if tileIdx == -1 {
				continue
			}

			for frame := 0; frame < int(s.tileSet.size); frame++ {
				for subframe := 0; subframe <= subframes; subframe++ {
					tilePosition := rl.Vector3Transform(rl.Vector3{X: float32(x), Y: float32(y)}, tileMatrix)
					rl.DrawTexturePro(
						s.tileSet.textures[tileIdx],
						rl.NewRectangle(float32(frame)*s.tileSet.size, 0, s.tileSet.size, s.tileSet.size),
						rl.NewRectangle(tilePosition.X, tilePosition.Y-frameStep*float32(frame)-float32(subframe), tileSize, tileSize),
						tileOrigin,
						tileRotation,
						rl.White,
					)
				}
			}
		}
	}
}

func (s Scene) MoveCamera(c Camera, v rl.Vector2) Camera {
	// Scale y in viewport to simulate vertical rotation
	cosCameraY := float32(math.Cos(float64(c.Rotation.Y)))

	d := rl.Vector2Transform(v, combine(
		rl.MatrixMultiply,
		rl.MatrixScale(1/c.Zoom/s.tileSet.size, 1/c.Zoom/s.tileSet.size/cosCameraY, 1),
		rl.MatrixRotateZ(c.Rotation.X),
	))
	newTarget := rl.Vector2Subtract(c.Target, d)
	c.Target = newTarget
	return c
}

func combine[T any](combineFn func(T, T) T, in ...T) T {
	out := in[0]
	for i := 1; i < len(in); i++ {
		out = combineFn(out, in[i])
	}
	return out
}

func iterator(from, to int) (func() int, func() bool) {
	currentVal, step := from, 1
	if from > to {
		step = -1
	}
	done := false

	nextFn := func() int {
		out := currentVal
		if out == to {
			done = true
		} else {
			currentVal += step
		}
		return out
	}

	doneFn := func() bool { return done }

	return nextFn, doneFn
}
