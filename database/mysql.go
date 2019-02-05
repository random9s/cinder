package database

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/random9s/cinder/database/config"
	db "github.com/random9s/cinder/database/marshal"
)

//Mysql ...
type Mysql struct {
	conf *config.MySQL
	*db.MySQL
}

//Register registers mysql connection
func (m *Mysql) Register(c config.Config) {
	m.conf = c.(*config.MySQL)
}

//Open opens mysql connection
func (m *Mysql) Open() error {
	//Open database
	dbconn, err := sql.Open("mysql", m.FormatDSN())
	if err != nil {
		return err
	}

	//check connection
	err = dbconn.Ping()
	if err != nil {
		return err
	}

	var wrapped = &Mysql{
		m.conf,
		db.NewMySQL(dbconn),
	}

	//wrap and return connection
	*m = *wrapped
	return err
}

//FormatDSN ...
func (m *Mysql) FormatDSN() string {
	var loc = new(time.Location)
	var err error
	if m.conf.Loc != "" {
		loc, err = time.LoadLocation(m.conf.Loc)
		if err != nil {
			panic(err)
		}
	}

	var conf = &mysql.Config{
		User:                    m.conf.User,
		Passwd:                  m.conf.Passwd,
		Net:                     m.conf.Net,
		Addr:                    m.conf.Addr,
		DBName:                  m.conf.DBName,
		Params:                  m.conf.Params,
		Collation:               m.conf.Collation,
		Loc:                     loc,
		MaxAllowedPacket:        m.conf.MaxAllowedPacket,
		TLSConfig:               m.conf.TLSConfig,
		AllowAllFiles:           m.conf.AllowAllFiles,
		AllowCleartextPasswords: m.conf.AllowCleartextPasswords,
		AllowNativePasswords:    m.conf.AllowNativePasswords,
		AllowOldPasswords:       m.conf.AllowOldPasswords,
		ClientFoundRows:         m.conf.ClientFoundRows,
		ColumnsWithAlias:        m.conf.ColumnsWithAlias,
		InterpolateParams:       m.conf.InterpolateParams,
		MultiStatements:         m.conf.MultiStatements,
		ParseTime:               m.conf.ParseTime,
		RejectReadOnly:          m.conf.RejectReadOnly,
	}
	return conf.FormatDSN()
}
