package gorm_masking

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
)

type TestTable struct {
	Id   int    `gorm:"primary_key;column:id"`
	Name string `gorm:"column:name" masking:"common"`
	Age  int    `gorm:"column:age"`
}

func (r TestTable) TableName() string {
	return "t_test_table"
}

var db *gorm.DB

func init() {
	dsn := "file:test?mode=memory&cache=shared"
	sqliteDb, _ := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	sqliteDb.Use(New(&Config{Key: "1234567890123456"}))
	sqliteDb.AutoMigrate(&TestTable{})
	db = sqliteDb

}

func Test_BaseUse(t *testing.T) {
	if err := db.Create(&TestTable{
		Name: "test",
		Age:  18,
	}).Error; err != nil {
		t.Error(err)
	}
}
