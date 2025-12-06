package minecraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAddress_HostOnly(t *testing.T) {
	host, port, err := ParseAddress("mc.example.com")

	assert.NoError(t, err)
	assert.Equal(t, "mc.example.com", host)
	assert.Equal(t, 25565, port)
}

func TestParseAddress_HostAndPort(t *testing.T) {
	host, port, err := ParseAddress("mc.example.com:25566")

	assert.NoError(t, err)
	assert.Equal(t, "mc.example.com", host)
	assert.Equal(t, 25566, port)
}

func TestParseAddress_CustomPort(t *testing.T) {
	host, port, err := ParseAddress("play.server.net:19132")

	assert.NoError(t, err)
	assert.Equal(t, "play.server.net", host)
	assert.Equal(t, 19132, port)
}

func TestParseAddress_IPAddress(t *testing.T) {
	host, port, err := ParseAddress("192.168.1.100:25565")

	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.100", host)
	assert.Equal(t, 25565, port)
}

func TestParseAddress_IPAddressOnly(t *testing.T) {
	host, port, err := ParseAddress("192.168.1.100")

	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.100", host)
	assert.Equal(t, 25565, port)
}

func TestFormatAddress_DefaultPort(t *testing.T) {
	addr := FormatAddress("mc.example.com", 25565)
	assert.Equal(t, "mc.example.com", addr)
}

func TestFormatAddress_CustomPort(t *testing.T) {
	addr := FormatAddress("mc.example.com", 25566)
	assert.Equal(t, "mc.example.com:25566", addr)
}

func TestNewClient(t *testing.T) {
	client := NewClient(5000000000) // 5 seconds in nanoseconds

	assert.NotNil(t, client)
	assert.Equal(t, 5000000000, int(client.timeout))
}

func TestConvertStatus_Nil(t *testing.T) {
	client := NewClient(5000000000)

	status := client.convertStatus(nil)

	assert.False(t, status.Online)
}


func TestServerStatus_Fields(t *testing.T) {
	status := ServerStatus{
		Online:      true,
		Version:     "1.20.4",
		Protocol:    765,
		Description: "Test Server",
		Players: PlayersInfo{
			Online: 5,
			Max:    20,
			Sample: []Player{
				{Name: "Player1", UUID: "uuid-1"},
				{Name: "Player2", UUID: "uuid-2"},
			},
		},
	}

	assert.True(t, status.Online)
	assert.Equal(t, "1.20.4", status.Version)
	assert.Equal(t, 5, status.Players.Online)
	assert.Equal(t, 20, status.Players.Max)
	assert.Len(t, status.Players.Sample, 2)
	assert.Equal(t, "Player1", status.Players.Sample[0].Name)
}
