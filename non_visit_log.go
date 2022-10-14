package gorm_masking

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type NonVisitLog struct {
}

func (n *NonVisitLog) AddVisitLog(db *gorm.DB, field *schema.Field, visitType VisitType) {

}
