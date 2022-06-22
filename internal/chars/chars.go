package chars

import (
	"fmt"
	"math/rand"
)

type ExclQuesMark struct {
	Value   string
	Style   StyleType
	Meaning MeaningType
}

type StyleType int
type MeaningType int
type TestMode struct {
	Pos int
}

const (
	styleTypeUnknown StyleType = iota
	styleTypeFullWidth
	styleTypeHalfWidth
	styleTypeEmoji
	styleTypeDoubleEmoji // !!
	styleTypeEQEmoji     // !?

	meaningTypeUnknown = iota
	meaningTypeExcl    // !
	meaningTypeQues    // ?
)

var (
	eqMarks = []ExclQuesMark{
		newExcl("！", styleTypeFullWidth),
		newExcl("!", styleTypeHalfWidth),
		newExcl("❗", styleTypeEmoji),
		newExcl("‼", styleTypeDoubleEmoji),
		newQues("？", styleTypeFullWidth),
		newQues("?", styleTypeHalfWidth),
		newQues("❓", styleTypeEmoji),
		newQues("⁉", styleTypeEQEmoji),
	}
)

func newExcl(v string, t StyleType) ExclQuesMark {
	return ExclQuesMark{
		Value:   v,
		Style:   t,
		Meaning: meaningTypeExcl,
	}
}

func newQues(v string, t StyleType) ExclQuesMark {
	return ExclQuesMark{
		Value:   v,
		Style:   t,
		Meaning: meaningTypeQues,
	}
}

func IsExclQuesMark(s string) (bool, *ExclQuesMark) {
	for _, v := range eqMarks {
		if v.Value == s {
			return true, &v
		}
	}
	return false, nil
}

func SampleExclQuesByValue(v string, t *TestMode) *ExclQuesMark {
	ok, got := IsExclQuesMark(v)
	if !ok {
		return nil
	}

	var s []ExclQuesMark
	for _, mark := range eqMarks {
		if mark.Meaning == got.Meaning {
			fmt.Println(mark, got)
			s = append(s, mark)
		}
	}
	if len(s) < 1 {
		return nil
	}

	if t != nil {
		// テスト用のパラメータがあるときは決め打ちで返す
		fmt.Println(s)
		return &s[t.Pos]
	}
	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
	return &s[0]
}

func FindExclQuesByStyleAndMeaning(s StyleType, m MeaningType) *ExclQuesMark {
	var eq []ExclQuesMark
	for _, mark := range eqMarks {
		if mark.Style == s {
			eq = append(eq, mark)
		}
	}
	if len(eq) < 1 {
		return nil
	}

	for _, mark := range eq {
		if mark.Meaning == m {
			return &mark
		}
	}

	return nil
}

func FindExclQuesByStyle(s StyleType) *ExclQuesMark {
	for _, mark := range eqMarks {
		if mark.Style == s {
			return &mark
		}
	}
	return nil
}
