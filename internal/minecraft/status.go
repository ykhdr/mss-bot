package minecraft

// ServerStatus represents the status of a Minecraft server
type ServerStatus struct {
	Online      bool
	Version     string
	Protocol    int
	Players     PlayersInfo
	Description string
}

// PlayersInfo contains player count and list information
type PlayersInfo struct {
	Online int
	Max    int
	Sample []Player
}

// Player represents a player on the server
type Player struct {
	Name string
	UUID string
}
