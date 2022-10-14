# gorm-masking
gorm plugin to desensitization and encryption data
[![test status](https://github.com/superqtqt/gorm-masking?branch=master "test status")](https://github.com/superqtqt/gorm-masking/actions)

## Quick Start

```go
//the code is in tests_test.go
type TestTable struct {
    Id     int    `gorm:"primary_key;column:id"`
    Name   string `gorm:"column:name" masking:"common" masking_encrypt_column:"name_id"` //masking this column
    Age    int    `gorm:"column:age"`
    NameId string `gorm:"column:name_id"` //this column store encrypt date for name
}
//init
var db *gorm.DB
var encryptKey = "1234567890123456"
dsn := "file:test?mode=memory&cache=shared"
sqliteDb, _ := gorm.Open(sqlite.Open(dsn), &gorm.Config{
Logger: logger.Default.LogMode(logger.Info),
})
sqliteDb.Use(New(&Config{Key: encryptKey}))

//query
queryCtx := context.WithValue(context.Background(), UnMaskingContext, []string{"Name"})
var queryDataByFind TestTable
if err := db.WithContext(queryCtx).Find(&queryDataByFind).Error; err != nil {
    t.Error(err)
}
```

