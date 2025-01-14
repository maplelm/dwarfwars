package types

type World struct {
	Chunks []Chunk
}

type Chunk struct {
	Tiles []Tile
}

type Tile struct {
	Type     TileType
	Contains []Entity
}

type TileType uint16
