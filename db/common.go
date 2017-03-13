package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

var Conn *gorm.DB

type Settings struct {
	Adapter  string `yaml:"adapter"`
	Database string `yaml:"database"`
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func InitConnection() {
	s := getDatabaseSettings()
	adapter := "mysql"
	port := "3306"
	charset := "utf8"
	connQuery := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		s.Username, s.Password, s.Host, port, s.Database, charset,
	)
	var err error
	Conn, err = gorm.Open(adapter, connQuery)
	Conn.LogMode(true) // FIXME: turn off in production

	if err != nil {
		panic(err)
	}
}

func getDatabaseSettings() Settings {
	path, _ := filepath.Abs("./db/database.yml")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var settings Settings
	err = yaml.Unmarshal(data, &settings)
	if err != nil {
		panic(err)
	}

	return settings
}
