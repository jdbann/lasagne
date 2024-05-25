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

	tileMap := lasagne.NewTileMap([][]int{
		{bridgeTile, bridgeTile},
		{bridgeTile, bridgeTile},
	})

	scene := lasagne.NewScene(lasagne.SceneParams{
		TileMap: tileMap,
		TileSet: *tileSet,
	})

	camera := lasagne.NewCamera()
	camera.Target = rl.Vector2{X: 1, Y: 1}

	for !rl.WindowShouldClose() {
		camera.Rotation.X = rl.Clamp((rl.GetMousePosition().X/float32(rl.GetScreenWidth()))*rl.Pi*2, 0, rl.Pi*2)
		camera.Rotation.Y = rl.Clamp((rl.GetMousePosition().Y/float32(rl.GetScreenHeight()))*(rl.Pi/2), 0, rl.Pi/2-0.001)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		scene.Draw(camera)

		rl.EndDrawing()
	}
}
