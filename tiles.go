package lasagne

import rl "github.com/gen2brain/raylib-go/raylib"

type TileSet struct {
	size     float32
	textures []rl.Texture2D
}

func NewTileSet(size float32) *TileSet {
	return &TileSet{
		size: size,
	}
}

func (ts *TileSet) AddTile(texture rl.Texture2D) int {
	ts.textures = append(ts.textures, texture)
	return len(ts.textures) - 1
}

type TileMap struct {
	tiles [][]int
}

func NewTileMap(tiles [][]int) TileMap {
	return TileMap{
		tiles: tiles,
	}
}
