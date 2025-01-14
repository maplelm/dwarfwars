package lobby

import "server/internal/client"

type Lobby struct {
	Clients []*client.Client
}
