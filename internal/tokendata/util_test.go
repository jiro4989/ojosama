package tokendata

import (
	"testing"

	"github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/jiro4989/ojosama/internal/feat"
	"github.com/stretchr/testify/assert"
)

func TestIsKuten(t *testing.T) {
	tests := []struct {
		desc string
		data tokenizer.TokenData
		want bool
	}{
		{
			desc: "正常系: 句点の場合はtrueですわ",
			data: tokenizer.TokenData{
				Features: feat.Kuten,
				Surface:  "。",
			},
			want: true,
		},
		{
			desc: "正常系: 句点出ない場合はfalseですわ",
			data: tokenizer.TokenData{
				Features: feat.Kuten,
				Surface:  "、",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := IsKuten(tt.data)
			assert.Equal(tt.want, got)
		})
	}
}

func TestIsPoliteWord(t *testing.T) {
	tests := []struct {
		desc string
		data tokenizer.TokenData
		want bool
	}{
		{
			desc: "正常系: 「オ」で始まることばは丁寧語ですわ",
			data: tokenizer.TokenData{
				Reading: "オニギリ",
			},
			want: true,
		},
		{
			desc: "正常系: 「オ」で始まっていない言葉は丁寧語ではありませんわ",
			data: tokenizer.TokenData{
				Reading: "メカブ",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := IsPoliteWord(tt.data)
			assert.Equal(tt.want, got)
		})
	}
}

func TestEqualsFeatures(t *testing.T) {
	tests := []struct {
		desc string
		a, b []string
		want bool
	}{
		{
			desc: "正常系: * 以降は無視されますわ",
			a:    []string{"名詞", "代名詞", "一般", "*", "*"},
			b:    []string{"名詞", "代名詞", "一般"},
			want: true,
		},
		{
			desc: "正常系: 途中までしかあっていない場合はfalseですわ",
			a:    []string{"名詞", "代名詞", "一般", "*", "*"},
			b:    []string{"名詞", "代名詞"},
			want: false,
		},
		{
			desc: "正常系: 1つでもずれてたらfalseですわ",
			a:    []string{"名詞", "代名詞", "一般", "*", "*"},
			b:    []string{"名詞", "寿司", "一般"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := EqualsFeatures(tt.a, tt.b)
			assert.Equal(tt.want, got)
		})
	}
}

func TestContainsFeatures(t *testing.T) {
	tests := []struct {
		desc string
		a    [][]string
		b    []string
		want bool
	}{
		{
			desc: "正常系: どれか1つと一致すればOKですわ",
			a: [][]string{
				{"名詞", "代名詞"},
				{"名詞", "一般"},
				{"名詞", "固有名詞"},
			},
			b:    []string{"名詞", "一般", "*", "*", "*"},
			want: true,
		},
		{
			desc: "正常系: 1つも一致しなければfalseですわ",
			a: [][]string{
				{"名詞", "代名詞"},
				{"名詞", "固有名詞"},
			},
			b:    []string{"名詞", "一般", "*", "*", "*"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := ContainsFeatures(tt.a, tt.b)
			assert.Equal(tt.want, got)
		})
	}
}

func TestConatinsString(t *testing.T) {
	tests := []struct {
		desc string
		a    []string
		b    string
		want bool
	}{
		{
			desc: "正常系: どれか1つと一致すればOKですわ",
			a:    []string{"a", "b"},
			b:    "b",
			want: true,
		},
		{
			desc: "正常系: 1つも一致しなければfalseですわ",
			a:    []string{"a", "b"},
			b:    "c",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := ContainsString(tt.a, tt.b)
			assert.Equal(tt.want, got)
		})
	}
}
