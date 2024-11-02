package engine

type EntityHandler interface {
	ID() int
	Health() float32
	Passable() bool
	OnDeath()
}

const (
	TileTypeDirt = iota
	TileTypeAir
	TileTypeGrass
	TileTypeSand
	TileTypeWater
)

type Tile struct {
	Entities []EntityHandler
	TileType int
}

type World struct {
	Tiles []Tile
}

func WorldGen(x, y, z, seed int) *World {

	w := &World{
		Tiles: make([]Tile, x*y*z),
	}

	for i, _ := range w.Tiles {
		if i > ((x * y) * z / 2) {
			w.Tiles[i].TileType = TileTypeAir
		} else {
			w.Tiles[i].TileType = TileTypeDirt
		}
	}

	return w
}
