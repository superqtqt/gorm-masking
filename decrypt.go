package gorm_masking

import (
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"reflect"
)

func decrypt(src, refValue reflect.Value, filedName string, db *gorm.DB) {
	dbField := db.Statement.Schema.LookUpField(filedName)
	if desensitizationType, ok := dbField.Tag.Lookup(MaskingTag); ok {
		dType := MaskingType(desensitizationType)
		dFunc := typeOptions[dType]
		encryptColumn, _ := dbField.Tag.Lookup(MaskingEncryptColumnTag)
		encryptValue := src.FieldByName(encryptColumn).String()
		if len(encryptValue) > 0 {

			actualValue := dFunc.UnMasking(encryptValue, dbField, db)
			log.Infof("encryptValue:%s,actualValue:%s", encryptValue, actualValue)
			refValue.SetString(actualValue)
		}
	}
}
