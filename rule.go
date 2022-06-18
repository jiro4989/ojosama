package ojosama

import "github.com/ikawaha/kagome/v2/tokenizer"

type Converter struct {
	Conditions                   []ConvertCondition
	BeforeIgnoreConditions       []ConvertCondition // 前のTokenで条件にマッチした場合は無視する
	AfterIgnoreConditions        []ConvertCondition // 次のTokenで条件にマッチした場合は無視する
	EnableWhenSentenceSeparation bool               // 文の区切り（単語の後に句点か読点がくる、あるいは何もない）場合だけ有効にする
	AppendLongNote               bool               // 波線を追加する
	Value                        string
}

// FIXME: 型が不適当
type MultiConverter struct {
	Conditions     [][]Converter
	AppendLongNote bool
	Value          string
}

type ConvertType int

type ConvertCondition struct {
	Type  ConvertType
	Value []string
}

const (
	ConvertTypeSurface ConvertType = iota + 1
	ConvertTypeReading
	ConvertTypeFeatures
)

var (
	multiConvertRules = []MultiConverter{
		{
			Value: "壱百満天原サロメ",
			Conditions: [][]Converter{
				{
					{
						Conditions: []ConvertCondition{
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
						Conditions: []ConvertCondition{
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
						Conditions: []ConvertCondition{
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
						Conditions: []ConvertCondition{
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
						Conditions: []ConvertCondition{
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
			Conditions: [][]Converter{
				{
					{
						Conditions: []ConvertCondition{
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
						Conditions: []ConvertCondition{
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
			Conditions: [][]Converter{
				{
					{
						Conditions: []ConvertCondition{
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
						Conditions: []ConvertCondition{
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
			Conditions: [][]Converter{
				{
					{
						Conditions: []ConvertCondition{
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
						Conditions: []ConvertCondition{
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
						Conditions: []ConvertCondition{
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
			Conditions: [][]Converter{
				{
					{
						Conditions: []ConvertCondition{
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
						Conditions: []ConvertCondition{
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
	excludeRules = []Converter{
		{
			Conditions: []ConvertCondition{
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

	convertRules = []Converter{
		{
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"です"},
				},
			},
			AfterIgnoreConditions: []ConvertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助詞", "副助詞／並立助詞／終助詞"},
				},
			},
			AppendLongNote: true,
			Value:          "ですわ",
		},
		{
			Conditions: []ConvertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"だ"},
				},
			},
			AfterIgnoreConditions: []ConvertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助詞", "副助詞／並立助詞／終助詞"},
				},
			},
			AppendLongNote: true,
			Value:          "ですわ",
		},
		{
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"まし"},
				},
			},
			BeforeIgnoreConditions: []ConvertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"動詞", "自立"},
				},
			},
			Value: "おりまし",
		},
		{
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"助動詞"},
				},
				{
					Type:  ConvertTypeSurface,
					Value: []string{"ない"},
				},
			},
			BeforeIgnoreConditions: []ConvertCondition{
				{
					Type:  ConvertTypeFeatures,
					Value: []string{"動詞", "自立"},
				},
			},
			Value: "ありません",
		},
		{
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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
			Conditions: []ConvertCondition{
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

func (c *ConvertCondition) equalsTokenData(data tokenizer.TokenData) bool {
	switch c.Type {
	case ConvertTypeFeatures:
		if equalsFeatures(data.Features, c.Value) {
			return true
		}
	case ConvertTypeSurface:
		if data.Surface == c.Value[0] {
			return true
		}
	case ConvertTypeReading:
		if data.Reading == c.Value[0] {
			return true
		}
	}
	return false
}

func (c *ConvertCondition) notEqualsTokenData(data tokenizer.TokenData) bool {
	switch c.Type {
	case ConvertTypeFeatures:
		if !equalsFeatures(data.Features, c.Value) {
			return true
		}
	case ConvertTypeSurface:
		if data.Surface != c.Value[0] {
			return true
		}
	case ConvertTypeReading:
		if data.Reading != c.Value[0] {
			return true
		}
	}
	return false
}
