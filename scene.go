package lasagne

import (
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

func (s Scene) Draw(_ Camera) {
	for y := range s.tileMap.tiles {
		for x, tile := range s.tileMap.tiles[y] {
			rl.DrawTexturePro(
				s.tileSet.textures[tile],
				rl.NewRectangle(240, 0, s.tileSet.size, s.tileSet.size),
				rl.NewRectangle(float32(x)*s.tileSet.size, float32(y)*s.tileSet.size, s.tileSet.size, s.tileSet.size),
				rl.Vector2{},
				0,
				rl.White,
			)
		}
	}
}
