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
		{bridgeTile},
	})

	scene := lasagne.NewScene(lasagne.SceneParams{
		TileMap: tileMap,
		TileSet: *tileSet,
	})

	camera := lasagne.NewCamera()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		scene.Draw(camera)

		rl.EndDrawing()
	}
}
