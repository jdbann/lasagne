package lasagne

import rl "github.com/gen2brain/raylib-go/raylib"

type Object struct {
	Position rl.Vector3
	Texture  rl.Texture2D
	Size     rl.Vector3
}

type renderValues struct {
	frameStep      float32
	subframes      int
	cameraMatrix   rl.Matrix
	cameraRotation float32
	zStep          float32
}

func objectsIterator(objects []Object, compareFn func(a, b rl.Vector3) int) (func() Object, func(rl.Vector3) bool, func() []Object) {
	objectIdx := 0

	nextFn := func() Object {
		nextIdx := objectIdx
		objectIdx++
		return objects[nextIdx]
	}

	doneFn := func(v rl.Vector3) bool {
		return objectIdx == len(objects) || compareFn(objects[objectIdx].Position, v) >= 0
	}

	drainFn := func() []Object {
		return objects[objectIdx:]
	}

	return nextFn, doneFn, drainFn
}

func renderObject(object Object, camera Camera, v renderValues) {
	objectSize := rl.Vector2{X: object.Size.X * camera.Zoom, Y: object.Size.Y * camera.Zoom}
	objectOrigin := rl.Vector2{X: objectSize.X / 2, Y: objectSize.Y / 2}
	for frame := 0; frame < int(object.Size.Z); frame++ {
		for subframe := 0; subframe <= v.subframes; subframe++ {
			position := rl.Vector3Transform(object.Position, v.cameraMatrix)
			rl.DrawTexturePro(
				object.Texture,
				rl.NewRectangle(float32(frame)*object.Size.X, 0, object.Size.X, object.Size.Y),
				rl.NewRectangle(position.X, position.Y-v.zStep*object.Position.Z-v.frameStep*float32(frame)-float32(subframe), objectSize.X, objectSize.Y),
				objectOrigin,
				v.cameraRotation,
				rl.White,
			)
		}
	}
}
