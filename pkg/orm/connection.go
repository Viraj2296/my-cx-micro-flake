package orm

import (
	"cx-micro-flake/pkg/util/encryption"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DatabaseConfig struct {
	Host                   string `json:"host"`
	Port                   int    `json:"port"`
	Type                   string `json:"type"`
	Name                   string `json:"name"`
	User                   string `json:"user"`
	Password               string `json:"password"`
	SecureConnectionString string `json:"secureConnectionString"`
	IsSingularTable        bool   `json:"isSingularTable"`
}

func (dc *DatabaseConfig) getConnectionFromSecureString() (*gorm.DB, error) {

	dbParam, err := encryption.FnDecrypt(dc.SecureConnectionString)

	if err != nil {
		return nil, err
	}

	if dc.Type == "mysql" {
		if dc.IsSingularTable {
			dbConnection, err := gorm.Open(mysql.Open(dbParam), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
			return dbConnection, err
		} else {
			dbConnection, err := gorm.Open(mysql.Open(dbParam), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: false}})
			return dbConnection, err
		}

	} else {
		return nil, errors.New("unsupported database type")
	}
}

func (dc *DatabaseConfig) NewConnection() (*gorm.DB, error) {
	if dc.SecureConnectionString != "" {
		return dc.getConnectionFromSecureString()
	}
	if dc.Type == "mysql" {
		dbURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", dc.User, dc.Password, dc.Host, dc.Port, dc.Name)
		if dc.IsSingularTable {
			dbConnection, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
			return dbConnection, err
		} else {
			dbConnection, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
			return dbConnection, err
		}

	} else {
		return nil, errors.New("unsupported database type")
	}
}

func (dc *DatabaseConfig) GetDbConnectionFromDbName(dbName string) (*gorm.DB, error) {

	if dc.Type == "mysql" {
		dbURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", dc.User, dc.Password, dc.Host, dc.Port, dbName)
		if dc.IsSingularTable {
			dbConnection, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
			return dbConnection, err
		} else {
			dbConnection, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
			return dbConnection, err
		}

	} else {
		return nil, errors.New("unsupported database type")
	}
}
