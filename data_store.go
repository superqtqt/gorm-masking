package gorm_masking

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
)

//数据存储的实现 每个表单独实现，其中有内部实现
type DataStore interface {
	Create(field *schema.Field, v reflect.Value, desensitizationValue, encryptionValue string, actualValue interface{}, desensitizationType MaskingType, db *gorm.DB) error
	Update(field *schema.Field, desensitizationValue, encryptionValue string, actualValue interface{}, desensitizationType MaskingType, db *gorm.DB) error
}
