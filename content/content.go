package content

import (
	"github.com/fc-testzone/apiserver/db"
	"github.com/fc-testzone/apiserver/utils"
)

type Post struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Content struct {
	cfg *utils.Configs
}

func NewContent(c *utils.Configs) *Content {
	return &Content{
		cfg: c,
	}
}

func (c *Content) Posts(posts *[]Post) error {
	var dbCfg = c.cfg.Settings().DB

	// Connect to DB
	var db = db.NewDatabase()
	var err = db.Connect(dbCfg.IP, dbCfg.Port, dbCfg.User, dbCfg.Passwd, dbCfg.DB)
	if err != nil {
		return err
	}

	// Find all posts
	err = db.Find(&Post{}, &Post{}, posts)
	if err != nil {
		return err
	}

	return nil
}
