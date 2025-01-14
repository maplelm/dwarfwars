package engine

import (
	"library/command"
)

/*
This struct will contain all of the relevent data to a single Game world instance
*/
type World struct {
	Tiles           []Tile
	CmdQueue        chan *command.Command
	updateCmdBuffer []*command.Command
}

func NewWorld(qs int) *World {
	return &World{
		Tiles:           nil,
		CmdQueue:        make(chan *command.Command, qs),
		updateCmdBuffer: make([]*command.Command, qs),
	}
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

func (w *World) Update() {
	s := len(w.CmdQueue) // number of commands to process this update
	for i := 0; i < s; i++ {
		w.updateCmdBuffer[i] = <-w.CmdQueue // pulling the commands from channel
	}
	// update the world
}
