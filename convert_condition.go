package ojosama

import (
	"regexp"

	"github.com/ikawaha/kagome/v2/tokenizer"
)

// convertCondition は変換に使う条件。
type convertCondition struct {
	Features   []string
	Reading    string
	ReadingRe  *regexp.Regexp // オプション。設定されてる時だけ使う
	Surface    string
	SurfaceRe  *regexp.Regexp // オプション。設定されてる時だけ使う
	BaseForm   string
	BaseFormRe *regexp.Regexp // オプション。設定されてる時だけ使う
}

type convertType int

// convertConditions は変換条件のスライス。
//
// ANDで評価するかORで評価するかは、この型を使う側に依存する
type convertConditions []convertCondition

func newCond(features []string, surface string) convertCondition {
	return convertCondition{
		Features: features,
		Surface:  surface,
	}
}

func newCondRe(features []string, surfaceRe *regexp.Regexp) convertCondition {
	return convertCondition{
		Features:  features,
		SurfaceRe: surfaceRe,
	}
}

func newCondSentenceEndingParticle(surface string) convertCondition {
	return convertCondition{
		Features: posSentenceEndingParticle,
		Surface:  surface,
	}
}

func newCondAuxiliaryVerb(surface string) convertCondition {
	return convertCondition{
		Features: posAuxiliaryVerb,
		Surface:  surface,
	}
}

func newConds(surfaces []string) convertConditions {
	var c convertConditions
	for _, s := range surfaces {
		cc := convertCondition{
			Surface: s,
		}
		c = append(c, cc)
	}
	return c
}

func neStr(a, b string) bool {
	return a != "" && a != b
}

func neRe(a *regexp.Regexp, b string) bool {
	return a != nil && !a.MatchString(b)
}

func (c *convertCondition) equalsTokenData(data tokenizer.TokenData) bool {
	if 0 < len(c.Features) && !equalsFeatures(data.Features, c.Features) {
		return false
	}
	if neStr(c.Surface, data.Surface) {
		return false
	}
	if neStr(c.Reading, data.Reading) {
		return false
	}
	if neStr(c.BaseForm, data.BaseForm) {
		return false
	}
	if neRe(c.SurfaceRe, data.Surface) {
		return false
	}
	if neRe(c.ReadingRe, data.Reading) {
		return false
	}
	if neRe(c.BaseFormRe, data.BaseForm) {
		return false
	}

	return true
}

// matchAllTokenData は data がすべての c と一致した時に true を返す。
func (c *convertConditions) matchAllTokenData(data tokenizer.TokenData) bool {
	for _, cond := range *c {
		if !cond.equalsTokenData(data) {
			return false
		}
	}
	return true
}

// matchAnyTokenData は data がいずれかの c と一致した時に true を返す。
func (c *convertConditions) matchAnyTokenData(data tokenizer.TokenData) bool {
	for _, cond := range *c {
		if cond.equalsTokenData(data) {
			return true
		}
	}
	return false
}
