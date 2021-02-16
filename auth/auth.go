package auth

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"time"

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

	// Hash password
	var sha = sha1.New()
	sha.Write([]byte(passwd))
	var hPasswd = hex.EncodeToString(sha.Sum(nil))

	return db.Insert(&User{}, &User{
		Login:  login,
		Passwd: hPasswd,
	})
}

func (a *Authorizer) CheckUser(login string, passwd string) error {
	var dbCfg = a.cfg.Settings().DB
	var usr []*User

	// Connect to DB
	var db = db.NewDatabase()
	var err = db.Connect(dbCfg.IP, dbCfg.Port, dbCfg.User, dbCfg.Passwd, dbCfg.DB)
	if err != nil {
		return err
	}

	// Hash password
	var sha = sha1.New()
	sha.Write([]byte(passwd))
	var hPasswd = hex.EncodeToString(sha.Sum(nil))

	// Find user in db
	err = db.Find(&User{}, &User{Login: login, Passwd: hPasswd}, &usr)
	if err != nil {
		return err
	}

	if len(usr) == 0 {
		return errors.New("User not found")
	}

	return nil
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
	sha.Write([]byte(login + passwd + time.Now().Format("20060201")))
	var token = hex.EncodeToString(sha.Sum(nil))

	// Update token in DB
	return token, db.Update(&User{}, "login", login, "token", token)
}

func (a *Authorizer) CheckToken(token string) error {
	var dbCfg = a.cfg.Settings().DB
	var usr []User

	// Connect to DB
	var db = db.NewDatabase()
	var err = db.Connect(dbCfg.IP, dbCfg.Port, dbCfg.User, dbCfg.Passwd, dbCfg.DB)
	if err != nil {
		return err
	}

	// Find token in db
	err = db.Find(&User{}, &User{Token: token}, &usr)
	if err != nil {
		return err
	}

	if len(usr) == 0 {
		return errors.New("User not found")
	}

	return nil
}
