package gorm_masking

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"math"
)

type MaskingPosition string

const (
	Left   MaskingPosition = "left"
	Right  MaskingPosition = "right"
	Middle MaskingPosition = "middle"
)

type CommonMasking struct {
	secretKey       string
	maskingRate     float32
	maskingPosition MaskingPosition
}

func NewChineseMasking(secretKey string, maskingRate float32, maskingPosition MaskingPosition) *CommonMasking {
	return &CommonMasking{secretKey: secretKey, maskingRate: maskingRate, maskingPosition: maskingPosition}
}

func (p *CommonMasking) Making(src string, field *schema.Field, db *gorm.DB) (string, string) {
	var desensitization, encryptValue string
	srcRune := []rune(src)
	length := len(srcRune)
	if length == 0 {
		return "", ""
	}
	if p.maskingRate >= 1 {
		desensitization = generate(length)
	} else if p.maskingRate <= 0 {
		desensitization = src
	} else {
		minMaskingLength := 1
		if int(float32(length)*p.maskingRate) > minMaskingLength {
			minMaskingLength = int(float32(length) * p.maskingRate)
		}
		switch p.maskingPosition {
		case Left:
			desensitization = generate(minMaskingLength) + string(srcRune[minMaskingLength-1:])
		case Right:
			desensitization = string(srcRune[:(length-minMaskingLength)]) + generate(minMaskingLength)
		case Middle:
			minLeft := int(math.Ceil(float64(length-minMaskingLength) / 2.0))
			desensitization = generate(minLeft) + string(srcRune[minLeft:(minLeft+minMaskingLength)]) + generate(length-minLeft-minMaskingLength)
		}
	}
	encryptValue = encrypt(src, p.secretKey)
	return desensitization, encryptValue
}

func (p *CommonMasking) UnMasking(v string, field *schema.Field, db *gorm.DB) string {
	if len(v) == 0 {
		return ""
	}
	return decryptValue(v, p.secretKey)
}
