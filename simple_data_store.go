package gorm_masking

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
)

type SimpleTableDataStore struct {
}

func (o *SimpleTableDataStore) Update(field *schema.Field, desensitizationValue, encryptionValue string, actualValue interface{}, desensitizationType MaskingType, db *gorm.DB) error {
	v, ok := field.Tag.Lookup(MaskingEncryptColumnTag)
	if !ok {
		return errors.New("未配置tag:" + MaskingEncryptColumnTag)
	}
	f := db.Statement.Schema.LookUpField(v)
	if f == nil {
		return errors.New("您配置的" + MaskingEncryptColumnTag + ":" + v + ";不存在")
	}

	if dest, ok := db.Statement.Dest.(map[string]interface{}); ok {
		dest[v] = encryptionValue
		return nil
	} else if reflect.TypeOf(db.Statement.Dest).Elem().Kind() == reflect.Struct {
		reflect.ValueOf(db.Statement.Dest).Elem().FieldByName(v).SetString(encryptionValue)
		return nil
	} else {
		return errors.New("不支持")
	}
}

func (o *SimpleTableDataStore) Create(field *schema.Field, vv reflect.Value, desensitizationValue, encryptionValue string, actualValue interface{}, desensitizationType MaskingType,
	db *gorm.DB) error {
	v, ok := field.Tag.Lookup(MaskingEncryptColumnTag)
	if !ok {
		return errors.New("未配置tag:" + MaskingEncryptColumnTag)
	}
	f := db.Statement.Schema.LookUpField(v)
	if f == nil {
		return errors.New("您配置的" + MaskingEncryptColumnTag + ":" + v + ";不存在")
	}
	f.Set(db.Statement.Context, vv, encryptionValue)
	return nil
}
