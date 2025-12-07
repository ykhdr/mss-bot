package minecraft

import (
	"context"
	"fmt"
	"time"

	"github.com/dreamscached/minequery/v2"
	"github.com/rs/zerolog/log"
)

// Client is a Minecraft server status client.
type Client struct {
	timeout time.Duration
}

// NewClient creates a new Minecraft client with the specified timeout.
func NewClient(timeout time.Duration) *Client {
	return &Client{
		timeout: timeout,
	}
}

// GetStatus queries the Minecraft server and returns its status.
func (c *Client) GetStatus(ctx context.Context, host string, port int) (*ServerStatus, error) {
	log.Debug().Str("host", host).Int("port", port).Dur("timeout", c.timeout).Msg("starting minecraft server query")

	pinger := minequery.NewPinger(
		minequery.WithTimeout(c.timeout),
		minequery.WithProtocolVersion17(minequery.Ping17ProtocolVersion119),
	)

	// Create a channel for the result
	type result struct {
		status *minequery.Status17
		err    error
	}
	resultCh := make(chan result, 1)

	go func() {
		status, err := pinger.Ping17(host, port)
		resultCh <- result{status: status, err: err}
	}()

	select {
	case <-ctx.Done():
		log.Warn().Str("host", host).Int("port", port).Msg("minecraft query canceled")
		return &ServerStatus{Online: false}, ctx.Err()
	case res := <-resultCh:
		if res.err != nil {
			log.Warn().Err(res.err).Str("host", host).Int("port", port).Msg("minecraft server ping failed")
			return &ServerStatus{Online: false}, nil
		}

		status := c.convertStatus(res.status)
		log.Debug().
			Str("host", host).
			Int("port", port).
			Bool("online", status.Online).
			Str("version", status.Version).
			Int("players", status.Players.Online).
			Msg("minecraft server query successful")
		return status, nil
	}
}

// convertStatus converts minequery status to our internal status format.
func (c *Client) convertStatus(status *minequery.Status17) *ServerStatus {
	if status == nil {
		return &ServerStatus{Online: false}
	}

	serverStatus := &ServerStatus{
		Online:      true,
		Version:     status.VersionName,
		Protocol:    status.ProtocolVersion,
		Description: status.Description.String(),
		Players: PlayersInfo{
			Online: status.OnlinePlayers,
			Max:    status.MaxPlayers,
			Sample: make([]Player, 0, len(status.SamplePlayers)),
		},
	}

	for _, p := range status.SamplePlayers {
		serverStatus.Players.Sample = append(serverStatus.Players.Sample, Player{
			Name: p.Nickname,
			UUID: p.UUID.String(),
		})
	}

	return serverStatus
}

// Ping checks if the server is reachable.
func (c *Client) Ping(ctx context.Context, host string, port int) (bool, error) {
	log.Debug().Str("host", host).Int("port", port).Msg("pinging minecraft server")
	status, err := c.GetStatus(ctx, host, port)
	if err != nil {
		log.Warn().Err(err).Str("host", host).Int("port", port).Msg("ping failed")
		return false, err
	}

	log.Debug().Str("host", host).Int("port", port).Bool("online", status.Online).Msg("ping result")
	return status.Online, nil
}

// FormatAddress formats host and port into a connection string.
func FormatAddress(host string, port int) string {
	if port == 25565 {
		return host
	}
	return fmt.Sprintf("%s:%d", host, port)
}

// ParseAddress parses a connection string into host and port.
func ParseAddress(address string) (host string, port int, err error) {
	log.Debug().Str("address", address).Msg("parsing minecraft address")

	port = 25565 // Default port

	// Check if port is specified
	var portStr string
	for i := len(address) - 1; i >= 0; i-- {
		if address[i] == ':' {
			host = address[:i]
			portStr = address[i+1:]
			break
		}
	}

	if host == "" {
		host = address
		log.Debug().Str("host", host).Int("port", port).Msg("address parsed (default port)")
		return host, port, nil
	}

	// Parse port
	if portStr != "" {
		_, err = fmt.Sscanf(portStr, "%d", &port)
		if err != nil {
			log.Error().Err(err).Str("address", address).Str("port_str", portStr).Msg("failed to parse port")
			return "", 0, fmt.Errorf("invalid port: %s", portStr)
		}
	}

	log.Debug().Str("host", host).Int("port", port).Msg("address parsed")
	return host, port, nil
}
