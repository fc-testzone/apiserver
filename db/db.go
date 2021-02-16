package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) Connect(ip string, port int, user string, passwd string, db string) error {
	var err error

	var dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, passwd, ip, port, db)
	d.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return err
}

func (d *Database) Insert(table interface{}, data interface{}) error {
	return d.db.Model(table).Create(data).Error
}

func (d *Database) Update(table interface{}, fieldCon string, condition interface{}, field string, newData interface{}) error {
	return d.db.Model(table).Where(fieldCon+" = ?", condition).Update(field, newData).Error
}

func (d *Database) Find(table interface{}, fieldCon string, condition interface{}, out interface{}) error {
	return d.db.Model(table).Where(fieldCon+" = ?", condition).Find(&out).Error
}

func (d *Database) FindAll(table interface{}, out interface{}) error {
	return d.db.Model(table).Find(out).Error
}
