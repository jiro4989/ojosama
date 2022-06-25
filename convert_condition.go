package ojosama

import (
	"regexp"

	"github.com/ikawaha/kagome/v2/tokenizer"
)

// convertCondition は変換に使う条件。
type convertCondition struct {
	Type        convertType
	Value       []string
	ValueRegexp *regexp.Regexp // オプション。設定されてる時だけ使う
}

type convertType int

// convertConditions は変換条件のスライス。
//
// この型の評価をする側は、すべての条件が true の時に true として解釈する。
// 一つでも false が混じったら false として解釈する。
type convertConditions []convertCondition

func newCond(features []string, surface string) convertConditions {
	return convertConditions{
		{
			Type:  convertTypeFeatures,
			Value: features,
		},
		{
			Type:  convertTypeSurface,
			Value: []string{surface},
		},
	}
}

func newCondRe(features []string, re *regexp.Regexp) convertConditions {
	return convertConditions{
		{
			Type:  convertTypeFeatures,
			Value: features,
		},
		{
			Type:        convertTypeSurface,
			ValueRegexp: re,
		},
	}
}

func newCondSentenceEndingParticle(surface string) convertConditions {
	return convertConditions{
		{
			Type:  convertTypeFeatures,
			Value: posSentenceEndingParticle,
		},
		{
			Type:  convertTypeSurface,
			Value: []string{surface},
		},
	}
}

func newCondAuxiliaryVerb(surface string) convertConditions {
	return convertConditions{
		{
			Type:  convertTypeFeatures,
			Value: posAuxiliaryVerb,
		},
		{
			Type:  convertTypeSurface,
			Value: []string{surface},
		},
	}
}

func newConds(surfaces []string) []convertConditions {
	var c []convertConditions
	for _, s := range surfaces {
		cc := convertConditions{
			{
				Type:  convertTypeSurface,
				Value: []string{s},
			},
		}
		c = append(c, cc)
	}
	return c
}

func (c *convertCondition) equalsTokenData(data tokenizer.TokenData) bool {
	switch c.Type {
	case convertTypeFeatures:
		if equalsFeatures(data.Features, c.Value) {
			return true
		}
	case convertTypeSurface:
		if c.ValueRegexp != nil {
			return c.ValueRegexp.MatchString(data.Surface)
		}
		if data.Surface == c.Value[0] {
			return true
		}
	case convertTypeBaseForm:
		if c.ValueRegexp != nil {
			return c.ValueRegexp.MatchString(data.BaseForm)
		}
		if data.BaseForm == c.Value[0] {
			return true
		}
	}
	return false
}

func (c *convertCondition) notEqualsTokenData(data tokenizer.TokenData) bool {
	switch c.Type {
	case convertTypeFeatures:
		if !equalsFeatures(data.Features, c.Value) {
			return true
		}
	case convertTypeSurface:
		if c.ValueRegexp != nil {
			return !c.ValueRegexp.MatchString(data.Surface)
		}
		if data.Surface != c.Value[0] {
			return true
		}
	case convertTypeBaseForm:
		if c.ValueRegexp != nil {
			return !c.ValueRegexp.MatchString(data.BaseForm)
		}
		if data.BaseForm != c.Value[0] {
			return true
		}
	}
	return false
}

// matchAllTokenData は data がすべての c と一致した時に true を返す。
func (c *convertConditions) matchAllTokenData(data tokenizer.TokenData) bool {
	for _, cond := range *c {
		if cond.notEqualsTokenData(data) {
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

// matchAnyMultiConvertConditions は ccs のどれかが1つでも一致すれば true を返す。
func matchAnyMultiConvertConditions(ccs []convertConditions, data tokenizer.TokenData) bool {
	for _, cs := range ccs {
		// 1つでも一致すればOK
		if cs.matchAllTokenData(data) {
			return true
		}
	}
	return false
}
