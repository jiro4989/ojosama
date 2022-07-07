package converter

import (
	"regexp"
	"testing"

	"github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/stretchr/testify/assert"
)

func TestEqualsTokenData(t *testing.T) {
	tests := []struct {
		desc string
		c    *ConvertCondition
		data tokenizer.TokenData
		want bool
	}{
		{
			desc: "正常系: 複数の条件が定義されている場合ANDで評価いたしますわ",
			c: &ConvertCondition{
				Features:   []string{"名詞"},
				Reading:    "a",
				Surface:    "b",
				BaseForm:   "c",
				ReadingRe:  regexp.MustCompile(`a+`),
				SurfaceRe:  regexp.MustCompile(`b+`),
				BaseFormRe: regexp.MustCompile(`c+`),
			},
			data: tokenizer.TokenData{
				Features: []string{"名詞"},
				Reading:  "a",
				Surface:  "b",
				BaseForm: "c",
			},
			want: true,
		},
		{
			desc: "正常系: Featuresが存在して、且つ不一致な場合は false ですわ",
			c: &ConvertCondition{
				Features: []string{"名詞"},
			},
			data: tokenizer.TokenData{
				Features: []string{"動詞"},
			},
			want: false,
		},
		{
			desc: "正常系: Surfaceが存在して、且つ不一致な場合は false ですわ",
			c: &ConvertCondition{
				Surface: "a",
			},
			data: tokenizer.TokenData{
				Surface: "b",
			},
			want: false,
		},
		{
			desc: "正常系: Readingが存在して、且つ不一致な場合は false ですわ",
			c: &ConvertCondition{
				Reading: "a",
			},
			data: tokenizer.TokenData{
				Reading: "b",
			},
			want: false,
		},
		{
			desc: "正常系: BaseFormが存在して、且つ不一致な場合は false ですわ",
			c: &ConvertCondition{
				BaseForm: "a",
			},
			data: tokenizer.TokenData{
				BaseForm: "b",
			},
			want: false,
		},
		{
			desc: "正常系: SurfaceReが存在して、且つ不一致な場合は false ですわ",
			c: &ConvertCondition{
				SurfaceRe: regexp.MustCompile(`a+`),
			},
			data: tokenizer.TokenData{
				Surface: "b",
			},
			want: false,
		},
		{
			desc: "正常系: ReadingReが存在して、且つ不一致な場合は false ですわ",
			c: &ConvertCondition{
				ReadingRe: regexp.MustCompile(`a+`),
			},
			data: tokenizer.TokenData{
				Reading: "b",
			},
			want: false,
		},
		{
			desc: "正常系: BaseFormReが存在して、且つ不一致な場合は false ですわ",
			c: &ConvertCondition{
				BaseFormRe: regexp.MustCompile(`a+`),
			},
			data: tokenizer.TokenData{
				BaseForm: "b",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := tt.c.EqualsTokenData(tt.data)
			assert.Equal(tt.want, got)
		})
	}
}

func TestMatchAllTokenData(t *testing.T) {
	tests := []struct {
		desc string
		c    ConvertConditions
		data tokenizer.TokenData
		want bool
	}{
		{
			desc: "正常系: すべての評価がtrueの場合にtrueを返しますわ",
			c: ConvertConditions{
				{
					Features: []string{"名詞"},
					Reading:  "a",
				},
				{
					Features: []string{"名詞"},
					Surface:  "b",
				},
			},
			data: tokenizer.TokenData{
				Features: []string{"名詞"},
				Reading:  "a",
				Surface:  "b",
				BaseForm: "c",
			},
			want: true,
		},
		{
			desc: "正常系: 一つで不一致の場合はfalseですわ",
			c: ConvertConditions{
				{
					Features: []string{"名詞"},
					Reading:  "a",
				},
				{
					Features: []string{"名詞"},
					Surface:  "c",
				},
			},
			data: tokenizer.TokenData{
				Features: []string{"名詞"},
				Reading:  "a",
				Surface:  "b",
				BaseForm: "c",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := tt.c.MatchAllTokenData(tt.data)
			assert.Equal(tt.want, got)
		})
	}
}

func TestMatchAnyTokenData(t *testing.T) {
	tests := []struct {
		desc string
		c    ConvertConditions
		data tokenizer.TokenData
		want bool
	}{
		{
			desc: "正常系: どれか1つがtrueになればtrueですわ",
			c: ConvertConditions{
				{
					Features: []string{"名詞"},
					Reading:  "z",
				},
				{
					Features: []string{"名詞"},
					Surface:  "b",
				},
			},
			data: tokenizer.TokenData{
				Features: []string{"名詞"},
				Reading:  "a",
				Surface:  "b",
				BaseForm: "c",
			},
			want: true,
		},
		{
			desc: "正常系: すべて不一致の場合はfalseですわ",
			c: ConvertConditions{
				{
					Features: []string{"名詞"},
					Reading:  "z",
				},
				{
					Features: []string{"名詞"},
					Surface:  "z",
				},
			},
			data: tokenizer.TokenData{
				Features: []string{"名詞"},
				Reading:  "a",
				Surface:  "b",
				BaseForm: "c",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := tt.c.MatchAnyTokenData(tt.data)
			assert.Equal(tt.want, got)
		})
	}
}
