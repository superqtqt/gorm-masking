package gorm_masking

import (
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
)

// base data type
type MaskingType string

type DataMasking interface {
	//
	// Making
	//  @Description:  desensitization and encryption data
	//  @param src
	//  @param field
	//  @param db
	//  @return string desensitization value
	//  @return string encryption value
	//
	Making(src string, field *schema.Field, db *gorm.DB) (string, string)
	//
	// UnMasking
	//  @Description: unmasking data
	//  @param v encryption value
	//  @param field
	//  @param db
	//  @return string unmasking value
	//
	UnMasking(v string, field *schema.Field, db *gorm.DB) string
}

// plugin has implemented plugin type
const (
	Phone        MaskingType = "phone"
	ChineseName  MaskingType = "chinese_name"
	IdentifyCard MaskingType = "identify_card"
	Common       MaskingType = "common"
)

// system const value
const (
	// struct tag name
	MaskingTag = "masking"
	// struct tag name
	MaskingEncryptColumnTag = "masking_encrypt_column"
	// context key,value is []string
	UnMaskingContext = "unmasking"
)

var typeOptions map[MaskingType]DataMasking
var visitLog VisitLog

func RegisterTypeOptions(maskingType MaskingType, dFunc DataMasking) {
	typeOptions[maskingType] = dFunc
}

func GetTypeOption(maskingType MaskingType) DataMasking {
	return typeOptions[maskingType]
}

type Config struct {
	Key      string
	VisitLog VisitLog
}

func New(config *Config) *Masking {
	if config == nil {
		panic("config is nil")
	}
	if len(config.Key) == 0 {
		panic("config key is nil")
	}
	if config.VisitLog == nil {
		visitLog = &NonVisitLog{}
	} else {
		visitLog = config.VisitLog
	}
	return &Masking{config: config}
}

type Masking struct {
	config *Config
}

func (m *Masking) Name() string {
	return "masking"
}

func (m *Masking) Initialize(db *gorm.DB) error {
	typeOptions = make(map[MaskingType]DataMasking)
	RegisterTypeOptions(ChineseName, NewChineseMasking(m.config.Key, 1))

	db.Callback().Create().Before("gorm:create").Register("update_created_at", CreateMasking)
	db.Callback().Update().Before("gorm:update").Register("my_plugin:before_update", UpdateMasking)
	db.Callback().Query().After("gorm:query").Register("my_plugin:after_query", QueryMasking)
	return nil
}

func CreateMasking(db *gorm.DB) {
	var storeFunc DataStore
	if v, ok := db.Statement.ReflectValue.Type().(DataStore); ok {
		storeFunc = v
	} else {
		storeFunc = &SimpleTableDataStore{}
	}
	fieldName := db.Statement.Schema.FieldsByName
	for k, _ := range fieldName {
		v := fieldName[k]
		if desensitizationType, ok := v.Tag.Lookup(MaskingTag); ok {
			dType := MaskingType(desensitizationType)
			dFunc := typeOptions[dType]
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice:
				for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
					sliceV := db.Statement.ReflectValue.Index(i)
					if actualValue, isZero := v.ValueOf(db.Statement.Context, sliceV); !isZero {
						dValue, encryptionValue := dFunc.Making(cast.ToString(actualValue), v, db)
						v.Set(db.Statement.Context, sliceV, dValue)
						if err := storeFunc.Create(v, sliceV, dValue, encryptionValue, actualValue, dType, db); err != nil {

						}
					} else {
						continue
					}
				}
			case reflect.Struct:
				if actualValue, isZero := v.ValueOf(db.Statement.Context, db.Statement.ReflectValue); !isZero {
					dValue, encryptionValue := dFunc.Making(cast.ToString(actualValue), v, db)
					v.Set(db.Statement.Context, db.Statement.ReflectValue, dValue)
					if err := storeFunc.Create(v, db.Statement.ReflectValue, dValue, encryptionValue, actualValue, dType, db); err != nil {

					}
				} else {
					continue
				}
			default:
				db.Config.Logger.Error(db.Statement.Context, "unsupport kind [%d] of value", db.Statement.ReflectValue.Kind())
			}

		}

	}
}

func UpdateMasking(db *gorm.DB) {
	var storeFunc DataStore
	if v, ok := db.Statement.ReflectValue.Type().(DataStore); ok {
		storeFunc = v
	} else {
		storeFunc = &SimpleTableDataStore{}
	}
	if updateInfo, ok := db.Statement.Dest.(map[string]interface{}); ok {
		for updateColumn := range updateInfo {
			updateV := updateInfo[updateColumn]
			updateField := db.Statement.Schema.LookUpField(updateColumn)
			if desensitizationType, ok := updateField.Tag.Lookup(MaskingTag); ok {
				dType := MaskingType(desensitizationType)
				dFunc := typeOptions[dType]
				dValue, encryptionValue := dFunc.Making(cast.ToString(updateV), updateField, db)
				if err := storeFunc.Update(updateField, dValue, encryptionValue, updateV, dType, db); err != nil {

				}
				updateInfo[updateColumn] = dValue
			}
		}
		return
	}
	destType := reflect.TypeOf(db.Statement.Dest)
	destValue := reflect.ValueOf(db.Statement.Dest)
	if destType.Elem().Kind() == reflect.Struct {
		destType = destType.Elem()
		destValue = destValue.Elem()
		for i := 0; i < destType.NumField(); i++ {
			field := destType.Field(i)
			if desensitizationType, ok := field.Tag.Lookup(MaskingTag); ok {
				val := destValue.Field(i).String()
				if len(val) == 0 {
					continue
				}
				dType := MaskingType(desensitizationType)
				dFunc := typeOptions[dType]
				dbField := db.Statement.Schema.LookUpField(field.Name)
				dValue, encryptionValue := dFunc.Making(cast.ToString(val), dbField, db)
				if err := storeFunc.Update(dbField, dValue, encryptionValue, val, dType, db); err != nil {

				}
				destValue.Field(i).SetString(dValue)
			}

		}
	}

}

func QueryMasking(db *gorm.DB) {
	val := db.Statement.Context.Value(UnMaskingContext)
	if val != nil {
		if columns, ok := val.([]string); ok && len(columns) > 0 {
			refVal := reflect.ValueOf(db.Statement.Model).Elem()
			switch refVal.Kind() {
			case reflect.Struct:
				for i := range columns {
					decrypt(refVal, refVal.FieldByName(columns[i]), columns[i], db)
				}
			case reflect.Slice:
				if refVal.Len() >= 1 {
					for i := 0; i < refVal.Len(); i++ {
						curVal := refVal.Index(i)
						if curVal.Kind() == reflect.Pointer {
							curVal = curVal.Elem()
						}
						for j := range columns {
							decrypt(curVal, curVal.FieldByName(columns[j]), columns[j], db)
						}
					}
				}
			default:
				db.Config.Logger.Error(db.Statement.Context, "无法支持的类型:"+refVal.Kind().String())
			}
		}
	}
}
