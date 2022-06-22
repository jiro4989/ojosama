package chars

type ExclQuesMark struct {
	Value   string
	Style   StyleType
	Meaning MeaningType
}

type StyleType int
type MeaningType int

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
		newExcl("？", styleTypeFullWidth),
		newExcl("?", styleTypeHalfWidth),
		newExcl("❓", styleTypeEmoji),
		newExcl("⁉", styleTypeEQEmoji),
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

func FindExclQuesByValue(v string) *ExclQuesMark {
	for _, mark := range eqMarks {
		if mark.Value == v {
			return &mark
		}
	}
	return nil
}

func FindExclQuesByMeaning(m MeaningType) *ExclQuesMark {
	for _, mark := range eqMarks {
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
