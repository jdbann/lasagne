package lasagne

import (
	"cmp"
	"math"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Scene struct {
	objects []*Object
	tileMap *TileMap
	tileSet *TileSet
}

type SceneParams struct {
	Objects []*Object
	TileMap *TileMap
	TileSet *TileSet
}

func NewScene(params SceneParams) *Scene {
	return &Scene{
		objects: params.Objects,
		tileMap: params.TileMap,
		tileSet: params.TileSet,
	}
}

func (s *Scene) AddObject(object *Object) {
	s.objects = append(s.objects, object)
}

func (s Scene) Draw(camera *Camera) {
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

	cameraMatrix := combine(rl.MatrixMultiply,
		rl.MatrixTranslate(-camera.Target.X, -camera.Target.Y, 0), // Focus camera on target
		rl.MatrixTranslate(0.5, 0.5, 0),                           // Adjust for origin of tiles
		rl.MatrixScale(s.tileSet.size, s.tileSet.size, 1),         // Scale to tile size
		rl.MatrixScale(camera.Zoom, camera.Zoom, 1),               // Scale to zoom level
		rl.MatrixRotateZ(-camera.Rotation.X),                      // Horizontal rotation
	)

	tileSize := s.tileSet.size * camera.Zoom
	tileOrigin := rl.Vector2{X: tileSize / 2, Y: tileSize / 2}
	frameStep := float32(math.Sin(float64(camera.Rotation.Y))) * camera.Zoom / cosCameraY

	v := renderValues{
		frameStep:      frameStep,
		subframes:      int(math.Ceil(float64(frameStep))),
		cameraMatrix:   cameraMatrix,
		cameraRotation: camera.Rotation.X * rl.Rad2deg,
		zStep:          frameStep * s.tileSet.size,
	}

	// Get order to render world in based on camera rotation
	yDir := cmp.Compare(math.Cos(float64(camera.Rotation.X)), 0)
	yFrom, yTo := 0, len(s.tileMap.tiles[0])-1
	if yDir < 0 {
		yFrom, yTo = yTo, yFrom
	}

	xDir := cmp.Compare(math.Sin(float64(camera.Rotation.X)), 0)
	xFrom, xTo := 0, len(s.tileMap.tiles[0][0])-1
	if xDir < 0 {
		xFrom, xTo = xTo, xFrom
	}

	// Setup objects iterator with a function to compare positions with distance from camera
	compareInCamera := cmpVector3Distance(xDir, yDir)
	slices.SortFunc(s.objects, func(a, b *Object) int {
		return compareInCamera(a.Position, b.Position)
	})
	objectNext, objectDone, objectDrain := objectsIterator(s.objects, compareInCamera)

	for z := range s.tileMap.tiles {
		for yNext, yDone := iterator(yFrom, yTo); !yDone(); {
			y := yNext()
			for xNext, xDone := iterator(xFrom, xTo); !xDone(); {
				x := xNext()

				// Render any objects which lower in the render order than the current tile
				for !objectDone(rl.Vector3{X: float32(x), Y: float32(y), Z: float32(z)}) {
					object := objectNext()
					renderObject(object, camera, v)
				}

				tileIdx := s.tileMap.tiles[z][y][x]
				if tileIdx == -1 {
					continue
				}

				for frame := 0; frame < int(s.tileSet.size); frame++ {
					for subframe := 0; subframe <= v.subframes; subframe++ {
						tilePosition := rl.Vector3Transform(rl.Vector3{X: float32(x), Y: float32(y)}, v.cameraMatrix)
						rl.DrawTexturePro(
							s.tileSet.textures[tileIdx],
							rl.NewRectangle(float32(frame)*s.tileSet.size, 0, s.tileSet.size, s.tileSet.size),
							rl.NewRectangle(tilePosition.X, tilePosition.Y-v.zStep*float32(z)-v.frameStep*float32(frame)-float32(subframe), tileSize, tileSize),
							tileOrigin,
							v.cameraRotation,
							rl.White,
						)
					}
				}
			}
		}
	}

	// Render remaining objects
	for _, object := range objectDrain() {
		renderObject(object, camera, v)
	}
}

func (s Scene) MoveCamera(c *Camera, v rl.Vector3) {
	// Scale y in viewport to simulate vertical rotation
	cosCameraY := float32(math.Cos(float64(c.Rotation.Y)))

	d := rl.Vector3Transform(v, combine(
		rl.MatrixMultiply,
		rl.MatrixScale(1/c.Zoom/s.tileSet.size, 1/c.Zoom/s.tileSet.size/cosCameraY, 1),
		rl.MatrixRotateZ(c.Rotation.X),
	))
	c.Target = rl.Vector3Subtract(c.Target, d)
}

func (s *Scene) RemoveObject(object *Object) bool {
	idx := slices.Index(s.objects, object)
	if idx == -1 {
		return false
	}
	s.objects = slices.Delete(s.objects, idx, idx+1)
	return true
}

func cmpVector3Distance(xDir, yDir int) func(a, b rl.Vector3) int {
	return func(a, b rl.Vector3) int {
		if v := cmp.Compare(a.Z, b.Z); v != 0 {
			return v
		}

		if v := cmp.Compare(a.Y, b.Y); v != 0 {
			return v * yDir
		}

		return cmp.Compare(a.Y, b.Y) * xDir
	}
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
