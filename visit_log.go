package gorm_masking

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type VisitType string

const (
	Visit  VisitType = "visit"
	Create VisitType = "create"
	Delete VisitType = "delete"
	Update VisitType = "update"
)

type VisitLog interface {
	AddVisitLog(db *gorm.DB, field *schema.Field, visitType VisitType)
}
