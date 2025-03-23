package user

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/woxQAQ/gim/internal/wsgateway/base"
)

type MessageError struct {
	UserID     string
	PlatformID int32
	Err        error
}

func (e *MessageError) Error() string {
	return fmt.Errorf("send message to user %s platform %d failed: %v",
		e.UserID, e.PlatformID, e.Err).Error()
}

// BroadcastMessage 实现 IUserManager 接口.
func (um *Manager) BroadcastMessage(msg base.IMessage) []error {
	um.mutex.RLock()
	defer um.mutex.RUnlock()

	var errors []error
	for _, up := range um.users {
		up.mutex.RLock()
		for _, conn := range up.Conns {
			if conn.State() == base.Connected {
				if err := conn.Send(websocket.TextMessage, msg.GetPayload()); err != nil {
					errors = append(errors, &MessageError{
						UserID:     up.UserID,
						PlatformID: conn.PlatformID(),
						Err:        err,
					})
				}
			}
		}
		up.mutex.RUnlock()
	}
	return errors
}

// SendMessage 实现 IUserManager 接口.
func (um *Manager) SendMessage(userID string, msg base.IMessage) []error {
	um.mutex.RLock()

	var errors []error
	up, exists := um.users[userID]
	um.mutex.RUnlock()
	if !exists {
		return nil
	}

	up.mutex.RLock()
	defer up.mutex.RUnlock()
	for _, conn := range up.Conns {
		if conn.State() == base.Connected {
			if err := conn.Send(websocket.TextMessage, msg.GetPayload()); err != nil {
				errors = append(errors, &MessageError{
					UserID:     up.UserID,
					PlatformID: conn.PlatformID(),
					Err:        err,
				})
			}
		}
	}
	return errors
}

// SendPlatformMessage 实现 IUserManager 接口.
func (um *Manager) SendPlatformMessage(userID string, platformID int32, msg base.IMessage) error {
	conn, err := um.GetConn(userID, platformID)
	if err != nil || conn == nil {
		return err
	}
	return conn.Send(websocket.TextMessage, msg.GetPayload())
}

// GetAll 实现 IUserManager 接口.
func (um *Manager) GetAll() ([]*State, error) {
	um.mutex.RLock()
	defer um.mutex.RUnlock()

	userStates := make([]*State, 0, len(um.users))
	for userID := range um.users {
		state, err := um.GetState(userID)
		if err != nil {
			continue
		}
		userStates = append(userStates, state)
	}
	return userStates, nil
}
