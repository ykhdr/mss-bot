package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateManager_SetAndGet(t *testing.T) {
	sm := NewStateManager()

	sm.SetState(12345, StateMainMenu, 100)

	assert.Equal(t, StateMainMenu, sm.GetState(12345))
	assert.Equal(t, 100, sm.GetMessageID(12345))
}

func TestStateManager_GetState_NotFound(t *testing.T) {
	sm := NewStateManager()

	assert.Equal(t, StateNone, sm.GetState(99999))
}

func TestStateManager_GetMessageID_NotFound(t *testing.T) {
	sm := NewStateManager()

	assert.Equal(t, 0, sm.GetMessageID(99999))
}

func TestStateManager_ClearState(t *testing.T) {
	sm := NewStateManager()

	sm.SetState(12345, StateSettings, 100)
	assert.Equal(t, StateSettings, sm.GetState(12345))

	sm.ClearState(12345)
	assert.Equal(t, StateNone, sm.GetState(12345))
}

func TestStateManager_IsInState(t *testing.T) {
	sm := NewStateManager()

	sm.SetState(12345, StateStatus, 100)

	assert.True(t, sm.IsInState(12345, StateStatus))
	assert.False(t, sm.IsInState(12345, StateMainMenu))
	assert.False(t, sm.IsInState(12345, StateSettings))
}

func TestStateManager_MultipleChatIDs(t *testing.T) {
	sm := NewStateManager()

	sm.SetState(111, StateMainMenu, 1)
	sm.SetState(222, StateStatus, 2)
	sm.SetState(333, StateSettings, 3)

	assert.Equal(t, StateMainMenu, sm.GetState(111))
	assert.Equal(t, StateStatus, sm.GetState(222))
	assert.Equal(t, StateSettings, sm.GetState(333))

	assert.Equal(t, 1, sm.GetMessageID(111))
	assert.Equal(t, 2, sm.GetMessageID(222))
	assert.Equal(t, 3, sm.GetMessageID(333))
}

func TestStateManager_UpdateState(t *testing.T) {
	sm := NewStateManager()

	sm.SetState(12345, StateMainMenu, 100)
	assert.Equal(t, StateMainMenu, sm.GetState(12345))

	sm.SetState(12345, StateSettings, 100)
	assert.Equal(t, StateSettings, sm.GetState(12345))
}
