package ojosama

import (
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// ConvertOption はお嬢様変換時のオプショナルな設定。
type ConvertOption struct {
	forceAppendLongNote forceAppendLongNote // 単体テスト用のパラメータ
}

// forceAppendLongNote は強制的に波線や感嘆符や疑問符を任意の数追加するための設定。
//
// 波線や感嘆符の付与には乱数が絡むため、単体テスト実行時に確実に等しい結果を得
// ることが難しい。この問題を回避するために、このパラメータを差し込むことで乱数
// の影響を受けないように制御する。単体テストでしか使うことを想定していないため、
// パブリックにはしない。
type forceAppendLongNote struct {
	enable               bool
	wavyLineCount        int
	exclamationMarkCount int
}

var (
	alnumRegexp = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Convert はテキストを壱百満天原サロメお嬢様風の口調に変換して返却する。
//
// 簡単に説明すると「ハーブですわ！」を「おハーブですわ～～！！！」と変換する。
// それ以外にもいくつかバリエーションがある。
//
// opt は挙動を微調整するためのオプショナルなパラメータ。不要であれば nil を渡せ
// ば良い。
func Convert(src string, opt *ConvertOption) (string, error) {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return "", err
	}

	// tokenize
	tokens := t.Tokenize(src)
	var result strings.Builder
	var nounKeep bool
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		data := tokenizer.NewTokenData(token)
		buf := data.Surface

		// 英数字のみの単語の場合は何もしない
		if alnumRegexp.MatchString(buf) {
			result.WriteString(buf)
			continue
		}

		// 動詞＋する＋終助詞の組み合わせに対して変換する
		if s, n, ok := convertSentenceEndingParticle(tokens, i); ok {
			i = n
			result.WriteString(s)
			continue
		}

		// 連続する条件による変換を行う
		if s, n, ok := convertContinuousConditions(tokens, i, opt); ok {
			i = n
			result.WriteString(s)
			continue
		}

		// 特定条件は優先して無視する
		if matchExcludeRule(data) {
			result.WriteString(buf)
			continue
		}

		// お嬢様言葉に変換
		buf, nounKeep = convert(data, tokens, i, buf, nounKeep, opt)

		// 形容詞、自立で文が終わった時は丁寧語ですわを追加する
		buf = appendPoliteWord(data, tokens, i, buf)

		result.WriteString(buf)
	}
	return result.String(), nil
}

func convertSentenceEndingParticle(tokens []tokenizer.Token, tokenPos int) (string, int, bool) {
	var result strings.Builder
	for _, r := range sentenceEndingParticleConvertRules {
		i := tokenPos
		data := tokenizer.NewTokenData(tokens[i])

		// 先頭が一致するならば次の単語に進む
		if !matchAnyMultiConvertConditions(r.conditions1, data) {
			continue
		}
		if len(tokens) <= i+1 {
			continue
		}
		result.WriteString(data.Surface)
		i++
		data = tokenizer.NewTokenData(tokens[i])

		// NOTE:
		// 2つ目以降は value の値で置き換えるため
		// result.WriteString(data.Surface) を実行しない。

		// 2つ目は動詞のいずれかとマッチする。マッチしなければふりだしに戻る
		if !matchAnyMultiConvertConditions(r.conditions2, data) {
			continue
		}
		if len(tokens) <= i+1 {
			continue
		}
		i++
		data = tokenizer.NewTokenData(tokens[i])

		// 助動詞はあった場合は無視する
		// なくても良い。
		if r.auxiliaryVerb.matchAllTokenData(data) {
			if len(tokens) <= i+1 {
				continue
			}
			i++
			data = tokenizer.NewTokenData(tokens[i])
		}

		// 最後、終助詞がどの意味分類に該当するかを取得
		mt, ok := getMeaningType(r.sentenceEndingParticle, data)
		if !ok {
			continue
		}

		// 意味分類に該当する変換候補の文字列を返す
		// TODO: 現状1個だけなので決め打ちで最初の1つ目を返す。
		result.WriteString(r.value[mt][0])
		return result.String(), i, true
	}
	return "", -1, false
}

func getMeaningType(typeMap map[meaningType][]convertConditions, data tokenizer.TokenData) (meaningType, bool) {
	for k, v := range typeMap {
		for _, cond := range v {
			if cond.matchAllTokenData(data) {
				return k, true
			}
		}
	}
	return 0, false
}

// convertContinuousConditions は連続する条件による変換ルールにマッチした変換結果を返す。
//
// 例えば「壱百満天原サロメ」や「横断歩道」のように、複数のTokenがこの順序で連続
// して初めて1つの意味になるような条件をすべて満たした時に結果を返す。
//
// 連続する条件にマッチした場合は tokenPos をその分だけ進める必要があるため、進
// めた後の tokenPos を返却する。
//
// 第三引数は変換ルールにマッチしたかどうかを返す。
func convertContinuousConditions(tokens []tokenizer.Token, tokenPos int, opt *ConvertOption) (string, int, bool) {
	for _, mc := range continuousConditionsConvertRules {
		if !matchContinuousConditions(tokens, tokenPos, mc.Conditions) {
			continue
		}

		n := tokenPos + len(mc.Conditions) - 1
		result := mc.Value
		if mc.AppendLongNote {
			result = appendLongNote(result, tokens, n, opt)
		}
		return result, n, true
	}
	return "", -1, false
}

// matchContinuousConditions は tokens の tokenPos の位置からのトークンが、連続する条件にすべてマッチするかを判定する。
//
// 次のトークンが存在しなかったり、1つでも条件が不一致になった場合 false を返す。
func matchContinuousConditions(tokens []tokenizer.Token, tokenPos int, ccs []convertConditions) bool {
	j := tokenPos
	for _, conds := range ccs {
		if len(tokens) <= j {
			return false
		}
		data := tokenizer.NewTokenData(tokens[j])
		if !conds.matchAllTokenData(data) {
			return false
		}
		j++
	}
	return true
}

// matchExcludeRule は除外ルールと一致するものが存在するかを判定する。
func matchExcludeRule(data tokenizer.TokenData) bool {
excludeLoop:
	for _, c := range excludeRules {
		if !c.Conditions.matchAllTokenData(data) {
			continue excludeLoop
		}
		return true
	}
	return false
}

// convert は基本的な変換を行う。
func convert(data tokenizer.TokenData, tokens []tokenizer.Token, i int, surface string, nounKeep bool, opt *ConvertOption) (string, bool) {
	var beforeToken tokenizer.TokenData
	var beforeTokenOK bool
	if 0 < i {
		beforeToken = tokenizer.NewTokenData(tokens[i-1])
		beforeTokenOK = true
	}

	var afterToken tokenizer.TokenData
	var afterTokenOK bool
	if i+1 < len(tokens) {
		afterToken = tokenizer.NewTokenData(tokens[i+1])
		afterTokenOK = true
	}

	for _, c := range convertRules {
		if !c.Conditions.matchAllTokenData(data) {
			continue
		}

		// 前に続く単語をみて変換を無視する
		if beforeTokenOK && c.BeforeIgnoreConditions.matchAnyTokenData(beforeToken) {
			break
		}

		// 次に続く単語をみて変換を無視する
		if afterTokenOK && c.AfterIgnoreConditions.matchAnyTokenData(afterToken) {
			break
		}

		// 文の区切りか、文の終わりの時だけ有効にする。
		// 次のトークンが存在して、且つ次のトークンが文を区切るトークンでない時
		// は変換しない。
		if c.EnableWhenSentenceSeparation && afterTokenOK && !isSentenceSeparation(afterToken) {
			break
		}

		result := c.Value

		// 波線伸ばしをランダムに追加する
		if c.AppendLongNote {
			result = appendLongNote(result, tokens, i, opt)
		}

		// 手前に「お」を付ける
		if !c.DisablePrefix {
			result, nounKeep = appendPrefix(data, tokens, i, result, nounKeep)
		}

		return result, nounKeep
	}

	// 手前に「お」を付ける
	result := surface
	result, nounKeep = appendPrefix(data, tokens, i, result, nounKeep)
	return result, nounKeep
}

// appendPrefix は surface の前に「お」を付ける。
func appendPrefix(data tokenizer.TokenData, tokens []tokenizer.Token, i int, surface string, nounKeep bool) (string, bool) {
	if !equalsFeatures(data.Features, []string{"名詞", "一般"}) && !equalsFeatures(data.Features[:2], []string{"名詞", "固有名詞"}) {
		return surface, false
	}

	// 丁寧語の場合は「お」を付けない
	if isPoliteWord(data) {
		return surface, false
	}

	// 次のトークンが動詞の場合は「お」を付けない。
	// 例: プレイする
	if i+1 < len(tokens) {
		data := tokenizer.NewTokenData(tokens[i+1])
		if equalsFeatures(data.Features, []string{"動詞", "自立"}) {
			return surface, nounKeep
		}
	}

	// すでに「お」を付与されているので、「お」を付与しない
	if nounKeep {
		return surface, false
	}

	if 0 < i {
		data := tokenizer.NewTokenData(tokens[i-1])

		// 手前のトークンが「お」の場合は付与しない
		if equalsFeatures(data.Features, []string{"接頭詞", "名詞接続"}) {
			return surface, false
		}

		// サ変接続が来ても付与しない。
		// 例: 横断歩道、解体新書
		if equalsFeatures(data.Features, []string{"名詞", "サ変接続"}) {
			return surface, false
		}
	}

	return "お" + surface, true
}

// appendPoliteWord は丁寧語を追加する。
func appendPoliteWord(data tokenizer.TokenData, tokens []tokenizer.Token, i int, surface string) string {
	if !equalsFeatures(data.Features, []string{"形容詞", "自立"}) {
		return surface
	}

	if len(tokens) <= i+1 {
		return surface
	}

	// 文の区切りのタイミングでは「ですわ」を差し込む
	if isSentenceSeparation(tokenizer.NewTokenData(tokens[i+1])) {
		return surface + "ですわ"
	}

	return surface
}

// isSentenceSeparation は data が文の区切りに使われる token かどうかを判定する。
func isSentenceSeparation(data tokenizer.TokenData) bool {
	return containsFeatures([][]string{{"記号", "句点"}, {"記号", "読点"}}, data.Features) ||
		containsString([]string{"！", "!", "？", "?"}, data.Surface)
}

// appendLongNote は次の token が感嘆符か疑問符の場合に波線、感嘆符、疑問符をランダムに追加する。
//
// 乱数が絡むと単体テストがやりづらくなるので、 opt を使うことで任意の数付与でき
// るようにしている。
func appendLongNote(src string, tokens []tokenizer.Token, i int, opt *ConvertOption) string {
	if len(tokens) <= i+1 {
		return src
	}

	data := tokenizer.NewTokenData(tokens[i+1])
	for _, s := range []string{"！", "？"} {
		if data.Surface != s {
			continue
		}

		var (
			w, e int
		)
		if opt != nil && opt.forceAppendLongNote.enable {
			// opt がある場合に限って任意の数付与できる。基本的に単体テスト用途。
			w = opt.forceAppendLongNote.wavyLineCount
			e = opt.forceAppendLongNote.exclamationMarkCount
		} else {
			w = rand.Intn(3)
			e = rand.Intn(3)
		}

		var suffix strings.Builder
		for i := 0; i < w; i++ {
			suffix.WriteString("～")
		}

		// 次の token は必ず感嘆符か疑問符のどちらかであることが確定しているため
		// -1 して数を調整している。
		for i := 0; i < e-1; i++ {
			suffix.WriteString(s)
		}

		src += suffix.String()
		break
	}
	return src
}

// isPoliteWord は丁寧語かどうかを判定する。
// 読みがオで始まる言葉も true になる。
func isPoliteWord(data tokenizer.TokenData) bool {
	return strings.HasPrefix(data.Reading, "オ")
}
