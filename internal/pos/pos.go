// pos は PartOfSpeech (品詞)パッケージ。
//
// 英語 文法 品詞
// https://ja.wikibooks.org/wiki/%E8%8B%B1%E8%AA%9E/%E6%96%87%E6%B3%95/%E5%93%81%E8%A9%9E
package pos

var (
	PronounGeneral            = []string{"名詞", "代名詞", "一般"}
	NounsGeneral              = []string{"名詞", "一般"}
	SpecificGeneral           = []string{"名詞", "固有名詞", "一般"}
	NotIndependenceGeneral    = []string{"名詞", "非自立", "一般"}
	AdnominalAdjective        = []string{"連体詞"}
	AdjectivesSelfSupporting  = []string{"形容詞", "自立"}
	Interjection              = []string{"感動詞"}
	VerbIndependence          = []string{"動詞", "自立"}
	VerbNotIndependence       = []string{"動詞", "非自立"}
	SentenceEndingParticle    = []string{"助詞", "終助詞"}
	SubPostpositionalParticle = []string{"助詞", "副助詞"}
	AssistantParallelParticle = []string{"助詞", "並立助詞"}
	SubParEndParticle         = []string{"助詞", "副助詞／並立助詞／終助詞"}
	ConnAssistant             = []string{"助詞", "接続助詞"}
	AuxiliaryVerb             = []string{"助動詞"}
	NounsSaDynamic            = []string{"名詞", "サ変接続"}
)
