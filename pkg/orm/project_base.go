package orm

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"strings"
)

type TableMeta struct {
	Field   string `json:"field"`
	Type    string `json:"type"`
	Null    string `json:"null"`
	Key     string `json:"key"`
	Default string `json:"default"`
	Extra   string `json:"extra"`
}

type TableField struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Null    string `json:"null"`
	Key     string `json:"key"`
	Default string `json:"default"`
	Extra   string `json:"extra"`
}

type Table struct {
	Id          int          `json:"id"`
	Name        string       `json:"name"`
	TableFields []TableField `json:"children"`
}

type ProjectDatabase struct {
	projectDbConnections map[string]*gorm.DB
	referenceDatabase    *gorm.DB
}

func (pb *ProjectDatabase) getConnection(connId string) *gorm.DB {
	return pb.projectDbConnections[connId]
}

func (pb *ProjectDatabase) ExecRawSql(connId, rawSql string) error {
	err := pb.getConnection(connId).Exec(rawSql).Error
	return err
}
func (pb *ProjectDatabase) TruncateTable(projectId string, tableName string) error {
	err := pb.getConnection(projectId).Exec("TRUNCATE TABLE " + tableName).Error
	return err
}

func (pb *ProjectDatabase) DropTable(projectId string, tableName string) error {
	err := pb.getConnection(projectId).Exec("DROP TABLE " + tableName).Error
	return err
}

func (pb *ProjectDatabase) CreateProjectDb(projectId string, name string, dbConfig DatabaseConfig) error {
	err := pb.referenceDatabase.Exec("CREATE DATABASE " + name).Error
	if err != nil {
		return err
	}
	// now connect db and add to map
	dbConnection, err := dbConfig.GetDbConnectionFromDbName(name)
	if err != nil {
		pb.referenceDatabase.Exec("DROP DATABASE " + name)
		return err
	}

	pb.projectDbConnections[projectId] = dbConnection
	return nil
}

//func (pb *ProjectBase) NoOfRecords(table string) int64 {
//	var numberOfRecords int64
//	pb.baseDb.Raw("SELECT COUNT(*) as numberOfRecords FROM " + table).Row().Scan(&numberOfRecords)
//
//	return numberOfRecords
//}

func (pb *ProjectDatabase) LoadDataFromFile(connId string, filename string, loadSql string) error {
	// create condition table
	loadSql = strings.Replace(loadSql, "%file_name%", filename, -1)
	fmt.Println("loading sql file :", loadSql)
	mysql.RegisterLocalFile(filename)
	err := pb.getConnection(connId).Exec(loadSql).Error
	return err

}

func (pb *ProjectDatabase) GetTableMetaInterface(connId string, name string) (*[]interface{}, error) {
	var dbObjects []TableMeta
	err := pb.getConnection(connId).Raw("DESCRIBE " + name).Scan(&dbObjects).Error
	if err == nil {
		interfaceObjects := make([]interface{}, len(dbObjects))
		for i, v := range dbObjects {
			interfaceObjects[i] = v
		}
		return &interfaceObjects, err
	}
	return nil, err
}

func (pb *ProjectDatabase) GetTableMeta(connId string, name string) (*[]TableMeta, error) {
	var dbObjects []TableMeta
	err := pb.getConnection(connId).Raw("DESCRIBE " + name).Scan(&dbObjects).Error
	if err == nil {

		return &dbObjects, err
	}
	return nil, err
}
