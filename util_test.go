package ojosama

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

			got := equalsFeatures(tt.a, tt.b)
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

			got := containsFeatures(tt.a, tt.b)
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

			got := containsString(tt.a, tt.b)
			assert.Equal(tt.want, got)
		})
	}
}
