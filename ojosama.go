// ojosama はお嬢様言葉への変換ロジックを管理する。
//
//
//
// 関数や構造体のアクセス範囲に関する方針
//
// 基本的にライブラリ用途としては Convert 関数1つのみを公開する。
// Convert関数の挙動の微調整はConvertOption構造体で制御する。
//
// ユーザ側で独自に変換ルールを追加出来たほうが良いかもしれないが、
// パッケージ構成や型名を変更する可能性が高いため、
// 破壊的変更として影響しないように公開する機能を必要最小限にする。
//
// 「ライブラリ呼び出し側で独自に変換ルールを追加したい」という人は
// おそらく現れないと想定しているが、要望がきた時に改めて考える。
//
// それまでは好き勝手パッケージ構成や名前を変更しまくれるように
// 公開範囲を極小にする。
// ゆえにConvert関数内の処理で関数を分割したくなった場合、
// すべてプライベート関数として実装する。
//
// 内部実装の都合でパッケージ等が必要になった場合は、
// すべて internal パッケージ配下にパッケージを切って配置する。
//
//
//
// 変換ルールの実装方針
//
// 変換ルールは internal/converter で定義する。
// 変換ロジックと変換ルールが密結合してるので
// 同じパッケージの方が自然にも感じるけれど分離している。
//
// 理由は、変換ルールの生成にしか使わないプライベート関数や、
// 変換条件の判定関数と、その内部実装用のプライベート関数などが増えてきたため、
// それらのプライベート関数に変換ロジック関数からアクセスさせないためにパッケージを分けている。
package ojosama

import (
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/jiro4989/ojosama/internal/chars"
	"github.com/jiro4989/ojosama/internal/converter"
	"github.com/jiro4989/ojosama/internal/feat"
	"github.com/jiro4989/ojosama/internal/pos"
	"github.com/jiro4989/ojosama/internal/tokendata"
)

// ConvertOption はお嬢様変換時のオプショナルな設定。
type ConvertOption struct {
	// 句点を！に変換する機能をOFFにする。句点を！に変換してしまうと変換元の文章
	// のニュアンスを破壊する可能性があるため、オプションパラメータで無効にでき
	// るようにする。
	DisableKutenToExclamation bool

	forceAppendLongNote     forceAppendLongNote // 単体テスト用のパラメータ
	forceCharsTestMode      *chars.TestMode     // 単体テスト用のパラメータ
	forceKutenToExclamation bool                // KutenToExclamationで強制的に3番目の要素を選択する
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

const (
	politeWord = "ですわ"
)

var (
	alnumRegexp = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

	shuffleElementsKutenToExclamation = []string{"。", "。", "！", "❗"}
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

		// 名詞＋動詞＋終助詞の組み合わせに対して変換する
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
		var kutenToEx bool
		buf, nounKeep, i, kutenToEx = convert(data, tokens, i, buf, nounKeep, opt)

		// 形容詞、自立で文が終わった時は丁寧語ですわを追加する
		if isAppendablePoliteWord(data, tokens, i) {
			buf += politeWord
			kutenToEx = true
		}

		if kutenToEx {
			if ok, s, pos := randomKutenToExclamation(tokens, i, opt); ok {
				buf += s
				i = pos
			}
		}

		result.WriteString(buf)
	}
	return result.String(), nil
}

// convertSentenceEndingParticle は名詞＋動詞（＋助動詞）＋終助詞の組み合わせすべてを満たす場合に変換する。
//
// 終助詞は文の終わりに、文を完結させつつ、文に「希望」「禁止」「詠嘆」「強意」等の意味を添える効果がある。
//
// 例えば「野球しようぜ」の場合、
// 「名詞：野球」「動詞：しよ」「助動詞：う」「終助詞：ぜ」という分解がされる。
//
// 終助詞の「ぜ」としては「希望」の意味合いが含まれるため、希望する意味合いのお嬢様言葉に変換する。
// 例：お野球をいたしませんこと
//
// その他にも「野球するな」だと「お野球をしてはいけませんわ」になる。
func convertSentenceEndingParticle(tokens []tokenizer.Token, tokenPos int) (string, int, bool) {
	for _, r := range converter.SentenceEndingParticleConvertRules {
		var result strings.Builder
		i := tokenPos
		data := tokenizer.NewTokenData(tokens[i])

		// 先頭が一致するならば次の単語に進む
		if !r.Conditions1.MatchAnyTokenData(data) {
			continue
		}
		if len(tokens) <= i+1 {
			continue
		}
		s := data.Surface
		// TODO: ベタ書きしててよくない
		if tokendata.EqualsFeatures(data.Features, pos.NounsGeneral) || tokendata.EqualsFeatures(data.Features[:2], pos.NounsSaDynamic) {
			s = "お" + s
		}
		result.WriteString(s)
		i++
		data = tokenizer.NewTokenData(tokens[i])

		// NOTE:
		// 2つ目以降は value の値で置き換えるため
		// result.WriteString(data.Surface) を実行しない。

		// 2つ目は動詞のいずれかとマッチする。マッチしなければふりだしに戻る
		if !r.Conditions2.MatchAnyTokenData(data) {
			continue
		}
		if len(tokens) <= i+1 {
			continue
		}
		i++
		data = tokenizer.NewTokenData(tokens[i])

		// 助動詞があった場合は無視してトークンを進める。
		// 別に無くても良い。
		if r.AuxiliaryVerb.MatchAllTokenData(data) {
			if len(tokens) <= i+1 {
				continue
			}
			i++
			data = tokenizer.NewTokenData(tokens[i])
		}

		// 最後、終助詞がどの意味分類に該当するかを取得
		mt, ok := converter.GetMeaningType(r.SentenceEndingParticle, data)
		if !ok {
			continue
		}

		// 意味分類に該当する変換候補の文字列を返す
		// TODO: 現状1個だけなので決め打ちで最初の1つ目を返す。
		result.WriteString(r.Value[mt][0])
		return result.String(), i, true
	}
	return "", -1, false
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
	for _, mc := range converter.ContinuousConditionsConvertRules {
		if !matchContinuousConditions(tokens, tokenPos, mc.Conditions) {
			continue
		}

		n := tokenPos + len(mc.Conditions) - 1
		result := mc.Value

		// FIXME: 書き方が汚い
		data := tokenizer.NewTokenData(tokens[tokenPos])
		surface := data.Surface
		if appendablePrefix(data) {
			surface = "お" + surface
		}
		result = strings.ReplaceAll(result, "@1", surface)

		// 句点と～が同時に発生することは無いので早期リターンで良い
		if ok, s, pos := randomKutenToExclamation(tokens, n, opt); ok {
			result += s
			n = pos
			return result, n, true
		}

		if mc.AppendLongNote {
			if note, pos := newLongNote(tokens, n, opt); note != "" {
				result += note
				n = pos
			}
		}

		return result, n, true
	}
	return "", -1, false
}

// matchContinuousConditions は tokens の tokenPos の位置からのトークンが、連続する条件にすべてマッチするかを判定する。
//
// 次のトークンが存在しなかったり、1つでも条件が不一致になった場合 false を返す。
func matchContinuousConditions(tokens []tokenizer.Token, tokenPos int, ccs converter.ConvertConditions) bool {
	j := tokenPos
	for _, conds := range ccs {
		if len(tokens) <= j {
			return false
		}
		data := tokenizer.NewTokenData(tokens[j])
		if !conds.EqualsTokenData(data) {
			return false
		}
		j++
	}
	return true
}

// matchExcludeRule は除外ルールと一致するものが存在するかを判定する。
func matchExcludeRule(data tokenizer.TokenData) bool {
excludeLoop:
	for _, c := range converter.ExcludeRules {
		if !c.Conditions.MatchAllTokenData(data) {
			continue excludeLoop
		}
		return true
	}
	return false
}

// convert は基本的な変換を行う。
func convert(data tokenizer.TokenData, tokens []tokenizer.Token, i int, surface string, nounKeep bool, opt *ConvertOption) (string, bool, int, bool) {
	var ok bool
	var c converter.ConvertRule
	if ok, c = matchConvertRule(data, tokens, i); !ok {
		result := surface
		result, nounKeep = appendPrefix(data, tokens, i, result, nounKeep)
		return result, nounKeep, i, false
	}

	result := c.Value
	pos := i

	// 波線伸ばしをランダムに追加する
	if c.AppendLongNote {
		if note, pos2 := newLongNote(tokens, i, opt); note != "" {
			result += note
			pos = pos2
		}
	}

	// 手前に「お」を付ける
	if !c.DisablePrefix {
		result, nounKeep = appendPrefix(data, tokens, i, result, nounKeep)
	}

	return result, nounKeep, pos, c.EnableKutenToExclamation
}

func matchConvertRule(data tokenizer.TokenData, tokens []tokenizer.Token, i int) (bool, converter.ConvertRule) {
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

	for _, c := range converter.ConvertRules {
		if !c.Conditions.MatchAllTokenData(data) {
			continue
		}

		// 前に続く単語をみて変換を無視する
		if beforeTokenOK && c.BeforeIgnoreConditions.MatchAnyTokenData(beforeToken) {
			break
		}

		// 次に続く単語をみて変換を無視する
		if afterTokenOK && c.AfterIgnoreConditions.MatchAnyTokenData(afterToken) {
			break
		}

		// 文の区切りか、文の終わりの時だけ有効にする。
		// 次のトークンが存在して、且つ次のトークンが文を区切るトークンでない時
		// は変換しない。
		if c.EnableWhenSentenceSeparation && afterTokenOK && !isSentenceSeparation(afterToken) {
			break
		}

		return true, c
	}
	return false, converter.ConvertRule{}
}

func appendablePrefix(data tokenizer.TokenData) bool {
	if !tokendata.EqualsFeatures(data.Features, []string{"名詞", "一般"}) && !tokendata.EqualsFeatures(data.Features[:2], []string{"名詞", "固有名詞"}) {
		return false
	}

	// 丁寧語の場合は「お」を付けない
	if tokendata.IsPoliteWord(data) {
		return false
	}

	return true
}

// appendPrefix は surface の前に「お」を付ける。
func appendPrefix(data tokenizer.TokenData, tokens []tokenizer.Token, i int, surface string, nounKeep bool) (string, bool) {
	if !appendablePrefix(data) {
		return surface, false
	}

	// 次のトークンが動詞の場合は「お」を付けない。
	// 例: プレイする
	if i+1 < len(tokens) {
		data := tokenizer.NewTokenData(tokens[i+1])
		if tokendata.EqualsFeatures(data.Features, []string{"動詞", "自立"}) {
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
		if tokendata.EqualsFeatures(data.Features, []string{"接頭詞", "名詞接続"}) {
			return surface, false
		}

		// サ変接続が来ても付与しない。
		// 例: 横断歩道、解体新書
		if tokendata.EqualsFeatures(data.Features, []string{"名詞", "サ変接続"}) {
			return surface, false
		}
	}

	return "お" + surface, true
}

// isAppendablePoliteWord は丁寧語を追加する。
func isAppendablePoliteWord(data tokenizer.TokenData, tokens []tokenizer.Token, i int) bool {
	if !tokendata.EqualsFeatures(data.Features, []string{"形容詞", "自立"}) {
		return false
	}

	if len(tokens) <= i+1 {
		return false
	}

	// 文の区切りのタイミングでは「ですわ」を差し込む
	if isSentenceSeparation(tokenizer.NewTokenData(tokens[i+1])) {
		return true
	}

	return false
}

// isSentenceSeparation は data が文の区切りに使われる token かどうかを判定する。
func isSentenceSeparation(data tokenizer.TokenData) bool {
	return tokendata.ContainsFeatures([][]string{feat.Kuten, feat.Toten}, data.Features) ||
		tokendata.ContainsString([]string{"！", "!", "？", "?"}, data.Surface)
}

// newLongNote は次の token が感嘆符か疑問符の場合に波線、感嘆符、疑問符をランダムに生成する。
//
// 乱数が絡むと単体テストがやりづらくなるので、 opt を使うことで任意の数付与できるようにしている。
func newLongNote(tokens []tokenizer.Token, i int, opt *ConvertOption) (string, int) {
	var ok bool
	var s string
	if ok, s = creatableLongNote(tokens, i); !ok {
		return "", -1
	}

	var tm *chars.TestMode
	if opt != nil {
		tm = opt.forceCharsTestMode
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

	// ！or？をどれかからランダムに選択する
	feq := chars.SampleExclamationQuestionByValue(s, tm)

	// 次の token は必ず感嘆符か疑問符のどちらかであることが確定しているため
	// -1 して数を調整している。
	for i := 0; i < e-1; i++ {
		suffix.WriteString(feq.Value)
	}

	// 後ろに！や？が連続する場合、それらをすべて feq と同じ種類（半角、全角、
	// 絵文字）の！や？に置き換えて返却する。
	excl, pos := getContinuousExclamationMark(tokens, i, feq)
	suffix.WriteString(excl)
	return suffix.String(), pos
}

func creatableLongNote(tokens []tokenizer.Token, i int) (bool, string) {
	if len(tokens) <= i+1 {
		return false, ""
	}

	data := tokenizer.NewTokenData(tokens[i+1])
	for _, s := range []string{"！", "？", "!", "?"} {
		if data.Surface != s {
			continue
		}
		return true, s
	}
	return false, ""
}

func getContinuousExclamationMark(tokens []tokenizer.Token, i int, feq *chars.ExclamationQuestionMark) (string, int) {
	var result strings.Builder
	pos := i

	for j := i + 1; j < len(tokens); j++ {
		token := tokens[j]
		data := tokenizer.NewTokenData(token)
		for _, r := range data.Surface {
			surface := string(r)
			if ok, eq := chars.IsExclamationQuestionMark(surface); !ok {
				return result.String(), pos
			} else {
				// e は！か？のどちらかなので、同じスタイルの文字を取得して追加
				if got := chars.FindExclamationQuestionByStyleAndMeaning(feq.Style, eq.Meaning); got != nil {

					result.WriteString(got.Value)
				}
			}
		}
		// トークンの位置を制御する変数なので、forループ内では変更しない
		pos = j
	}

	return result.String(), pos
}

// randomKutenToExclamation はランダムで句点を！に変換する。
func randomKutenToExclamation(tokens []tokenizer.Token, tokenPos int, opt *ConvertOption) (bool, string, int) {
	if opt != nil && opt.DisableKutenToExclamation {
		return false, "", tokenPos
	}

	pos := tokenPos + 1
	if len(tokens) <= pos {
		return false, "", tokenPos
	}

	data := tokenizer.NewTokenData(tokens[pos])
	if !tokendata.IsKuten(data) {
		return false, "", tokenPos
	}

	// テスト用に値をすげ替えられるようにする
	var s []string
	if opt != nil && opt.forceKutenToExclamation {
		s = []string{"❗", "❗"}
	} else {
		s = shuffleElementsKutenToExclamation
	}

	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
	return true, s[0], pos
}
