package service

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/ykhdr/mss-bot/internal/minecraft"
	"github.com/ykhdr/mss-bot/internal/storage"
	"github.com/ykhdr/mss-bot/internal/storage/models"
)

// ServerService provides business logic for server operations
type ServerService struct {
	storage storage.ServerStorage
	mc      *minecraft.Client
}

// NewServerService creates a new server service
func NewServerService(storage storage.ServerStorage, mc *minecraft.Client) *ServerService {
	return &ServerService{
		storage: storage,
		mc:      mc,
	}
}

// GetServerConfig returns the server configuration for a chat
func (s *ServerService) GetServerConfig(ctx context.Context, chatID int64) (*models.Server, error) {
	log.Debug().Int64("chat_id", chatID).Msg("getting server config")
	return s.storage.GetByChatID(ctx, chatID)
}

// SetServerConfig sets or updates the server configuration for a chat
func (s *ServerService) SetServerConfig(ctx context.Context, chatID int64, ip string, port int, name string) error {
	log.Info().Int64("chat_id", chatID).Str("ip", ip).Int("port", port).Str("name", name).Msg("setting server config")

	server := &models.Server{
		ChatID: chatID,
		IP:     ip,
		Port:   port,
		Name:   name,
	}

	return s.storage.Upsert(ctx, server)
}

// GetServerStatus returns the status of the configured server for a chat
func (s *ServerService) GetServerStatus(ctx context.Context, chatID int64) (*ServerStatusResult, error) {
	log.Debug().Int64("chat_id", chatID).Msg("getting server status")

	server, err := s.storage.GetByChatID(ctx, chatID)
	if err != nil {
		log.Warn().Err(err).Int64("chat_id", chatID).Msg("server config not found for status check")
		return nil, err
	}

	log.Debug().Int64("chat_id", chatID).Str("ip", server.IP).Int("port", server.Port).Msg("querying minecraft server")
	status, err := s.mc.GetStatus(ctx, server.IP, server.Port)
	if err != nil {
		log.Warn().Err(err).Int64("chat_id", chatID).Str("ip", server.IP).Int("port", server.Port).Msg("minecraft server query failed")
		return &ServerStatusResult{
			Server: server,
			Status: &minecraft.ServerStatus{Online: false},
			Error:  err,
		}, nil
	}

	log.Info().Int64("chat_id", chatID).Str("ip", server.IP).Int("port", server.Port).Bool("online", status.Online).Int("players", status.Players.Online).Msg("minecraft server status retrieved")
	return &ServerStatusResult{
		Server: server,
		Status: status,
	}, nil
}

// ServerStatusResult contains both server config and its current status
type ServerStatusResult struct {
	Server *models.Server
	Status *minecraft.ServerStatus
	Error  error
}

// FormatStatus formats the server status for display
func (r *ServerStatusResult) FormatStatus() string {
	if r.Server == nil {
		return "Сервер не настроен. Используйте настройки для добавления сервера."
	}

	serverName := r.Server.Name
	if serverName == "" {
		serverName = minecraft.FormatAddress(r.Server.IP, r.Server.Port)
	}

	if !r.Status.Online {
		return fmt.Sprintf("🔴 *%s*\n\n"+
			"Адрес: `%s`\n"+
			"Статус: Недоступен",
			escapeMarkdown(serverName),
			minecraft.FormatAddress(r.Server.IP, r.Server.Port),
		)
	}

	playersStr := ""
	if len(r.Status.Players.Sample) > 0 {
		playersStr = "\n\n👥 *Игроки онлайн:*\n"
		for _, p := range r.Status.Players.Sample {
			playersStr += fmt.Sprintf("• %s\n", escapeMarkdown(p.Name))
		}
	}

	return fmt.Sprintf("🟢 *%s*\n\n"+
		"Адрес: `%s`\n"+
		"Версия: %s\n"+
		"Онлайн: %d/%d%s",
		escapeMarkdown(serverName),
		minecraft.FormatAddress(r.Server.IP, r.Server.Port),
		escapeMarkdown(r.Status.Version),
		r.Status.Players.Online,
		r.Status.Players.Max,
		playersStr,
	)
}

// FormatConfig formats the server configuration for display
func FormatConfig(server *models.Server) string {
	if server == nil {
		return "⚙️ *Настройки сервера*\n\n" +
			"Сервер не настроен.\n\n" +
			"Для настройки отправьте команду:\n" +
			"`/mss-set <ip>:<port> <name>`\n\n" +
			"Пример:\n" +
			"`/mss-set mc.example.com:25565 My Server`"
	}

	serverName := server.Name
	if serverName == "" {
		serverName = "Не указано"
	}

	return fmt.Sprintf("⚙️ *Настройки сервера*\n\n"+
		"IP: `%s`\n"+
		"Порт: `%d`\n"+
		"Название: %s\n\n"+
		"Для изменения отправьте команду:\n"+
		"`/mss-set <ip>:<port> <name>`",
		server.IP,
		server.Port,
		escapeMarkdown(serverName),
	)
}

// escapeMarkdown escapes special Markdown characters
func escapeMarkdown(s string) string {
	replacer := []struct {
		old, new string
	}{
		{"_", "\\_"},
		{"*", "\\*"},
		{"[", "\\["},
		{"]", "\\]"},
		{"(", "\\("},
		{")", "\\)"},
		{"~", "\\~"},
		{"`", "\\`"},
		{">", "\\>"},
		{"#", "\\#"},
		{"+", "\\+"},
		{"-", "\\-"},
		{"=", "\\="},
		{"|", "\\|"},
		{"{", "\\{"},
		{"}", "\\}"},
		{".", "\\."},
		{"!", "\\!"},
	}

	result := s
	for _, r := range replacer {
		result = replaceAll(result, r.old, r.new)
	}
	return result
}

func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); i++ {
		if string(s[i]) == old {
			result += new
		} else {
			result += string(s[i])
		}
	}
	return result
}
