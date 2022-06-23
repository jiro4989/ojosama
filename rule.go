package ojosama

import (
	"regexp"

	"github.com/ikawaha/kagome/v2/tokenizer"
)

// convertRule は 単独のTokenに対して、Conditionsがすべてマッチしたときに変換するルール。
//
// 基本的な変換はこの型で定義する。
type convertRule struct {
	Conditions                   convertConditions
	BeforeIgnoreConditions       convertConditions // 前のTokenで条件にマッチした場合は無視する
	AfterIgnoreConditions        convertConditions // 次のTokenで条件にマッチした場合は無視する
	EnableWhenSentenceSeparation bool              // 文の区切り（単語の後に句点か読点がくる、あるいは何もない）場合だけ有効にする
	AppendLongNote               bool              // 波線を追加する
	DisablePrefix                bool              // 「お」を手前に付与しない
	Value                        string
}

// continuousConditionsConvertRule は連続する条件がすべてマッチしたときに変換するルール。
type continuousConditionsConvertRule struct {
	Conditions     []convertConditions
	AppendLongNote bool
	Value          string
}

// convertConditions は変換条件のスライス。
//
// この型の評価をする側は、すべての条件が true の時に true として解釈する。
// 一つでも false が混じったら false として解釈する。
type convertConditions []convertCondition

type convertType int

// convertCondition は変換に使う条件。
type convertCondition struct {
	Type        convertType
	Value       []string
	ValueRegexp *regexp.Regexp // オプション。設定されてる時だけ使う
}

// meaningType は言葉の意味分類。
type meaningType int

// sentenceEndingParticleConvertRule は「名詞」＋「動詞」＋「終助詞」の組み合わせによる変換ルール。
type sentenceEndingParticleConvertRule struct {
	conditions1            []convertConditions                 // 一番最初に評価されるルール
	conditions2            []convertConditions                 // 二番目に評価されるルール
	auxiliaryVerb          convertConditions                   // 助動詞。マッチしなくても次にすすむ
	sentenceEndingParticle map[meaningType][]convertConditions // 終助詞
	value                  map[meaningType][]string
}

const (
	convertTypeSurface convertType = iota + 1
	convertTypeFeatures
	convertTypeBaseForm

	meaningTypeUnknown    meaningType = iota
	meaningTypeHope                   // 希望
	meaningTypePoem                   // 詠嘆
	meaningTypeProhibiton             // 禁止
	meaningTypeCoercion               // 強制
)

var (
	sentenceEndingParticleConvertRules = []sentenceEndingParticleConvertRule{
		{
			conditions1: []convertConditions{
				{
					{Type: convertTypeFeatures, Value: nounsGeneral},
				},
				{
					{Type: convertTypeFeatures, Value: nounsSaDynamic},
				},
			},
			conditions2: []convertConditions{
				{
					{Type: convertTypeFeatures, Value: verbIndependence},
					{Type: convertTypeBaseForm, Value: []string{"する"}},
				},
				{
					{Type: convertTypeFeatures, Value: verbIndependence},
					{Type: convertTypeBaseForm, Value: []string{"やる"}},
				},
			},
			sentenceEndingParticle: map[meaningType][]convertConditions{
				meaningTypeHope: {
					newCondSentenceEndingParticle("ぜ"),
					newCondSentenceEndingParticle("よ"),
					newCondSentenceEndingParticle("べ"),
				},
				meaningTypePoem: {
					// これだけ特殊
					newCond([]string{"助詞", "副助詞／並立助詞／終助詞"}, "か"),
				},
				meaningTypeProhibiton: {
					newCondSentenceEndingParticle("な"),
				},
				meaningTypeCoercion: {
					newCondSentenceEndingParticle("ぞ"),
					newCondSentenceEndingParticle("の"),
				},
			},
			auxiliaryVerb: newCondAuxiliaryVerb("う"),
			value: map[meaningType][]string{
				meaningTypeHope: {
					"をいたしませんこと",
				},
				meaningTypePoem: {
					"をいたしますわ",
				},
				meaningTypeProhibiton: {
					"をしてはいけませんわ",
				},
				meaningTypeCoercion: {
					"をいたしますわよ",
				},
			},
		},
	}

	// continuousConditionsConvertRules は連続する条件がすべてマッチしたときに変換するルール。
	//
	// 例えば「壱百満天原サロメ」や「横断歩道」のように、複数のTokenがこの順序で連続
	// して初めて1つの意味になるような条件を定義する。
	continuousConditionsConvertRules = []continuousConditionsConvertRule{
		{
			Value:      "壱百満天原サロメ",
			Conditions: newConds([]string{"壱", "百", "満天", "原", "サロメ"}),
		},

		{
			Value:      "壱百満天原",
			Conditions: newConds([]string{"壱", "百", "満天", "原"}),
		},

		{
			Value:      "壱百満点",
			Conditions: newConds([]string{"壱", "百", "満点"}),
		},

		{
			Value:          "いたしますわ",
			AppendLongNote: true,
			Conditions: []convertConditions{
				newCond([]string{"動詞", "自立"}, "し"),
				newCond([]string{"助動詞"}, "ます"),
			},
		},

		{
			Value: "ですので",
			Conditions: []convertConditions{
				newCond([]string{"助動詞"}, "だ"),
				newCond([]string{"助詞", "接続助詞"}, "から"),
			},
		},

		{
			Value: "なんですの",
			Conditions: []convertConditions{
				newCond([]string{"助動詞"}, "な"),
				newCond([]string{"名詞", "非自立", "一般"}, "ん"),
				newCond([]string{"助動詞"}, "だ"),
			},
		},

		{
			Value: "ですわ",
			Conditions: []convertConditions{
				newCond([]string{"助動詞"}, "だ"),
				newCond([]string{"助詞", "終助詞"}, "よ"),
			},
		},
	}

	// excludeRules は変換処理を無視するルール。
	// このルールは convertRules よりも優先して評価される。
	excludeRules = []convertRule{
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"名詞", "固有名詞", "一般"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"カス"},
				},
			},
		},
		{
			Conditions: newCondRe(nounsGeneral, regexp.MustCompile(`^(ー+|～+)$`)),
		},
	}

	// convertRules は 単独のTokenに対して、Conditionsがすべてマッチしたときに変換するルール。
	//
	// 基本的な変換はここに定義する。
	convertRules = []convertRule{
		// 一人称
		newRulePronounGeneral("俺", "私"),
		newRulePronounGeneral("オレ", "ワタクシ"),
		newRulePronounGeneral("おれ", "わたくし"),
		newRulePronounGeneral("僕", "私"),
		newRulePronounGeneral("ボク", "ワタクシ"),
		newRulePronounGeneral("ぼく", "わたくし"),
		newRulePronounGeneral("あたし", "わたくし"),
		newRulePronounGeneral("わたし", "わたくし"),

		// 二人称
		newRulePronounGeneral("あなた", "貴方"),

		// 三人称
		// TODO: AfterIgnore系も簡単に定義できるようにしたい
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"名詞", "一般"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"パパ"},
				},
			},
			AfterIgnoreConditions: convertConditions{
				{
					Type:  convertTypeSurface,
					Value: []string{"上"},
				},
			},
			Value: "パパ上",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"名詞", "一般"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"ママ"},
				},
			},
			AfterIgnoreConditions: convertConditions{
				{
					Type:  convertTypeSurface,
					Value: []string{"上"},
				},
			},
			Value: "ママ上",
		},
		newRulePronounGeneral("皆", "皆様方"),
		newRuleNounsGeneral("皆様", "皆様方").disablePrefix(true),

		// こそあど言葉
		newRulePronounGeneral("これ", "こちら"),
		newRulePronounGeneral("それ", "そちら"),
		newRulePronounGeneral("あれ", "あちら"),
		newRulePronounGeneral("どれ", "どちら"),
		newRuleAdnominalAdjective("この", "こちらの"),
		newRuleAdnominalAdjective("その", "そちらの"),
		newRuleAdnominalAdjective("あの", "あちらの"),
		newRuleAdnominalAdjective("どの", "どちらの"),
		newRulePronounGeneral("ここ", "こちら"),
		newRulePronounGeneral("そこ", "そちら"),
		newRulePronounGeneral("あそこ", "あちら"),
		newRulePronounGeneral("どこ", "どちら"),
		newRuleAdnominalAdjective("こんな", "このような"),
		newRuleAdnominalAdjective("そんな", "そのような"),
		newRuleAdnominalAdjective("あんな", "あのような"),
		newRuleAdnominalAdjective("どんな", "どのような"),

		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"名詞", "非自立", "一般"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"もん"},
				},
			},
			Value: "もの",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"です"},
				},
			},
			AfterIgnoreConditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助詞", "副助詞／並立助詞／終助詞"},
				},
			},
			AppendLongNote: true,
			Value:          "ですわ",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"だ"},
				},
			},
			AfterIgnoreConditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助詞", "副助詞／並立助詞／終助詞"},
				},
			},
			AppendLongNote: true,
			Value:          "ですわ",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"動詞", "自立"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"する"},
				},
			},
			EnableWhenSentenceSeparation: true,
			AppendLongNote:               true,
			Value:                        "いたしますわ",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"動詞", "自立"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"なる"},
				},
			},
			EnableWhenSentenceSeparation: true,
			AppendLongNote:               true,
			Value:                        "なりますわ",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"動詞", "自立"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"ある"},
				},
			},
			Value: "あります",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助詞", "副助詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"じゃ"},
				},
			},
			Value: "では",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助詞", "副助詞／並立助詞／終助詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"か"},
				},
			},
			Value: "の",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助詞", "終助詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"わ"},
				},
			},
			AppendLongNote: true,
			Value:          "ですわ",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助詞", "終助詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"な"},
				},
			},
			Value: "ね",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助詞", "終助詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"さ"},
				},
			},
			Value: "",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助詞", "接続助詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"から"},
				},
			},
			Value: "ので",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助詞", "接続助詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"けど"},
				},
			},
			Value: "けれど",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助詞", "接続助詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"し"},
				},
			},
			Value: "ですし",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"まし"},
				},
			},
			BeforeIgnoreConditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"動詞", "自立"},
				},
			},
			Value: "おりまし",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"ます"},
				},
			},
			AppendLongNote: true,
			Value:          "ますわ",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"た"},
				},
			},
			EnableWhenSentenceSeparation: true,
			AppendLongNote:               true,
			Value:                        "たわ",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"だろ"},
				},
			},
			Value: "でしょう",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"ない"},
				},
			},
			BeforeIgnoreConditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"動詞", "自立"},
				},
			},
			Value: "ありません",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"動詞", "非自立"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"ください"},
				},
			},
			Value: "くださいまし",
		},
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"動詞", "非自立"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"くれ"},
				},
			},
			Value: "くださいまし",
		},
		newRuleVerbs("ありがとう", "ありがとうございますわ"),
		newRuleVerbs("じゃぁ", "それでは"),
		newRuleVerbs("じゃあ", "それでは"),
		{
			Conditions: convertConditions{
				{
					Type:  convertTypeFeatures,
					Value: []string{"動詞", "非自立"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"くれる"},
				},
			},
			Value: "くれます",
		},
		newRuleAdjectivesSelfSupporting("汚い", "きったねぇ"),
		newRuleAdjectivesSelfSupporting("きたない", "きったねぇ"),
		newRuleAdjectivesSelfSupporting("臭い", "くっせぇ"),
		newRuleAdjectivesSelfSupporting("くさい", "くっせぇ"),
		newRuleVerbs("うふ", "おほ"),
		newRuleVerbs("うふふ", "おほほ"),
		newRuleVerbs("う", "お"),
		newRuleVerbs("ふふふ", "ほほほ"),
	}
)

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

func matchAnyMultiConvertConditions(ccs []convertConditions, data tokenizer.TokenData) bool {
	for _, cs := range ccs {
		// 1つでも一致すればOK
		if cs.matchAllTokenData(data) {
			return true
		}
	}
	return false
}

var (
	pronounGeneral           = []string{"名詞", "代名詞", "一般"}
	nounsGeneral             = []string{"名詞", "一般"}
	adnominalAdjective       = []string{"連体詞"}
	adjectivesSelfSupporting = []string{"形容詞", "自立"}
	verbs                    = []string{"感動詞"}
	verbIndependence         = []string{"動詞", "自立"}
	sentenceEndingParticle   = []string{"助詞", "終助詞"}
	auxiliaryVerb            = []string{"助動詞"}
	nounsSaDynamic           = []string{"名詞", "サ変接続"}
)

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
			Value: sentenceEndingParticle,
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
			Value: auxiliaryVerb,
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

func newRule(features []string, surface, value string) convertRule {
	return convertRule{
		Conditions: convertConditions{
			{
				Type:  convertTypeFeatures,
				Value: features,
			},
			{
				Type:  convertTypeSurface,
				Value: []string{surface},
			},
		},
		Value: value,
	}
}

func newRulePronounGeneral(surface, value string) convertRule {
	return newRule(pronounGeneral, surface, value)
}

func newRuleNounsGeneral(surface, value string) convertRule {
	return newRule(nounsGeneral, surface, value)
}

func newRuleAdnominalAdjective(surface, value string) convertRule {
	return newRule(adnominalAdjective, surface, value)
}

func newRuleAdjectivesSelfSupporting(surface, value string) convertRule {
	return newRule(adjectivesSelfSupporting, surface, value)
}

func newRuleVerbs(surface, value string) convertRule {
	return newRule(verbs, surface, value)
}

func (c convertRule) disablePrefix(v bool) convertRule {
	c.DisablePrefix = v
	return c
}
