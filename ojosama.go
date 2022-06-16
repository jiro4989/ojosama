package ojosama

import (
	"fmt"
	"strings"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

type ConvertOption struct{}

type Converter struct {
	Conditions []ConvertCondition
	Value      string
}

type ConvertType int

type ConvertCondition struct {
	Type  ConvertType
	Value []string
}

const (
	ConvertTypeSurface ConvertType = iota + 1
	ConvertTypeReading
	ConvertTypeFeatures
)

var (
	convertRules = []Converter{
		{
			Conditions: []ConvertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"名詞", "一般"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"草"},
				},
			},
			Value: "ハーブ",
		},
		{
			Conditions: []ConvertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"名詞", "代名詞", "一般"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"あなた"},
				},
			},
			Value: "貴方",
		},
		{
			Conditions: []ConvertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"です"},
				},
			},
			Value: "ですわ",
		},
		{
			Conditions: []ConvertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"動詞", "自立"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"する"},
				},
			},
			Value: "いたしますわ",
		},
	}
)

func Convert(src string, opt *ConvertOption) (string, error) {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return "", err
	}

	// tokenize
	tokens := t.Tokenize(src)
	var result strings.Builder
	for i, token := range tokens {
		data := tokenizer.NewTokenData(token)
		buf := data.Surface
		fmt.Println(data.Features, buf)
	converterLoop:
		for _, c := range convertRules {
			for _, cond := range c.Conditions {
				switch cond.Type {
				case ConvertTypeFeatures:
					if !equalsFeatures(data.Features, cond.Value) {
						continue converterLoop
					}
				case ConvertTypeSurface:
					if data.Surface != cond.Value[0] {
						continue converterLoop
					}
				case ConvertTypeReading:
					if data.Reading != cond.Value[0] {
						continue converterLoop
					}
				}
			}

			buf = c.Value
			break
		}

		// 名詞 一般の場合は手前に「お」をつける。
		if equalsFeatures(data.Features, []string{"名詞", "一般"}) {
			// ただし、直後に動詞 自立がくるときは「お」をつけない。
			// 例: プレイする
			if i+1 < len(tokens) {
				data := tokenizer.NewTokenData(tokens[i+1])
				if equalsFeatures(data.Features, []string{"動詞", "自立"}) {
					goto endLoop
				}
			}
			buf = "お" + buf
		}

	endLoop:
		result.WriteString(buf)
	}
	return result.String(), nil
}

func equalsFeatures(a, b []string) bool {
	var a2 []string
	for _, v := range a {
		if v == "*" {
			break
		}
		a2 = append(a2, v)
	}

	a = a2
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
