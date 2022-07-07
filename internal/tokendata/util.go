package tokendata

import (
	"strings"

	"github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/jiro4989/ojosama/internal/feat"
)

func IsKuten(data tokenizer.TokenData) bool {
	return EqualsFeatures(data.Features, feat.Kuten) && data.Surface == "。"
}

// IsPoliteWord は丁寧語かどうかを判定する。
// 読みがオで始まる言葉も true になる。
func IsPoliteWord(data tokenizer.TokenData) bool {
	return strings.HasPrefix(data.Reading, "オ")
}

// EqualsFeatures は Features が等しいかどうかを判定する。
//
// featuresが空の時は * が設定されているため、 * が出現したら以降は無視する。
//
// 例えば {名詞,代名詞,一般,*,*} と {名詞,代名詞,一般} を比較したとき、単純にス
// ライスを完全一致で比較すると false になるが、この関数に関しては * 以降を無視
// するため true になる。
func EqualsFeatures(a, b []string) bool {
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

// ContainsFeatures は a の中に b が含まれるかを判定する。
//
// features用。
func ContainsFeatures(a [][]string, b []string) bool {
	for _, a2 := range a {
		if EqualsFeatures(b, a2) {
			return true
		}
	}
	return false
}

// ContainsString は a の中に b が含まれるかを判定する。
func ContainsString(a []string, b string) bool {
	for _, a2 := range a {
		if a2 == b {
			return true
		}
	}
	return false
}
