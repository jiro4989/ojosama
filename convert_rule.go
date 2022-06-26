package ojosama

import (
	"regexp"
)

// convertRule は 単独のTokenに対して、Conditionsがすべてマッチしたときに変換するルール。
//
// 基本的な変換はこの型で定義する。
type convertRule struct {
	Conditions                   convertConditions // 起点になる変換条件
	BeforeIgnoreConditions       convertConditions // 前のTokenで条件にマッチした場合は無視する
	AfterIgnoreConditions        convertConditions // 次のTokenで条件にマッチした場合は無視する
	EnableWhenSentenceSeparation bool              // 文の区切り（単語の後に句点か読点がくる、あるいは何もない）場合だけ有効にする
	AppendLongNote               bool              // 波線を追加する
	DisablePrefix                bool              // 「お」を手前に付与しない
	Value                        string            // この文字列に置換する
}

/*
英語 文法 品詞
https://ja.wikibooks.org/wiki/%E8%8B%B1%E8%AA%9E/%E6%96%87%E6%B3%95/%E5%93%81%E8%A9%9E

品詞 part of speech (pos)
*/
var (
	posPronounGeneral            = []string{"名詞", "代名詞", "一般"}
	posNounsGeneral              = []string{"名詞", "一般"}
	posSpecificGeneral           = []string{"名詞", "固有名詞", "一般"}
	posNotIndependenceGeneral    = []string{"名詞", "非自立", "一般"}
	posAdnominalAdjective        = []string{"連体詞"}
	posAdjectivesSelfSupporting  = []string{"形容詞", "自立"}
	posInterjection              = []string{"感動詞"}
	posVerbIndependence          = []string{"動詞", "自立"}
	posVerbNotIndependence       = []string{"動詞", "非自立"}
	posSentenceEndingParticle    = []string{"助詞", "終助詞"}
	posSubPostpositionalParticle = []string{"助詞", "副助詞"}
	posAssistantParallelParticle = []string{"助詞", "並立助詞"}
	posSubParEndParticle         = []string{"助詞", "副助詞／並立助詞／終助詞"}
	posConnAssistant             = []string{"助詞", "接続助詞"}
	posAuxiliaryVerb             = []string{"助動詞"}
	posNounsSaDynamic            = []string{"名詞", "サ変接続"}
)

func newRule(features []string, surface, value string) convertRule {
	return convertRule{
		Conditions: convertConditions{
			newCond(features, surface),
		},
		Value: value,
	}
}

func newRulePronounGeneral(surface, value string) convertRule {
	return newRule(posPronounGeneral, surface, value)
}

func newRuleNounsGeneral(surface, value string) convertRule {
	return newRule(posNounsGeneral, surface, value)
}

func newRuleAdnominalAdjective(surface, value string) convertRule {
	return newRule(posAdnominalAdjective, surface, value)
}

func newRuleAdjectivesSelfSupporting(surface, value string) convertRule {
	return newRule(posAdjectivesSelfSupporting, surface, value)
}

func newRuleInterjection(surface, value string) convertRule {
	return newRule(posInterjection, surface, value)
}

func (c convertRule) disablePrefix(v bool) convertRule {
	c.DisablePrefix = v
	return c
}

// continuousConditionsConvertRule は連続する条件がすべてマッチしたときに変換するルール。
type continuousConditionsConvertRule struct {
	Conditions     convertConditions
	AppendLongNote bool
	Value          string
}

// sentenceEndingParticleConvertRule は「名詞」＋「動詞」＋「終助詞」の組み合わせによる変換ルール。
type sentenceEndingParticleConvertRule struct {
	conditions1            convertConditions                 // 一番最初に評価されるルール
	conditions2            convertConditions                 // 二番目に評価されるルール
	auxiliaryVerb          convertConditions                 // 助動詞。マッチしなくても次にすすむ
	sentenceEndingParticle map[meaningType]convertConditions // 終助詞
	value                  map[meaningType][]string
}

// meaningType は言葉の意味分類。
type meaningType int

const (
	convertTypeSurface convertType = iota + 1
	convertTypeFeatures
	convertTypeBaseForm

	meaningTypeUnknown     meaningType = iota
	meaningTypeHope                    // 希望
	meaningTypePoem                    // 詠嘆
	meaningTypeProhibition             // 禁止
	meaningTypeCoercion                // 強制
)

var (
	sentenceEndingParticleConvertRules = []sentenceEndingParticleConvertRule{
		{
			conditions1: convertConditions{
				{Features: posNounsGeneral},
				{Features: posNounsSaDynamic},
			},
			conditions2: convertConditions{
				{Features: posVerbIndependence, BaseForm: "する"},
				{Features: posVerbIndependence, BaseForm: "やる"},
			},
			sentenceEndingParticle: map[meaningType]convertConditions{
				meaningTypeHope: {
					newCondSentenceEndingParticle("ぜ"),
					newCondSentenceEndingParticle("よ"),
					newCondSentenceEndingParticle("べ"),
				},
				meaningTypePoem: {
					// これだけ特殊
					newCond([]string{"助詞", "副助詞／並立助詞／終助詞"}, "か"),
				},
				meaningTypeProhibition: {
					newCondSentenceEndingParticle("な"),
				},
				meaningTypeCoercion: {
					newCondSentenceEndingParticle("ぞ"),
					newCondSentenceEndingParticle("の"),
				},
			},
			auxiliaryVerb: convertConditions{
				newCondAuxiliaryVerb("う"),
			},
			value: map[meaningType][]string{
				meaningTypeHope: {
					"をいたしませんこと",
				},
				meaningTypePoem: {
					"をいたしますわ",
				},
				meaningTypeProhibition: {
					"をしてはいけませんわ",
				},
				meaningTypeCoercion: {
					"をいたしますわよ",
				},
			},
		},
	}

	condNounsGeneral    = convertCondition{Features: posNounsGeneral}
	condPronounsGeneral = convertCondition{Features: posPronounGeneral}

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
			Conditions: convertConditions{
				newCond([]string{"動詞", "自立"}, "し"),
				newCond([]string{"助動詞"}, "ます"),
			},
		},

		{
			Value: "ですので",
			Conditions: convertConditions{
				newCond([]string{"助動詞"}, "だ"),
				newCond([]string{"助詞", "接続助詞"}, "から"),
			},
		},

		{
			Value: "なんですの",
			Conditions: convertConditions{
				newCond([]string{"助動詞"}, "な"),
				newCond([]string{"名詞", "非自立", "一般"}, "ん"),
				newCond([]string{"助動詞"}, "だ"),
			},
		},

		{
			Value: "ですわ",
			Conditions: convertConditions{
				newCond([]string{"助動詞"}, "だ"),
				newCond([]string{"助詞", "終助詞"}, "よ"),
			},
		},

		{
			Value: "なんですの",
			Conditions: convertConditions{
				newCond(posPronounGeneral, "なん"),
				newCond(posSubPostpositionalParticle, "じゃ"),
			},
		},
		{
			Value: "なんですの",
			Conditions: convertConditions{
				newCond(posPronounGeneral, "なん"),
				newCond(posAuxiliaryVerb, "だ"),
			},
		},
		{
			Value: "なんですの",
			Conditions: convertConditions{
				newCond(posPronounGeneral, "なん"),
				newCond(posAssistantParallelParticle, "や"),
			},
		},

		{
			Value: "@1ですの",
			Conditions: convertConditions{
				condNounsGeneral,
				newCond(posAuxiliaryVerb, "じゃ"),
			},
		},
		{
			Value: "@1ですの",
			Conditions: convertConditions{
				condNounsGeneral,
				newCond(posAuxiliaryVerb, "だ"),
			},
		},
		{
			Value: "@1ですの",
			Conditions: convertConditions{
				condNounsGeneral,
				newCond(posAuxiliaryVerb, "や"),
			},
		},

		{
			Value: "@1ですの",
			Conditions: convertConditions{
				condPronounsGeneral,
				newCond(posAuxiliaryVerb, "じゃ"),
			},
		},
		{
			Value: "@1ですの",
			Conditions: convertConditions{
				condPronounsGeneral,
				newCond(posAuxiliaryVerb, "だ"),
			},
		},
		{
			Value: "@1ですの",
			Conditions: convertConditions{
				condPronounsGeneral,
				newCond(posAuxiliaryVerb, "や"),
			},
		},
	}

	// excludeRules は変換処理を無視するルール。
	// このルールは convertRules よりも優先して評価される。
	excludeRules = []convertRule{
		{
			Conditions: convertConditions{
				newCond(posSpecificGeneral, "カス"),
			},
		},
		{
			Conditions: convertConditions{
				newCondRe(posNounsGeneral, regexp.MustCompile(`^(ー+|～+)$`)),
			},
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
				newCond(posNounsGeneral, "パパ"),
			},
			AfterIgnoreConditions: convertConditions{
				{Surface: "上"},
			},
			Value: "パパ上",
		},
		{
			Conditions: convertConditions{
				newCond(posNounsGeneral, "ママ"),
			},
			AfterIgnoreConditions: convertConditions{
				{Surface: "上"},
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
				newCond(posNotIndependenceGeneral, "もん"),
			},
			Value: "もの",
		},
		{
			Conditions: convertConditions{
				newCond(posAuxiliaryVerb, "です"),
			},
			AfterIgnoreConditions: convertConditions{
				{Features: posSubParEndParticle},
			},
			AppendLongNote: true,
			Value:          "ですわ",
		},
		{
			Conditions: convertConditions{
				newCond(posAuxiliaryVerb, "だ"),
			},
			AfterIgnoreConditions: convertConditions{
				{Features: posSubParEndParticle},
			},
			AppendLongNote: true,
			Value:          "ですわ",
		},
		{
			Conditions: convertConditions{
				newCond(posVerbIndependence, "する"),
			},
			EnableWhenSentenceSeparation: true,
			AppendLongNote:               true,
			Value:                        "いたしますわ",
		},
		{
			Conditions: convertConditions{
				newCond(posVerbIndependence, "なる"),
			},
			EnableWhenSentenceSeparation: true,
			AppendLongNote:               true,
			Value:                        "なりますわ",
		},
		{
			Conditions: convertConditions{
				newCond(posVerbIndependence, "ある"),
			},
			Value: "あります",
		},
		{
			Conditions: convertConditions{
				newCond(posSubPostpositionalParticle, "じゃ"),
			},
			Value: "では",
		},
		{
			Conditions: convertConditions{
				newCond(posSubParEndParticle, "か"),
			},
			Value: "の",
		},
		{
			Conditions: convertConditions{
				newCond(posSentenceEndingParticle, "わ"),
			},
			AppendLongNote: true,
			Value:          "ですわ",
		},
		{
			Conditions: convertConditions{
				newCond(posSentenceEndingParticle, "な"),
			},
			Value: "ね",
		},
		{
			Conditions: convertConditions{
				newCond(posSentenceEndingParticle, "さ"),
			},
			Value: "",
		},
		{
			Conditions: convertConditions{
				newCond(posConnAssistant, "から"),
			},
			Value: "ので",
		},
		{
			Conditions: convertConditions{
				newCond(posConnAssistant, "けど"),
			},
			Value: "けれど",
		},
		{
			Conditions: convertConditions{
				newCond(posConnAssistant, "し"),
			},
			Value: "ですし",
		},
		{
			Conditions: convertConditions{
				newCond(posAuxiliaryVerb, "まし"),
			},
			BeforeIgnoreConditions: convertConditions{
				{Features: posVerbIndependence},
			},
			Value: "おりまし",
		},
		{
			Conditions: convertConditions{
				newCond(posAuxiliaryVerb, "ます"),
			},
			AppendLongNote: true,
			Value:          "ますわ",
		},
		{
			Conditions: convertConditions{
				newCond(posAuxiliaryVerb, "た"),
			},
			EnableWhenSentenceSeparation: true,
			AppendLongNote:               true,
			Value:                        "たわ",
		},
		{
			Conditions: convertConditions{
				newCond(posAuxiliaryVerb, "だろ"),
			},
			Value: "でしょう",
		},
		{
			Conditions: convertConditions{
				newCond(posAuxiliaryVerb, "ない"),
			},
			BeforeIgnoreConditions: convertConditions{
				{Features: posVerbIndependence},
			},
			Value: "ありません",
		},
		{
			Conditions: convertConditions{
				newCond(posVerbNotIndependence, "ください"),
			},
			Value: "くださいまし",
		},
		{
			Conditions: convertConditions{
				newCond(posVerbNotIndependence, "くれ"),
			},
			Value: "くださいまし",
		},
		newRuleInterjection("ありがとう", "ありがとうございますわ"),
		newRuleInterjection("じゃぁ", "それでは"),
		newRuleInterjection("じゃあ", "それでは"),
		{
			Conditions: convertConditions{
				newCond(posVerbNotIndependence, "くれる"),
			},
			Value: "くれます",
		},
		newRuleAdjectivesSelfSupporting("汚い", "きったねぇ"),
		newRuleAdjectivesSelfSupporting("きたない", "きったねぇ"),
		newRuleAdjectivesSelfSupporting("臭い", "くっせぇ"),
		newRuleAdjectivesSelfSupporting("くさい", "くっせぇ"),
		newRuleInterjection("うふ", "おほ"),
		newRuleInterjection("うふふ", "おほほ"),
		newRuleInterjection("う", "お"),
		newRuleInterjection("ふふふ", "ほほほ"),
	}
)
