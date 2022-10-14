package gorm_masking

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
)

type TestTable struct {
	Id     int    `gorm:"primary_key;column:id"`
	Name   string `gorm:"column:name" masking:"common" masking_encrypt_column:"name_id"`
	Age    int    `gorm:"column:age"`
	NameId string `gorm:"column:name_id"`
}

func (r TestTable) TableName() string {
	return "t_test_table"
}

var db *gorm.DB
var encryptKey = "1234567890123456"

func init() {
	dsn := "file:test?mode=memory&cache=shared"
	sqliteDb, _ := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	sqliteDb.Use(New(&Config{Key: encryptKey}))
	sqliteDb.AutoMigrate(&TestTable{})
	db = sqliteDb

}

func Test_BaseUse(t *testing.T) {
	addData := &TestTable{
		Id:   1,
		Name: "base use",
		Age:  18,
	}
	if err := db.Create(addData).Error; err != nil {
		t.Error(err)
	}
	assert.Equal(t, "ba****se", addData.Name)

	queryCtx := context.WithValue(context.Background(), UnMaskingContext, []string{"Name"})
	var queryDataByFind TestTable
	if err := db.WithContext(queryCtx).Find(&queryDataByFind).Error; err != nil {
		t.Error(err)
	}
	assert.Equal(t, "base use", queryDataByFind.Name)
	var queryDataByFirst TestTable
	if err := db.WithContext(queryCtx).First(&queryDataByFirst).Error; err != nil {
		t.Error(err)
	}
	assert.Equal(t, "base use", queryDataByFirst.Name)

	var queryDataByRaw TestTable
	if err := db.WithContext(queryCtx).
		Raw("select * from t_test_table").Find(&queryDataByRaw).Error; err != nil {
		t.Error(err)
	}
	assert.Equal(t, "base use", queryDataByRaw.Name)

	saveValue := "test save"
	addData.Name = saveValue
	if err := db.Save(addData).Error; err != nil {
		t.Error(err)
	}
	var saveResult TestTable
	if err := db.WithContext(queryCtx).First(&saveResult).Error; err != nil {
		t.Error(err)
	}
	assert.Equal(t, saveValue, saveResult.Name)

	updateNameValue := "test update"
	addData.Name = updateNameValue
	if err := db.Updates(addData).Error; err != nil {
		t.Error(err)
	}
	var updateResult TestTable
	if err := db.WithContext(queryCtx).First(&updateResult).Error; err != nil {
		t.Error(err)
	}
	assert.Equal(t, updateNameValue, updateResult.Name)

	mapUpdateNameValue := "map update value"
	if err := db.Model(&TestTable{}).Where("id = ?", addData.Id).Updates(map[string]interface{}{
		"name": mapUpdateNameValue,
	}).Error; err != nil {
		t.Error(err)
	}
	var updateMapResult TestTable
	if err := db.WithContext(queryCtx).First(&updateMapResult).Error; err != nil {
		t.Error(err)
	}
	assert.Equal(t, mapUpdateNameValue, updateMapResult.Name)
}
