package ojosama

import "github.com/ikawaha/kagome/v2/tokenizer"

type convertRule struct {
	Conditions                   []convertCondition
	BeforeIgnoreConditions       []convertCondition // 前のTokenで条件にマッチした場合は無視する
	AfterIgnoreConditions        []convertCondition // 次のTokenで条件にマッチした場合は無視する
	EnableWhenSentenceSeparation bool               // 文の区切り（単語の後に句点か読点がくる、あるいは何もない）場合だけ有効にする
	AppendLongNote               bool               // 波線を追加する
	Value                        string
}

// FIXME: 型が不適当
type multiConvertRule struct {
	Conditions     [][]convertRule
	AppendLongNote bool
	Value          string
}

type convertType int

type convertCondition struct {
	Type  convertType
	Value []string
}

const (
	convertTypeSurface convertType = iota + 1
	convertTypeFeatures
)

var (
	multiConvertRules = []multiConvertRule{
		{
			Value: "壱百満天原サロメ",
			Conditions: [][]convertRule{
				{
					{
						Conditions: []convertCondition{
							{
								Type:  convertTypeFeatures,
								Value: []string{"名詞", "一般"},
							},
							{
								Type:  convertTypeSurface,
								Value: []string{"壱"},
							},
						},
					},
					{
						Conditions: []convertCondition{
							{
								Type:  convertTypeFeatures,
								Value: []string{"名詞", "数"},
							},
							{
								Type:  convertTypeSurface,
								Value: []string{"百"},
							},
						},
					},
					{
						Conditions: []convertCondition{
							{
								Type:  convertTypeFeatures,
								Value: []string{"名詞", "一般"},
							},
							{
								Type:  convertTypeSurface,
								Value: []string{"満天"},
							},
						},
					},
					{
						Conditions: []convertCondition{
							{
								Type:  convertTypeFeatures,
								Value: []string{"接頭詞", "名詞接続"},
							},
							{
								Type:  convertTypeSurface,
								Value: []string{"原"},
							},
						},
					},
					{
						Conditions: []convertCondition{
							{
								Type:  convertTypeFeatures,
								Value: []string{"名詞", "一般"},
							},
							{
								Type:  convertTypeSurface,
								Value: []string{"サロメ"},
							},
						},
					},
				},
			},
		},

		{
			Value:          "いたしますわ",
			AppendLongNote: true,
			Conditions: [][]convertRule{
				{
					{
						Conditions: []convertCondition{
							{
								Type:  convertTypeFeatures,
								Value: []string{"動詞", "自立"},
							},
							{
								Type:  convertTypeSurface,
								Value: []string{"し"},
							},
						},
					},
					{
						Conditions: []convertCondition{
							{
								Type:  convertTypeFeatures,
								Value: []string{"助動詞"},
							},
							{
								Type:  convertTypeSurface,
								Value: []string{"ます"},
							},
						},
					},
				},
			},
		},

		{
			Value: "ですので",
			Conditions: [][]convertRule{
				{
					{
						Conditions: []convertCondition{
							{
								Type:  convertTypeFeatures,
								Value: []string{"助動詞"},
							},
							{
								Type:  convertTypeSurface,
								Value: []string{"だ"},
							},
						},
					},
					{
						Conditions: []convertCondition{
							{
								Type:  convertTypeFeatures,
								Value: []string{"助詞", "接続助詞"},
							},
							{
								Type:  convertTypeSurface,
								Value: []string{"から"},
							},
						},
					},
				},
			},
		},

		{
			Value: "なんですの",
			Conditions: [][]convertRule{
				{
					{
						Conditions: []convertCondition{
							{
								Type:  convertTypeFeatures,
								Value: []string{"助動詞"},
							},
							{
								Type:  convertTypeSurface,
								Value: []string{"な"},
							},
						},
					},
					{
						Conditions: []convertCondition{
							{
								Type:  convertTypeFeatures,
								Value: []string{"名詞", "非自立", "一般"},
							},
							{
								Type:  convertTypeSurface,
								Value: []string{"ん"},
							},
						},
					},
					{
						Conditions: []convertCondition{
							{
								Type:  convertTypeFeatures,
								Value: []string{"助動詞"},
							},
							{
								Type:  convertTypeSurface,
								Value: []string{"だ"},
							},
						},
					},
				},
			},
		},

		{
			Value: "ですわ",
			Conditions: [][]convertRule{
				{
					{
						Conditions: []convertCondition{
							{
								Type:  convertTypeFeatures,
								Value: []string{"助動詞"},
							},
							{
								Type:  convertTypeSurface,
								Value: []string{"だ"},
							},
						},
					},
					{
						Conditions: []convertCondition{
							{
								Type:  convertTypeFeatures,
								Value: []string{"助詞", "終助詞"},
							},
							{
								Type:  convertTypeSurface,
								Value: []string{"よ"},
							},
						},
					},
				},
			},
		},
	}

	// excludeRules は変換処理を無視するルール。
	// このルールは convertRules よりも優先して評価される。
	excludeRules = []convertRule{
		{
			Conditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"名詞", "一般"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"お嬢様"},
				},
			},
		},
		{
			Conditions: []convertCondition{
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
	}

	convertRules = []convertRule{
		// 名詞、代名詞、一般
		// 三人称
		{
			Conditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"名詞", "代名詞", "一般"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"あなた"},
				},
			},
			Value: "貴方",
		},

		// 一人称
		newCondPronounGeneral("俺", "私"),
		newCondPronounGeneral("オレ", "ワタクシ"),
		newCondPronounGeneral("おれ", "わたくし"),
		newCondPronounGeneral("僕", "私"),
		newCondPronounGeneral("ボク", "ワタクシ"),
		newCondPronounGeneral("ぼく", "わたくし"),
		newCondPronounGeneral("あたし", "わたくし"),
		newCondPronounGeneral("わたし", "わたくし"),
		// TODO: AfterIgnore系も簡単に定義できるようにしたい
		{
			Conditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"名詞", "一般"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"パパ"},
				},
			},
			AfterIgnoreConditions: []convertCondition{
				{
					Type:  convertTypeSurface,
					Value: []string{"上"},
				},
			},
			Value: "パパ上",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"名詞", "一般"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"ママ"},
				},
			},
			AfterIgnoreConditions: []convertCondition{
				{
					Type:  convertTypeSurface,
					Value: []string{"上"},
				},
			},
			Value: "ママ上",
		},

		// こそあど言葉
		newCondPronounGeneral("これ", "こちら"),
		newCondPronounGeneral("それ", "そちら"),
		newCondPronounGeneral("あれ", "あちら"),
		newCondPronounGeneral("どれ", "どちら"),
		newCondAdnominalAdjective("この", "こちらの"),
		newCondAdnominalAdjective("その", "そちらの"),
		newCondAdnominalAdjective("あの", "あちらの"),
		newCondAdnominalAdjective("どの", "どちらの"),
		newCondPronounGeneral("ここ", "こちら"),
		newCondPronounGeneral("そこ", "そちら"),
		newCondPronounGeneral("あそこ", "あちら"),
		newCondPronounGeneral("どこ", "どちら"),
		newCondAdnominalAdjective("こんな", "このような"),
		newCondAdnominalAdjective("そんな", "そのような"),
		newCondAdnominalAdjective("あんな", "あのような"),
		newCondAdnominalAdjective("どんな", "どのような"),

		{
			Conditions: []convertCondition{
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
			Conditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"です"},
				},
			},
			AfterIgnoreConditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助詞", "副助詞／並立助詞／終助詞"},
				},
			},
			AppendLongNote: true,
			Value:          "ですわ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"だ"},
				},
			},
			AfterIgnoreConditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助詞", "副助詞／並立助詞／終助詞"},
				},
			},
			AppendLongNote: true,
			Value:          "ですわ",
		},
		{
			Conditions: []convertCondition{
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
			Conditions: []convertCondition{
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
			Conditions: []convertCondition{
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
			Conditions: []convertCondition{
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
			Conditions: []convertCondition{
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
			Conditions: []convertCondition{
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
			Conditions: []convertCondition{
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
			Conditions: []convertCondition{
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
			Conditions: []convertCondition{
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
			Conditions: []convertCondition{
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
			Conditions: []convertCondition{
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
			Conditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"まし"},
				},
			},
			BeforeIgnoreConditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"動詞", "自立"},
				},
			},
			Value: "おりまし",
		},
		{
			Conditions: []convertCondition{
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
			Conditions: []convertCondition{
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
			Conditions: []convertCondition{
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
			Conditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"ない"},
				},
			},
			BeforeIgnoreConditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"動詞", "自立"},
				},
			},
			Value: "ありません",
		},
		{
			Conditions: []convertCondition{
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
			Conditions: []convertCondition{
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
		{
			Conditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"感動詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"ありがとう"},
				},
			},
			Value: "ありがとうございますわ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"感動詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"じゃぁ"},
				},
			},
			Value: "それでは",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"感動詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"じゃあ"},
				},
			},
			Value: "それでは",
		},
		{
			Conditions: []convertCondition{
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
		{
			Conditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"形容詞", "自立"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"汚い"},
				},
			},
			Value: "きったねぇ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"形容詞", "自立"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"きたない"},
				},
			},
			Value: "きったねぇ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"形容詞", "自立"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"臭い"},
				},
			},
			Value: "くっせぇ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"形容詞", "自立"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"くさい"},
				},
			},
			Value: "くっせぇ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"感動詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"うふ"},
				},
			},
			Value: "おほ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"感動詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"うふふ"},
				},
			},
			Value: "おほほ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"感動詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"う"},
				},
			},
			Value: "お",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  convertTypeFeatures,
					Value: []string{"感動詞"},
				},
				{
					Type:  convertTypeSurface,
					Value: []string{"ふふふ"},
				},
			},
			Value: "ほほほ",
		},
	}
)

func (c *convertCondition) equalsTokenData(data tokenizer.TokenData) bool {
	switch c.Type {
	case convertTypeFeatures:
		if equalsFeatures(data.Features, c.Value) {
			return true
		}
	case convertTypeSurface:
		if data.Surface == c.Value[0] {
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
		if data.Surface != c.Value[0] {
			return true
		}
	}
	return false
}

var (
	pronounGeneral     = []string{"名詞", "代名詞", "一般"}
	nounsGeneral       = []string{"名詞", "一般"}
	adnominalAdjective = []string{"連体詞"}
)

func newCond(features []string, surface, value string) convertRule {
	return convertRule{
		Conditions: []convertCondition{
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

func newCondPronounGeneral(surface, value string) convertRule {
	return newCond(pronounGeneral, surface, value)
}

func newCondNounsGeneral(surface, value string) convertRule {
	return newCond(nounsGeneral, surface, value)
}

func newCondAdnominalAdjective(surface, value string) convertRule {
	return newCond(adnominalAdjective, surface, value)
}
