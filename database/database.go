package database

import (
	"github.com/random9s/cinder/database/config"
	"github.com/random9s/cinder/database/marshal"
)

//Database ...
type Database interface {
	Register(config.Config)
	Open() (db.UnmarshalMarshaler, error)
	Close() error
}

//Register registers new database
func Register(c config.Config, d Database) {
	d.Register(c)
}
