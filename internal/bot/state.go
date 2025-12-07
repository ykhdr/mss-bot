package bot

import "sync"

// State represents the current state of the bot for a specific chat
type State int

const (
	// StateNone - no active interaction
	StateNone State = iota
	// StateMainMenu - main menu is displayed
	StateMainMenu
	// StateStatus - server status is displayed
	StateStatus
	// StatePlayers - players list is displayed
	StatePlayers
	// StateSettings - settings menu is displayed
	StateSettings
)

// StateManager manages bot states for different chats
type StateManager struct {
	mu     sync.RWMutex
	states map[int64]chatState
}

type chatState struct {
	state     State
	messageID int
}

// NewStateManager creates a new state manager
func NewStateManager() *StateManager {
	return &StateManager{
		states: make(map[int64]chatState),
	}
}

// GetState returns the current state for a chat
func (sm *StateManager) GetState(chatID int64) State {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if cs, ok := sm.states[chatID]; ok {
		return cs.state
	}
	return StateNone
}

// SetState sets the state for a chat
func (sm *StateManager) SetState(chatID int64, state State, messageID int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.states[chatID] = chatState{
		state:     state,
		messageID: messageID,
	}
}

// GetMessageID returns the message ID for a chat's current state
func (sm *StateManager) GetMessageID(chatID int64) int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if cs, ok := sm.states[chatID]; ok {
		return cs.messageID
	}
	return 0
}

// ClearState removes the state for a chat
func (sm *StateManager) ClearState(chatID int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.states, chatID)
}

// IsInState checks if a chat is in a specific state
func (sm *StateManager) IsInState(chatID int64, state State) bool {
	return sm.GetState(chatID) == state
}
