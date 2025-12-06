package models

import "time"

// Server represents a Minecraft server configuration for a chat
type Server struct {
	ID        int64
	ChatID    int64
	IP        string
	Port      int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// DefaultPort is the default Minecraft server port
const DefaultPort = 25565

// Address returns the server address in ip:port format
func (s *Server) Address() string {
	if s.Port == DefaultPort {
		return s.IP
	}
	return s.IP + ":" + string(rune(s.Port))
}
