package chars

import (
	"math/rand"
)

type ExclamationQuestionMark struct {
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

	meaningTypeUnknown = iota
	meaningTypeExcl    // !
	meaningTypeQues    // ?
	meaningTypeEQ      // !?
)

var (
	eqMarks = []ExclamationQuestionMark{
		newExcl("！", styleTypeFullWidth),
		newExcl("!", styleTypeHalfWidth),
		newExcl("❗", styleTypeEmoji),
		newExcl("‼", styleTypeDoubleEmoji),
		newQues("？", styleTypeFullWidth),
		newQues("?", styleTypeHalfWidth),
		newQues("❓", styleTypeEmoji),
		newEQ("!?", styleTypeHalfWidth),
		newEQ("⁉", styleTypeEmoji),
	}
)

func newExcl(v string, t StyleType) ExclamationQuestionMark {
	return ExclamationQuestionMark{
		Value:   v,
		Style:   t,
		Meaning: meaningTypeExcl,
	}
}

func newQues(v string, t StyleType) ExclamationQuestionMark {
	return ExclamationQuestionMark{
		Value:   v,
		Style:   t,
		Meaning: meaningTypeQues,
	}
}

func newEQ(v string, t StyleType) ExclamationQuestionMark {
	return ExclamationQuestionMark{
		Value:   v,
		Style:   t,
		Meaning: meaningTypeEQ,
	}
}

func IsExclamationQuestionMark(s string) (bool, *ExclamationQuestionMark) {
	for _, v := range eqMarks {
		if v.Value == s {
			return true, &v
		}
	}
	return false, nil
}

func SampleExclamationQuestionByValue(v string, t *TestMode) *ExclamationQuestionMark {
	ok, got := IsExclamationQuestionMark(v)
	if !ok {
		return nil
	}

	var s []ExclamationQuestionMark
	for _, mark := range eqMarks {
		if mark.Meaning == got.Meaning {
			s = append(s, mark)
		}
	}
	// 到達しないはずだけれど一応いれてる
	if len(s) < 1 {
		return nil
	}

	if t != nil {
		// テスト用のパラメータがあるときは決め打ちで返す
		return &s[t.Pos]
	}
	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
	return &s[0]
}

func FindExclamationQuestionByStyleAndMeaning(s StyleType, m MeaningType) *ExclamationQuestionMark {
	var eq []ExclamationQuestionMark
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
