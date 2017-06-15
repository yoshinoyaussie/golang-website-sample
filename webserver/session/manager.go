package session

import (
	"errors"
	"time"

	"github.com/labstack/echo"
)

// ID はセッションを一意に識別するIDです。
type ID string

// Store はセッションデータと整合性トークンを保持する構造体です。
type Store struct {
	Data             map[string]string
	ConsistencyToken string
}

// Manager は Sessionの操作・管理を行います。
type Manager struct {
	stopCh    chan struct{}
	commandCh chan command
	stopGCCh  chan struct{}
}

// Start は Managerの開始を行います。
func (m *Manager) Start(echo *echo.Echo) {
	e = echo
	go m.mainLoop()
	time.Sleep(100 * time.Millisecond)
	go m.gcLoop()
}

// Stop は Managerの停止を行います。
func (m *Manager) Stop() {
	m.stopGCCh <- struct{}{}
	time.Sleep(100 * time.Millisecond)
	m.stopCh <- struct{}{}
}

// Create は セッションの作成を行います。
func (m *Manager) Create() (ID, error) {
	respCh := make(chan response, 1)
	defer close(respCh)
	cmd := command{commandCreate, nil, respCh}
	m.commandCh <- cmd
	resp := <-respCh
	var res ID
	if resp.err != nil {
		e.Logger.Debugf("Session Create Error. [%s]", resp.err)
		return res, resp.err
	}
	if res, ok := resp.result[0].(ID); ok {
		return res, nil
	}
	e.Logger.Debugf("Session Create Error. [%s]", ErrorOther)
	return res, ErrorOther
}

// LoadStore は データストアの読み出しを行います。
func (m *Manager) LoadStore(sessionID ID) (Store, error) {
	respCh := make(chan response, 1)
	defer close(respCh)
	req := []interface{}{sessionID}
	cmd := command{commandLoadStore, req, respCh}
	m.commandCh <- cmd
	resp := <-respCh
	var res Store
	if resp.err != nil {
		e.Logger.Debugf("Session[%s] Load store Error. [%s]", sessionID, resp.err)
		return res, resp.err
	}
	if res, ok := resp.result[0].(Store); ok {
		return res, nil
	}
	e.Logger.Debugf("Session[%s] Load store Error. [%s]", sessionID, ErrorOther)
	return res, ErrorOther
}

// SaveStore は データストアの保存を行います。
func (m *Manager) SaveStore(sessionID ID, sessionStore Store) error {
	respCh := make(chan response, 1)
	defer close(respCh)
	req := []interface{}{sessionID, sessionStore}
	cmd := command{commandSaveStore, req, respCh}
	m.commandCh <- cmd
	resp := <-respCh
	if resp.err != nil {
		e.Logger.Debugf("Session[%s] Save store Error. [%s]", sessionID, resp.err)
		return resp.err
	}
	return nil
}

// Delete は セッションの削除を行います。
func (m *Manager) Delete(sessionID ID) error {
	respCh := make(chan response, 1)
	defer close(respCh)
	req := []interface{}{sessionID}
	cmd := command{commandDelete, req, respCh}
	m.commandCh <- cmd
	resp := <-respCh
	if resp.err != nil {
		e.Logger.Debugf("Session[%s] Delete Error. [%s]", sessionID, resp.err)
		return resp.err
	}
	return nil
}

// DeleteExpired は 期限切れセッションの削除を行います。
func (m *Manager) DeleteExpired() error {
	respCh := make(chan response, 1)
	defer close(respCh)
	cmd := command{commandDelete, nil, respCh}
	m.commandCh <- cmd
	resp := <-respCh
	if resp.err != nil {
		e.Logger.Debugf("Session DeleteExpired Error. [%s]", resp.err)
		return resp.err
	}
	return nil
}

// Managerが返す各エラーのインスタンスを生成します。
var (
	ErrorBadParameter   = errors.New("Bad Parameter")
	ErrorNotFound       = errors.New("Not Found")
	ErrorInvalidToken   = errors.New("Invalid Token")
	ErrorInvalidCommand = errors.New("Invalid Command")
	ErrorNotImplemented = errors.New("Not Implemented")
	ErrorOther          = errors.New("Other")
)
