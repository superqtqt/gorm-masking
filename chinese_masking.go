package gorm_masking

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ChineseDataMasking struct {
	secretKey        string
	displayManLength int
}

func NewChineseMasking(secretKey string, displayManLength int) *ChineseDataMasking {
	return &ChineseDataMasking{secretKey: secretKey, displayManLength: displayManLength}
}

func (p *ChineseDataMasking) Making(src string, field *schema.Field, db *gorm.DB) (string, string) {
	var desensitization, encryptValue string
	srcRune := []rune(src)
	length := len(srcRune)
	switch {
	case length == 0:
		return "", ""
	default:
		desensitization = string(srcRune[0]) + generate(length-1)
	}
	encryptValue = encrypt(src, p.secretKey)
	return desensitization, encryptValue
}

func (p *ChineseDataMasking) UnMasking(v string, field *schema.Field, db *gorm.DB) string {
	if len(v) == 0 {
		return ""
	}
	return decryptValue(v, p.secretKey)
}
