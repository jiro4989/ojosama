package converter

import (
	"regexp"

	"github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/jiro4989/ojosama/internal/pos"
	"github.com/jiro4989/ojosama/internal/tokendata"
)

// ConvertCondition は変換に使う条件。
type ConvertCondition struct {
	Features   []string
	Reading    string
	ReadingRe  *regexp.Regexp // オプション。設定されてる時だけ使う
	Surface    string
	SurfaceRe  *regexp.Regexp // オプション。設定されてる時だけ使う
	BaseForm   string
	BaseFormRe *regexp.Regexp // オプション。設定されてる時だけ使う
}

// ConvertConditions は変換条件のスライス。
//
// ANDで評価するかORで評価するかは、この型を使う側に依存する
type ConvertConditions []ConvertCondition

func newCond(features []string, surface string) ConvertCondition {
	return ConvertCondition{
		Features: features,
		Surface:  surface,
	}
}

func newCondRe(features []string, surfaceRe *regexp.Regexp) ConvertCondition {
	return ConvertCondition{
		Features:  features,
		SurfaceRe: surfaceRe,
	}
}

func newCondSentenceEndingParticle(surface string) ConvertCondition {
	return ConvertCondition{
		Features: pos.SentenceEndingParticle,
		Surface:  surface,
	}
}

func newCondAuxiliaryVerb(surface string) ConvertCondition {
	return ConvertCondition{
		Features: pos.AuxiliaryVerb,
		Surface:  surface,
	}
}

func newConds(surfaces []string) ConvertConditions {
	var c ConvertConditions
	for _, s := range surfaces {
		cc := ConvertCondition{
			Surface: s,
		}
		c = append(c, cc)
	}
	return c
}

func isNotEmptyStringAndDoesntEqualString(a, b string) bool {
	return a != "" && a != b
}

func isNotNilAndDoesntMatchString(a *regexp.Regexp, b string) bool {
	return a != nil && !a.MatchString(b)
}

func (c *ConvertCondition) EqualsTokenData(data tokenizer.TokenData) bool {
	if 0 < len(c.Features) && !tokendata.EqualsFeatures(data.Features, c.Features) {
		return false
	}
	if isNotEmptyStringAndDoesntEqualString(c.Surface, data.Surface) {
		return false
	}
	if isNotEmptyStringAndDoesntEqualString(c.Reading, data.Reading) {
		return false
	}
	if isNotEmptyStringAndDoesntEqualString(c.BaseForm, data.BaseForm) {
		return false
	}
	if isNotNilAndDoesntMatchString(c.SurfaceRe, data.Surface) {
		return false
	}
	if isNotNilAndDoesntMatchString(c.ReadingRe, data.Reading) {
		return false
	}
	if isNotNilAndDoesntMatchString(c.BaseFormRe, data.BaseForm) {
		return false
	}

	return true
}

// matchAllTokenData は data がすべての c と一致した時に true を返す。
func (c *ConvertConditions) MatchAllTokenData(data tokenizer.TokenData) bool {
	for _, cond := range *c {
		if !cond.EqualsTokenData(data) {
			return false
		}
	}
	return true
}

// matchAnyTokenData は data がいずれかの c と一致した時に true を返す。
func (c *ConvertConditions) MatchAnyTokenData(data tokenizer.TokenData) bool {
	for _, cond := range *c {
		if cond.EqualsTokenData(data) {
			return true
		}
	}
	return false
}
