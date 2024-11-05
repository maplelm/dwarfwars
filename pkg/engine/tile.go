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
