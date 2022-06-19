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
	ConvertTypeSurface convertType = iota + 1
	ConvertTypeFeatures
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
								Type:  ConvertTypeFeatures,
								Value: []string{"名詞", "一般"},
							},
							{
								Type:  ConvertTypeSurface,
								Value: []string{"壱"},
							},
						},
					},
					{
						Conditions: []convertCondition{
							{
								Type:  ConvertTypeFeatures,
								Value: []string{"名詞", "数"},
							},
							{
								Type:  ConvertTypeSurface,
								Value: []string{"百"},
							},
						},
					},
					{
						Conditions: []convertCondition{
							{
								Type:  ConvertTypeFeatures,
								Value: []string{"名詞", "一般"},
							},
							{
								Type:  ConvertTypeSurface,
								Value: []string{"満天"},
							},
						},
					},
					{
						Conditions: []convertCondition{
							{
								Type:  ConvertTypeFeatures,
								Value: []string{"接頭詞", "名詞接続"},
							},
							{
								Type:  ConvertTypeSurface,
								Value: []string{"原"},
							},
						},
					},
					{
						Conditions: []convertCondition{
							{
								Type:  ConvertTypeFeatures,
								Value: []string{"名詞", "一般"},
							},
							{
								Type:  ConvertTypeSurface,
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
								Type:  ConvertTypeFeatures,
								Value: []string{"動詞", "自立"},
							},
							{
								Type:  ConvertTypeSurface,
								Value: []string{"し"},
							},
						},
					},
					{
						Conditions: []convertCondition{
							{
								Type:  ConvertTypeFeatures,
								Value: []string{"助動詞"},
							},
							{
								Type:  ConvertTypeSurface,
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
								Type:  ConvertTypeFeatures,
								Value: []string{"助動詞"},
							},
							{
								Type:  ConvertTypeSurface,
								Value: []string{"だ"},
							},
						},
					},
					{
						Conditions: []convertCondition{
							{
								Type:  ConvertTypeFeatures,
								Value: []string{"助詞", "接続助詞"},
							},
							{
								Type:  ConvertTypeSurface,
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
								Type:  ConvertTypeFeatures,
								Value: []string{"助動詞"},
							},
							{
								Type:  ConvertTypeSurface,
								Value: []string{"な"},
							},
						},
					},
					{
						Conditions: []convertCondition{
							{
								Type:  ConvertTypeFeatures,
								Value: []string{"名詞", "非自立", "一般"},
							},
							{
								Type:  ConvertTypeSurface,
								Value: []string{"ん"},
							},
						},
					},
					{
						Conditions: []convertCondition{
							{
								Type:  ConvertTypeFeatures,
								Value: []string{"助動詞"},
							},
							{
								Type:  ConvertTypeSurface,
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
								Type:  ConvertTypeFeatures,
								Value: []string{"助動詞"},
							},
							{
								Type:  ConvertTypeSurface,
								Value: []string{"だ"},
							},
						},
					},
					{
						Conditions: []convertCondition{
							{
								Type:  ConvertTypeFeatures,
								Value: []string{"助詞", "終助詞"},
							},
							{
								Type:  ConvertTypeSurface,
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
					Type:  ConvertTypeFeatures,
					Value: []string{"名詞", "一般"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"お嬢様"},
				},
			},
		},
	}

	convertRules = []convertRule{
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"名詞", "代名詞", "一般"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"あなた"},
				},
			},
			Value: "貴方",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"名詞", "代名詞", "一般"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"俺"},
				},
			},
			Value: "私",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"名詞", "代名詞", "一般"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"オレ"},
				},
			},
			Value: "ワタクシ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"名詞", "代名詞", "一般"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"おれ"},
				},
			},
			Value: "わたくし",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"名詞", "代名詞", "一般"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"僕"},
				},
			},
			Value: "私",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"名詞", "代名詞", "一般"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"ボク"},
				},
			},
			Value: "ワタクシ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"名詞", "代名詞", "一般"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"ぼく"},
				},
			},
			Value: "わたくし",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"名詞", "代名詞", "一般"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"あたし"},
				},
			},
			Value: "わたくし",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"名詞", "代名詞", "一般"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"わたし"},
				},
			},
			Value: "わたくし",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"名詞", "代名詞", "一般"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"どこ"},
				},
			},
			Value: "どちら",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"名詞", "非自立", "一般"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"もん"},
				},
			},
			Value: "もの",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"です"},
				},
			},
			AfterIgnoreConditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助詞", "副助詞／並立助詞／終助詞"},
				},
			},
			AppendLongNote: true,
			Value:          "ですわ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"だ"},
				},
			},
			AfterIgnoreConditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助詞", "副助詞／並立助詞／終助詞"},
				},
			},
			AppendLongNote: true,
			Value:          "ですわ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"動詞", "自立"},
				},
				{
					Type:  ConvertTypeSurface,
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
					Type:  ConvertTypeFeatures,
					Value: []string{"動詞", "自立"},
				},
				{
					Type:  ConvertTypeSurface,
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
					Type:  ConvertTypeFeatures,
					Value: []string{"動詞", "自立"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"ある"},
				},
			},
			Value: "あります",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助詞", "副助詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"じゃ"},
				},
			},
			Value: "では",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助詞", "副助詞／並立助詞／終助詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"か"},
				},
			},
			Value: "の",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助詞", "終助詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"わ"},
				},
			},
			AppendLongNote: true,
			Value:          "ですわ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助詞", "終助詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"な"},
				},
			},
			Value: "ね",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助詞", "終助詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"さ"},
				},
			},
			Value: "",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助詞", "接続助詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"から"},
				},
			},
			Value: "ので",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助詞", "接続助詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"けど"},
				},
			},
			Value: "けれど",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助詞", "接続助詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"し"},
				},
			},
			Value: "ですし",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"まし"},
				},
			},
			BeforeIgnoreConditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"動詞", "自立"},
				},
			},
			Value: "おりまし",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"ます"},
				},
			},
			AppendLongNote: true,
			Value:          "ますわ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  ConvertTypeSurface,
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
					Type:  ConvertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"だろ"},
				},
			},
			Value: "でしょう",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"ない"},
				},
			},
			BeforeIgnoreConditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"動詞", "自立"},
				},
			},
			Value: "ありません",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"動詞", "非自立"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"ください"},
				},
			},
			Value: "くださいまし",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"動詞", "非自立"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"くれ"},
				},
			},
			Value: "くださいまし",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"感動詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"ありがとう"},
				},
			},
			Value: "ありがとうございますわ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"感動詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"じゃぁ"},
				},
			},
			Value: "それでは",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"感動詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"じゃあ"},
				},
			},
			Value: "それでは",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"動詞", "非自立"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"くれる"},
				},
			},
			Value: "くれます",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"形容詞", "自立"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"汚い"},
				},
			},
			Value: "きったねぇ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"形容詞", "自立"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"きたない"},
				},
			},
			Value: "きったねぇ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"形容詞", "自立"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"臭い"},
				},
			},
			Value: "くっせぇ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"形容詞", "自立"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"くさい"},
				},
			},
			Value: "くっせぇ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"感動詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"うふ"},
				},
			},
			Value: "おほ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"感動詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"うふふ"},
				},
			},
			Value: "おほほ",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"感動詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"う"},
				},
			},
			Value: "お",
		},
		{
			Conditions: []convertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"感動詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"ふふふ"},
				},
			},
			Value: "ほほほ",
		},
	}
)

func (c *convertCondition) equalsTokenData(data tokenizer.TokenData) bool {
	switch c.Type {
	case ConvertTypeFeatures:
		if equalsFeatures(data.Features, c.Value) {
			return true
		}
	case ConvertTypeSurface:
		if data.Surface == c.Value[0] {
			return true
		}
	}
	return false
}

func (c *convertCondition) notEqualsTokenData(data tokenizer.TokenData) bool {
	switch c.Type {
	case ConvertTypeFeatures:
		if !equalsFeatures(data.Features, c.Value) {
			return true
		}
	case ConvertTypeSurface:
		if data.Surface != c.Value[0] {
			return true
		}
	}
	return false
}
