package lobby

import "github.com/maplelm/dwarfwars/cmd/server/pkg/client"

type Lobby struct {
	Clients []*client.Client
}
