package model

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	"github.com/labstack/echo"
)

// User はユーザーの情報を表します。
type User struct {
	ID       ID        `json:"id"`
	UserID   string    `json:"user_id"`
	Password StringMD5 `json:"password"`
	FullName string    `json:"full_name"`
	Roles    []Role    `json:"roles"`
}

// Copy は情報のコピーを行います。
func (u *User) Copy(f *User) {
	u.ID = f.ID
	u.UserID = f.UserID
	u.Password = f.Password
	u.FullName = f.FullName
	u.Roles = make([]Role, len(f.Roles))
	copy(u.Roles, f.Roles)
}

// UserDataAccessor はユーザーの情報を操作するAPIを提供します。
type UserDataAccessor struct {
	stopCh    chan struct{}
	commandCh chan command
}

// ID は情報を一意に識別するためのIDです。
type ID string

// StringMD5 はMD5ハッシュ化された文字列です。
type StringMD5 string

// Role はユーザーの権限を表します。
type Role string

// ユーザー権限の定義
const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// Start はAccessorの開始を行います。
func (a *UserDataAccessor) Start(echo *echo.Echo) error {
	e = echo
	users = make(map[ID]User)
	if err := a.decodeJSON(); err != nil {
		return err
	}
	go a.mainLoop()
	return nil
}

// Stop はAccessorの停止を行います。
func (a *UserDataAccessor) Stop() {
	a.stopCh <- struct{}{}
}

// FindAll はユーザーを全件検索します。
func (a *UserDataAccessor) FindAll() ([]User, error) {
	respCh := make(chan response, 1)
	defer close(respCh)
	req := []interface{}{}
	cmd := command{commandFindAll, req, respCh}
	a.commandCh <- cmd
	resp := <-respCh
	var res []User
	if resp.err != nil {
		e.Logger.Debugf("User Find Error. [%s]", resp.err)
		return res, resp.err
	}
	if res, ok := resp.result[0].([]User); ok {
		return res, nil
	}
	e.Logger.Debugf("User Find Error. [%s]", ErrorOther)
	return res, ErrorOther
}

// FindByUserID はUserIDでユーザーを検索します。
func (a *UserDataAccessor) FindByUserID(reqUserID string, option FindOption) ([]User, error) {
	respCh := make(chan response, 1)
	defer close(respCh)
	req := []interface{}{reqUserID, option}
	cmd := command{commandFindByUserID, req, respCh}
	a.commandCh <- cmd
	resp := <-respCh
	var res []User
	if resp.err != nil {
		e.Logger.Debugf("User[UserID=%s] Find Error. [%s]", reqUserID, resp.err)
		return res, resp.err
	}
	if res, ok := resp.result[0].([]User); ok {
		return res, nil
	}
	e.Logger.Debugf("User[UserID=%s] Find Error. [%s]", reqUserID, ErrorOther)
	return res, ErrorOther
}

// EncodeStringMD5 は、MD5エンコードした文字列を返します。
func EncodeStringMD5(str string) StringMD5 {
	h := md5.New()
	io.WriteString(h, str)
	encodeStr := hex.EncodeToString(h.Sum(nil))
	res := StringMD5(encodeStr)

	return res
}

// FindOption は検索時のオプションを定義します。
type FindOption int

// 検索時のオプション
const (
	FIndAll    FindOption = iota // 全件検索
	FindFirst                    // 1件目のみ返す
	FindUnique                   // 結果が1件のみでない場合にはエラーを返す
)

// DataAccessorが返す各エラーのインスタンスを生成します。
var (
	ErrorNotFound        = errors.New("Not found")
	ErrorMultipleResults = errors.New("Multiple results")
	ErrorInvalidCommand  = errors.New("Invalid Command")
	ErrorBadParameter    = errors.New("Bad Parameter")
	ErrorNotImplemented  = errors.New("Not Implemented")
	ErrorOther           = errors.New("Other")
)

func (a *UserDataAccessor) decodeJSON() error {
	// JSONファイル読み込み
	bytes, err := ioutil.ReadFile("data/users.json")
	if err != nil {
		return err
	}
	// JSONをデコードする
	var records []User
	if err := json.Unmarshal(bytes, &records); err != nil {
		return err
	}
	// 結果をmapにセットする
	for _, x := range records {
		users[x.ID] = x
	}
	return nil
}

// echoのインスタンス
var e *echo.Echo

// 情報をメモリ上に持つためのmap
var users map[ID]User

// コマンド種別の定義
type commandType int

const (
	commandFindAll      commandType = iota // 全件検索
	commandFindByID                        // IDで検索
	commandFindByUserID                    // UserIDで検索
)

// コマンド実行のためのパラメータ
type command struct {
	cmdType    commandType
	req        []interface{}
	responseCh chan response
}

// コマンド実行の結果
type response struct {
	result []interface{}
	err    error
}

// UserDataAccessor のメインループ処理
func (a *UserDataAccessor) mainLoop() {
	a.stopCh = make(chan struct{}, 1)
	a.commandCh = make(chan command, 1)
	defer close(a.commandCh)
	defer close(a.stopCh)
	e.Logger.Info("model.UserDataAccessor:start")
loop:
	for {
		// 受信したコマンドによって処理を振り分ける
		select {
		case cmd := <-a.commandCh:
			switch cmd.cmdType {
			// 全件検索
			case commandFindAll:
				results := []User{}
				for _, x := range users {
					user := User{}
					user.Copy(&x)
					results = append(results, user)
				}
				res := []interface{}{results}
				cmd.responseCh <- response{res, nil}
				break
			// IDで検索
			case commandFindByID:
				// 未実装
				cmd.responseCh <- response{nil, ErrorNotImplemented}
				break
			// UserIDで検索
			case commandFindByUserID:
				reqUserID, ok := cmd.req[0].(string)
				if !ok {
					cmd.responseCh <- response{nil, ErrorBadParameter}
					break
				}
				reqOption, ok := cmd.req[1].(FindOption)
				if !ok {
					cmd.responseCh <- response{nil, ErrorBadParameter}
					break
				}
				results := []User{}
				for _, x := range users {
					if x.UserID == reqUserID {
						user := User{}
						user.Copy(&x)
						results = append(results, user)
						if reqOption == FindFirst {
							break
						}
					}
				}
				if len(results) <= 0 {
					cmd.responseCh <- response{nil, ErrorNotFound}
					break
				}
				if reqOption == FindUnique && len(results) > 1 {
					cmd.responseCh <- response{nil, ErrorMultipleResults}
					break
				}
				res := []interface{}{results}
				cmd.responseCh <- response{res, nil}
			// それ以外（エラー）
			default:
				cmd.responseCh <- response{nil, ErrorInvalidCommand}
			}
		case <-a.stopCh:
			break loop
		}
	}
	e.Logger.Info("model.UserDataAccessor:stop")
}
