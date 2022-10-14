package gorm_masking

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
)

// SameTableDataStore
//  @Description:  save the desensitization value and encryption value in the same table
type SameTableDataStore struct {
}

func (o *SameTableDataStore) Update(field *schema.Field, desensitizationValue, encryptionValue string, actualValue interface{}, desensitizationType MaskingType, db *gorm.DB) error {
	v, ok := field.Tag.Lookup(MaskingEncryptColumnTag)
	if !ok {
		return errors.New(field.Name + " has no config tag:" + MaskingEncryptColumnTag)
	}
	f := db.Statement.Schema.LookUpField(v)
	if f == nil {
		return errors.New(field.Name + "config the tag " + MaskingEncryptColumnTag + ":" + v + " not exist")
	}

	if dest, ok := db.Statement.Dest.(map[string]interface{}); ok {
		dest[v] = encryptionValue
		return nil
	} else if reflect.TypeOf(db.Statement.Dest).Elem().Kind() == reflect.Struct {
		fieldName := field.Schema.LookUpField(v).StructField.Name
		reflect.ValueOf(db.Statement.Dest).Elem().FieldByName(fieldName).SetString(encryptionValue)
		return nil
	} else {
		return errors.New("不支持")
	}
}

func (o *SameTableDataStore) Create(field *schema.Field, vv reflect.Value, desensitizationValue, encryptionValue string, actualValue interface{}, desensitizationType MaskingType,
	db *gorm.DB) error {
	v, ok := field.Tag.Lookup(MaskingEncryptColumnTag)
	if !ok {
		return errors.New(field.Name + " has no config tag:" + MaskingEncryptColumnTag)
	}
	f := db.Statement.Schema.LookUpField(v)
	if f == nil {
		return errors.New(field.Name + "config the tag " + MaskingEncryptColumnTag + ":" + v + " not exist")
	}
	f.Set(db.Statement.Context, vv, encryptionValue)
	return nil
}
