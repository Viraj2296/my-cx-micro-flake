package repository

import (
	"go.cerex.io/transcendflow/component"
	"go.cerex.io/transcendflow/logging"
	"go.cerex.io/transcendflow/orm"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

/*
Repository Layer
The repository layer interacts with the database or any external data source.
Place your repository interfaces and implementations in the internal/repository directory.
*/
type Repository interface {
	// GetResource define the functions here, then we can have own implementation
	GetResource(tableName string, resourceId int) (error, component.GeneralObject)
	UpdateResource(tableName string, resourceId int, serialisedData datatypes.JSON, updatedBy int) error
	GetResourceWithCondition(tableName string, condition string) (error, []component.GeneralObject)
	GetResources(tableName string) (error, []component.GeneralObject)
	GetCountByCondition(tableName string, condition string) int64
	CreateResource(tableName string, serialisedData datatypes.JSON, createdBy int) (error, int)
}
type sparePartInventoryRepository struct {
	Database *gorm.DB
	Logger   *logging.Logger
}

func (v sparePartInventoryRepository) CreateResource(tableName string, serialisedData datatypes.JSON, createdBy int) (error, int) {
	err, c := orm.CreateFromResource(v.Database, tableName, createdBy, serialisedData)
	return err, c
}
func (v sparePartInventoryRepository) GetResource(tableName string, resourceId int) (error, component.GeneralObject) {
	err, c := orm.Get(v.Database, tableName, resourceId)
	return err, c
}

func (v sparePartInventoryRepository) GetResources(tableName string) (error, []component.GeneralObject) {
	err, c := orm.GetObjects(v.Database, tableName)
	return err, c
}

func (v sparePartInventoryRepository) GetResourceWithCondition(tableName string, condition string) (error, []component.GeneralObject) {
	err, c := orm.GetConditionalObjects(v.Database, tableName, condition)
	return err, c
}

func (v sparePartInventoryRepository) UpdateResource(tableName string, resourceId int, serialisedData datatypes.JSON, updatedBy int) error {
	err := orm.UpdateSerialisedResourceFromId(v.Database, tableName, resourceId, updatedBy, serialisedData)
	return err
}
func NewSparePartInventoryRepository(database *gorm.DB, logger *logging.Logger) Repository {
	return &sparePartInventoryRepository{
		Database: database,
		Logger:   logger,
	}
}
func (v sparePartInventoryRepository) GetCountByCondition(tableName string, condition string) int64 {
	count := orm.CountByCondition(v.Database, tableName, condition)
	return count
}
