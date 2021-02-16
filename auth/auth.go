package auth

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"

	"github.com/fc-testzone/apiserver/db"
	"github.com/fc-testzone/apiserver/utils"
)

type User struct {
	Login  string
	Passwd string
	Token  string
}

type Authorizer struct {
	cfg *utils.Configs
}

func NewAuthorizer(c *utils.Configs) *Authorizer {
	return &Authorizer{
		cfg: c,
	}
}

func (a *Authorizer) CreateUser(login string, passwd string) error {
	var dbCfg = a.cfg.Settings().DB

	var db = db.NewDatabase()
	var err = db.Connect(dbCfg.IP, dbCfg.Port, dbCfg.User, dbCfg.Passwd, dbCfg.DB)
	if err != nil {
		return err
	}

	return db.Insert(&User{}, &User{
		Login:  login,
		Passwd: passwd,
	})
}

func (a *Authorizer) CreateToken(login string, passwd string) (string, error) {
	var dbCfg = a.cfg.Settings().DB

	// Connect to DB
	var db = db.NewDatabase()
	var err = db.Connect(dbCfg.IP, dbCfg.Port, dbCfg.User, dbCfg.Passwd, dbCfg.DB)
	if err != nil {
		return "", err
	}

	// Generate token
	var sha = sha1.New()
	sha.Write([]byte(login + passwd))
	var token = hex.EncodeToString(sha.Sum(nil))

	// Update token in DB
	return token, db.Update(&User{}, "login", login, "token", token)
}

func (a *Authorizer) CheckToken(token string) error {
	var dbCfg = a.cfg.Settings().DB
	var usr []*User

	// Connect to DB
	var db = db.NewDatabase()
	var err = db.Connect(dbCfg.IP, dbCfg.Port, dbCfg.User, dbCfg.Passwd, dbCfg.DB)
	if err != nil {
		return err
	}

	// Check find token in db
	err = db.Find(&User{}, "token", token, usr)
	if err != nil {
		return err
	}

	if len(usr) == 0 {
		return errors.New("User not found")
	}

	return nil
}
