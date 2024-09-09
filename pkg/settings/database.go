package settings

import (
	"fmt"
	"github.com/maplelm/dwarfwars/pkg/logging"
)

const (
	DbDriverMySQL  = "mysql"
	DbDriverMsSQL  = "mssql"
	DbDriverSQLite = "sqlite3"
)

type Database struct {
	Addr      string                 `toml:"Address"`
	Driver    string                 `toml:"Driver"`
	Trusted   bool                   `toml:"Trusted"`
	Port      int                    `toml:"Port"`
	Username  string                 `toml:"Username"`
	Password  string                 `toml:"Password"`
	DefaultDB string                 `toml:"DB"`
	Views     map[string][]SQLColumn `toml:"Views"`
	Tables    map[string][]SQLColumn `toml:"Tables"`
}

func (d Database) ConnString() (str string, err error) {
	switch d.Driver {
	case DbDriverMySQL:
		str = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", d.Username, d.Password, d.Addr, d.Port, d.DefaultDB)
	case DbDriverMsSQL:
		if !d.Trusted {
			str = fmt.Sprintf("Server=%s:%s; Database=%s; User Id=%s; Password=%s;", d.Addr, d.Port, d.DefaultDB, d.Username, d.Password)
		} else {
			str = fmt.Sprintf("Server=%s:%s; Database=%s; Trusted_Connection=True", d.Addr, d.Port, d.DefaultDB)
		}
	case DbDriverSQLite:
		logging.Info("SQLite driver support not implemented yet")
	default:
		err = fmt.Errorf("Unsupported Driver: %s", d.Driver)
	}
	return
}

func (d Database) ViewKeys() (l []string) {
	for k, _ := range d.Views {
		l = append(l, k)
	}
	return
}

func (d Database) TableKeys() (l []string) {
	for k, _ := range d.Tables {
		l = append(l, k)
	}
	return
}

type SQLColumn struct {
	Name          string  `toml:"Name"`
	ColType       string  `toml:"Type"`
	AutoIncrement bool    `toml:"Auto Increment"`
	PrimaryKey    bool    `toml:"Primary Key"`
	ForiegnKey    bool    `toml:"Foriegn Key"`
	Nullable      bool    `toml:"Null Allowed"`
	DefaultValue  *string `toml:"Default"`
}
