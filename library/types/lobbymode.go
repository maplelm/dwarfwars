package types

type LobbyMode byte

func (l *LobbyMode) IsPublic() bool {
	return (*l & 0b001) > 0
}

func (l *LobbyMode) IsPrivate() bool {
	return (*l & 0b010) > 0
}

func (l *LobbyMode) IsPasswordProtected() bool {
	return (*l & 0b100) > 0
}

func (l *LobbyMode) SetPublic() {
	*l = (*l & 0b100) + 0b001
}

func (l *LobbyMode) SetPrivate() {
	*l = (*l & 0b100) + 0b010
}

func (l *LobbyMode) SetPasswordProtected(t bool) {
	if t {
		*l = *l | 0b100
	} else {
		*l = *l & 0b011
	}
}
