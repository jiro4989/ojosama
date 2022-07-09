package converter

import (
	"regexp"

	"github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/jiro4989/ojosama/internal/pos"
)

// ConvertRule は 単独のTokenに対して、Conditionsがすべてマッチしたときに変換するルール。
//
// 基本的な変換はこの型で定義する。
type ConvertRule struct {
	Conditions                   ConvertConditions // 起点になる変換条件
	BeforeIgnoreConditions       ConvertConditions // 前のTokenで条件にマッチした場合は無視する
	AfterIgnoreConditions        ConvertConditions // 次のTokenで条件にマッチした場合は無視する
	EnableWhenSentenceSeparation bool              // 文の区切り（単語の後に句点か読点がくる、あるいは何もない）場合だけ有効にする
	AppendLongNote               bool              // 波線を追加する
	DisablePrefix                bool              // 「お」を手前に付与しない
	EnableKutenToExclamation     bool              // 直後に句点が来たとき確率で！に変換する
	Value                        string            // この文字列に置換する
}

func newRule(features []string, surface, value string) ConvertRule {
	return ConvertRule{
		Conditions: ConvertConditions{
			newCond(features, surface),
		},
		Value: value,
	}
}

func newRulePronounGeneral(surface, value string) ConvertRule {
	return newRule(pos.PronounGeneral, surface, value)
}

func newRuleNounsGeneral(surface, value string) ConvertRule {
	return newRule(pos.NounsGeneral, surface, value)
}

func newRuleAdnominalAdjective(surface, value string) ConvertRule {
	return newRule(pos.AdnominalAdjective, surface, value)
}

func newRuleAdjectivesSelfSupporting(surface, value string) ConvertRule {
	return newRule(pos.AdjectivesSelfSupporting, surface, value)
}

func newRuleInterjection(surface, value string) ConvertRule {
	return newRule(pos.Interjection, surface, value)
}

func (c ConvertRule) disablePrefix(v bool) ConvertRule {
	c.DisablePrefix = v
	return c
}

// ContinuousConditionsConvertRule は連続する条件がすべてマッチしたときに変換するルール。
type ContinuousConditionsConvertRule struct {
	Conditions               ConvertConditions
	AppendLongNote           bool
	EnableKutenToExclamation bool
	Value                    string
}

// SentenceEndingParticleConvertRule は「名詞」＋「動詞」＋「終助詞」の組み合わせによる変換ルール。
type SentenceEndingParticleConvertRule struct {
	Conditions1            ConvertConditions                 // 一番最初に評価されるルール
	Conditions2            ConvertConditions                 // 二番目に評価されるルール
	AuxiliaryVerb          ConvertConditions                 // 助動詞。マッチしなくても次にすすむ
	SentenceEndingParticle map[MeaningType]ConvertConditions // 終助詞
	Value                  map[MeaningType][]string
}

// MeaningType は言葉の意味分類。
type MeaningType int

const (
	meaningTypeUnknown     MeaningType = iota
	meaningTypeHope                    // 希望
	meaningTypePoem                    // 詠嘆
	meaningTypeProhibition             // 禁止
	meaningTypeCoercion                // 強制
)

var (
	SentenceEndingParticleConvertRules = []SentenceEndingParticleConvertRule{
		{
			Conditions1: ConvertConditions{
				{Features: pos.NounsGeneral},
				{Features: pos.NounsSaDynamic},
			},
			Conditions2: ConvertConditions{
				{Features: pos.VerbIndependence, BaseForm: "する"},
				{Features: pos.VerbIndependence, BaseForm: "やる"},
			},
			SentenceEndingParticle: map[MeaningType]ConvertConditions{
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
			AuxiliaryVerb: ConvertConditions{
				newCondAuxiliaryVerb("う"),
			},
			Value: map[MeaningType][]string{
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

	condNounsGeneral    = ConvertCondition{Features: pos.NounsGeneral}
	condPronounsGeneral = ConvertCondition{Features: pos.PronounGeneral}

	// continuousConditionsConvertRules は連続する条件がすべてマッチしたときに変換するルール。
	//
	// 例えば「壱百満天原サロメ」や「横断歩道」のように、複数のTokenがこの順序で連続
	// して初めて1つの意味になるような条件を定義する。
	ContinuousConditionsConvertRules = []ContinuousConditionsConvertRule{
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
			Conditions: ConvertConditions{
				newCond([]string{"動詞", "自立"}, "し"),
				newCond([]string{"助動詞"}, "ます"),
			},
			EnableKutenToExclamation: true,
		},

		{
			Value: "ですので",
			Conditions: ConvertConditions{
				newCond([]string{"助動詞"}, "だ"),
				newCond([]string{"助詞", "接続助詞"}, "から"),
			},
			EnableKutenToExclamation: true,
		},

		{
			Value: "なんですの",
			Conditions: ConvertConditions{
				newCond([]string{"助動詞"}, "な"),
				newCond([]string{"名詞", "非自立", "一般"}, "ん"),
				newCond([]string{"助動詞"}, "だ"),
			},
			EnableKutenToExclamation: true,
		},

		{
			Value: "ですわ",
			Conditions: ConvertConditions{
				newCond([]string{"助動詞"}, "だ"),
				newCond([]string{"助詞", "終助詞"}, "よ"),
			},
			EnableKutenToExclamation: true,
		},

		{
			Value: "なんですの",
			Conditions: ConvertConditions{
				newCond(pos.PronounGeneral, "なん"),
				newCond(pos.SubPostpositionalParticle, "じゃ"),
			},
			EnableKutenToExclamation: true,
		},
		{
			Value: "なんですの",
			Conditions: ConvertConditions{
				newCond(pos.PronounGeneral, "なん"),
				newCond(pos.AuxiliaryVerb, "だ"),
			},
			EnableKutenToExclamation: true,
		},
		{
			Value: "なんですの",
			Conditions: ConvertConditions{
				newCond(pos.PronounGeneral, "なん"),
				newCond(pos.AssistantParallelParticle, "や"),
			},
			EnableKutenToExclamation: true,
		},

		{
			Value: "@1ですの",
			Conditions: ConvertConditions{
				condNounsGeneral,
				newCond(pos.AuxiliaryVerb, "じゃ"),
			},
			EnableKutenToExclamation: true,
		},
		{
			Value: "@1ですの",
			Conditions: ConvertConditions{
				condNounsGeneral,
				newCond(pos.AuxiliaryVerb, "だ"),
			},
			EnableKutenToExclamation: true,
		},
		{
			Value: "@1ですの",
			Conditions: ConvertConditions{
				condNounsGeneral,
				newCond(pos.AuxiliaryVerb, "や"),
			},
			EnableKutenToExclamation: true,
		},

		{
			Value: "@1ですの",
			Conditions: ConvertConditions{
				condPronounsGeneral,
				newCond(pos.AuxiliaryVerb, "じゃ"),
			},
			EnableKutenToExclamation: true,
		},
		{
			Value: "@1ですの",
			Conditions: ConvertConditions{
				condPronounsGeneral,
				newCond(pos.AuxiliaryVerb, "だ"),
			},
			EnableKutenToExclamation: true,
		},
		{
			Value: "@1ですの",
			Conditions: ConvertConditions{
				condPronounsGeneral,
				newCond(pos.AuxiliaryVerb, "や"),
			},
			EnableKutenToExclamation: true,
		},
	}

	// ExcludeRules は変換処理を無視するルール。
	// このルールは ConvertRules よりも優先して評価される。
	ExcludeRules = []ConvertRule{
		{
			Conditions: ConvertConditions{
				newCond(pos.SpecificGeneral, "カス"),
			},
		},
		{
			Conditions: ConvertConditions{
				newCondRe(pos.NounsGeneral, regexp.MustCompile(`^(ー+|～+)$`)),
			},
		},
	}

	// ConvertRules は 単独のTokenに対して、Conditionsがすべてマッチしたときに変換するルール。
	//
	// 基本的な変換はここに定義する。
	ConvertRules = []ConvertRule{
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
		newRulePronounGeneral("あんた", "貴方"),
		newRulePronounGeneral("おまえ", "貴方"),
		newRulePronounGeneral("お前", "貴方"),
		newRulePronounGeneral("てめぇ", "貴方"),
		newRulePronounGeneral("てめえ", "貴方"),
		newRuleNounsGeneral("貴様", "貴方").disablePrefix(true),
		// newRulePronounGeneral("きさま", "貴方"),
		// newRulePronounGeneral("そなた", "貴方"),
		newRulePronounGeneral("君", "貴方"),

		// 三人称
		// TODO: AfterIgnore系も簡単に定義できるようにしたい
		{
			Conditions: ConvertConditions{
				newCond(pos.NounsGeneral, "パパ"),
			},
			AfterIgnoreConditions: ConvertConditions{
				{Surface: "上"},
			},
			Value: "パパ上",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.NounsGeneral, "ママ"),
			},
			AfterIgnoreConditions: ConvertConditions{
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
			Conditions: ConvertConditions{
				newCond(pos.NotIndependenceGeneral, "もん"),
			},
			Value: "もの",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.AuxiliaryVerb, "です"),
			},
			AfterIgnoreConditions: ConvertConditions{
				{Features: pos.SubParEndParticle},
			},
			AppendLongNote:           true,
			EnableKutenToExclamation: true,
			Value:                    "ですわ",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.AuxiliaryVerb, "だ"),
			},
			AfterIgnoreConditions: ConvertConditions{
				{Features: pos.SubParEndParticle},
			},
			AppendLongNote:           true,
			EnableKutenToExclamation: true,
			Value:                    "ですわ",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.VerbIndependence, "する"),
			},
			EnableWhenSentenceSeparation: true,
			AppendLongNote:               true,
			EnableKutenToExclamation:     true,
			Value:                        "いたしますわ",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.VerbIndependence, "なる"),
			},
			EnableWhenSentenceSeparation: true,
			AppendLongNote:               true,
			EnableKutenToExclamation:     true,
			Value:                        "なりますわ",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.VerbIndependence, "ある"),
			},
			Value: "あります",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.SubPostpositionalParticle, "じゃ"),
			},
			Value: "では",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.SubParEndParticle, "か"),
			},
			Value: "の",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.SentenceEndingParticle, "わ"),
			},
			AppendLongNote:           true,
			EnableKutenToExclamation: true,
			Value:                    "ですわ",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.SentenceEndingParticle, "な"),
			},
			Value: "ね",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.SentenceEndingParticle, "さ"),
			},
			Value: "",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.ConnAssistant, "から"),
			},
			Value: "ので",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.ConnAssistant, "けど"),
			},
			Value: "けれど",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.ConnAssistant, "し"),
			},
			Value: "ですし",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.AuxiliaryVerb, "まし"),
			},
			BeforeIgnoreConditions: ConvertConditions{
				{Features: pos.VerbIndependence},
			},
			Value: "おりまし",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.AuxiliaryVerb, "ます"),
			},
			AppendLongNote:           true,
			EnableKutenToExclamation: true,
			Value:                    "ますわ",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.AuxiliaryVerb, "た"),
			},
			EnableWhenSentenceSeparation: true,
			AppendLongNote:               true,
			EnableKutenToExclamation:     true,
			Value:                        "たわ",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.AuxiliaryVerb, "だろ"),
			},
			Value: "でしょう",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.AuxiliaryVerb, "ない"),
			},
			BeforeIgnoreConditions: ConvertConditions{
				{Features: pos.VerbIndependence},
			},
			Value: "ありません",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.VerbNotIndependence, "ください"),
			},
			EnableKutenToExclamation: true,
			Value:                    "くださいまし",
		},
		{
			Conditions: ConvertConditions{
				newCond(pos.VerbNotIndependence, "くれ"),
			},
			EnableKutenToExclamation: true,
			Value:                    "くださいまし",
		},
		newRuleInterjection("ありがとう", "ありがとうございますわ"),
		newRuleInterjection("じゃぁ", "それでは"),
		newRuleInterjection("じゃあ", "それでは"),
		{
			Conditions: ConvertConditions{
				newCond(pos.VerbNotIndependence, "くれる"),
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

func GetMeaningType(typeMap map[MeaningType]ConvertConditions, data tokenizer.TokenData) (MeaningType, bool) {
	for k, cond := range typeMap {
		if cond.MatchAnyTokenData(data) {
			return k, true
		}
	}
	return meaningTypeUnknown, false
}
