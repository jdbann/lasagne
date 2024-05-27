package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jdbann/lasagne"
)

func main() {
	rl.InitWindow(1280, 720, "lasagne - dungeon example")
	defer rl.CloseWindow()

	tileSet := lasagne.NewTileSet(16)
	bridgeTexture := rl.LoadTexture("assets/Bridge_strip16.png")
	bridgeTile := tileSet.AddTile(bridgeTexture)
	floorTexture := rl.LoadTexture("assets/FloorCentrePlain_strip16.png")
	floorTile := tileSet.AddTile(floorTexture)
	floorPatternTexture := rl.LoadTexture("assets/FloorCentrePattern_strip16.png")
	floorPatternTile := tileSet.AddTile(floorPatternTexture)
	wallTexture := rl.LoadTexture("assets/WallCentreA_strip16.png")
	wallTile := tileSet.AddTile(wallTexture)

	tileMap := lasagne.NewTileMap([][][]int{
		{
			{floorTile, floorTile, floorTile, -1},
			{floorTile, floorTile, floorTile, bridgeTile},
			{floorPatternTile, floorTile, floorTile, -1},
			{floorPatternTile, floorPatternTile, floorTile, -1},
		},
		{
			{-1, -1, -1, -1},
			{-1, wallTile, -1, -1},
			{-1, -1, -1, -1},
			{-1, -1, -1, -1},
		},
	})

	barrelTexture := rl.LoadTexture("assets/Barrel_strip8.png")
	barrelObject := &lasagne.Object{
		Position: rl.Vector3{
			X: 2,
			Y: 2,
			Z: 1,
		},
		Texture: barrelTexture,
		Size: rl.Vector3{
			X: 14,
			Y: 14,
			Z: 8,
		},
	}
	chairTexture := rl.LoadTexture("assets/Chair_strip12.png")
	chairObject := &lasagne.Object{
		Position: rl.Vector3{
			X: 1,
			Y: 3,
			Z: 1,
		},
		Texture: chairTexture,
		Size: rl.Vector3{
			X: 12,
			Y: 12,
			Z: 12,
		},
	}

	scene := lasagne.NewScene(lasagne.SceneParams{
		TileMap: tileMap,
		TileSet: tileSet,
	})
	scene.AddObject(barrelObject)
	scene.AddObject(chairObject)

	camera := lasagne.NewCamera()
	camera.Target = rl.Vector3{X: 2, Y: 2, Z: 1}

	for !rl.WindowShouldClose() {
		mouseDelta := rl.GetMouseDelta()
		if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			camera.Rotation.X = rl.Wrap(camera.Rotation.X+(mouseDelta.X/float32(rl.GetScreenWidth()))*rl.Pi*2, 0, rl.Pi*2)
			camera.Rotation.Y = rl.Clamp(camera.Rotation.Y+(mouseDelta.Y/float32(rl.GetScreenHeight()))*(rl.Pi/2), 0, rl.Pi/2-0.001)
		} else if rl.IsKeyDown(rl.KeySpace) {
			scene.MoveCamera(camera, rl.Vector3{X: mouseDelta.X, Y: mouseDelta.Y})
		}

		camera.Zoom = rl.Clamp(camera.Zoom+rl.GetMouseWheelMove(), 0.5, 8)

		rl.BeginDrawing()
		rl.ClearBackground(rl.NewColor(10, 10, 24, 255))

		scene.Draw(camera)

		rl.EndDrawing()
	}
}
