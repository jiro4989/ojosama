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
		if alnumRegexp.MatchString(data.Surface) {
			result.WriteString(buf)
			continue
		}

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
		buf = convert(data, tokens, i, buf, opt)

		// 名詞 一般の場合は手前に「お」をつける。
		// ただし、直後に動詞 自立がくるときは「お」をつけない。
		// 例: プレイする
		buf, nounKeep = appendPrefix(data, tokens, i, buf, nounKeep)

		// 形容詞、自立で文が終わった時は丁寧語ですわを追加する
		buf = appendPoliteWord(data, tokens, i, buf)

		result.WriteString(buf)
	}
	return result.String(), nil
}

// convertContinuousConditions は連続する条件がすべてマッチした時に変換する。
//
// 例えば「壱百満天原サロメ」や「横断歩道」のように、複数のTokenがこの順序で連続
// して初めて1つの意味になるような条件をすべて満たした時に変換を行う。
func convertContinuousConditions(tokens []tokenizer.Token, i int, opt *ConvertOption) (string, int, bool) {
ruleLoop:
	for _, mc := range continuousConditionsConvertRules {
		j := i
		var s strings.Builder

		// conditionsのすべての評価がtrueの場合だけ変換する。
		// マッチすると次のTokenにアクセスするために、ループカウンタを1進める。
		//
		// conditions が1つでも不一致の場合、
		// そのruleの以降のconditionsは評価する意味が無い。
		// よって conditions の評価ループを脱出し、次のruleの評価に移行する。
		//
		// 次のトークンが存在しない場合も評価する意味が無いので conditions 評価
		// ループを脱出する。
		for _, conds := range mc.Conditions {
			if len(tokens) <= j {
				continue ruleLoop
			}
			data := tokenizer.NewTokenData(tokens[j])
			if !conds.matchAllTokenData(data) {
				continue ruleLoop
			}
			s.WriteString(data.Surface)
			j++
		}

		result := mc.Value
		if mc.AppendLongNote {
			result = appendLongNote(result, tokens, j-1, opt)
		}
		return result, j - 1, true
	}
	return "", -1, false
}

// matchExcludeRule は除外ルールと一致するものが存在するかを判定する。
func matchExcludeRule(data tokenizer.TokenData) bool {
excludeLoop:
	for _, c := range excludeRules {
		for _, cond := range c.Conditions {
			if cond.notEqualsTokenData(data) {
				continue excludeLoop
			}
		}
		return true
	}
	return false
}

func convert(data tokenizer.TokenData, tokens []tokenizer.Token, i int, surface string, opt *ConvertOption) string {
converterLoop:
	for _, c := range convertRules {
		for _, cond := range c.Conditions {
			if cond.notEqualsTokenData(data) {
				continue converterLoop
			}
		}

		// 前に続く単語をみて変換を無視する
		for _, cond := range c.BeforeIgnoreConditions {
			if 0 < i {
				data := tokenizer.NewTokenData(tokens[i-1])
				if cond.equalsTokenData(data) {
					break converterLoop
				}
			}
		}

		// 次に続く単語をみて変換を無視する
		for _, cond := range c.AfterIgnoreConditions {
			if i+1 < len(tokens) {
				data := tokenizer.NewTokenData(tokens[i+1])
				if cond.equalsTokenData(data) {
					break converterLoop
				}
			}
		}

		// 文の区切りか、文の終わりの時だけ有効にする。
		if c.EnableWhenSentenceSeparation {
			if i+1 < len(tokens) {
				// 次のトークンが区切りではない場合は変換しない
				data := tokenizer.NewTokenData(tokens[i+1])
				if !isSentenceSeparation(data) {
					break
				}
			}
			// 次のトークンが存在しない場合は文の終わりなので変換する
		}

		result := c.Value

		// 波線伸ばしをランダムに追加する
		if c.AppendLongNote {
			result = appendLongNote(result, tokens, i, opt)
		}

		return result
	}
	return surface
}

func appendPrefix(data tokenizer.TokenData, tokens []tokenizer.Token, i int, surface string, nounKeep bool) (string, bool) {
	if equalsFeatures(data.Features, []string{"名詞", "一般"}) || equalsFeatures(data.Features[:2], []string{"名詞", "固有名詞"}) {
		if i+1 < len(tokens) {
			data := tokenizer.NewTokenData(tokens[i+1])
			if equalsFeatures(data.Features, []string{"動詞", "自立"}) {
				return surface, nounKeep
			}
		}

		// 名詞 一般が連続する場合は最初の1つ目にだけ「お」を付ける
		if !nounKeep {
			// 1つ手前にすでに「お」が付いている場合は付与しない
			if 0 < i {
				data := tokenizer.NewTokenData(tokens[i-1])
				if equalsFeatures(data.Features, []string{"接頭詞", "名詞接続"}) {
					return surface, false
				}
			}
			return "お" + surface, true
		}
	}
	return surface, false
}

// 丁寧語を差し込む
func appendPoliteWord(data tokenizer.TokenData, tokens []tokenizer.Token, i int, surface string) string {
	if equalsFeatures(data.Features, []string{"形容詞", "自立"}) {
		if i+1 < len(tokens) {
			// 文の区切りのタイミングでは「ですわ」を差し込む
			data := tokenizer.NewTokenData(tokens[i+1])
			if isSentenceSeparation(data) {
				return surface + "ですわ"
			}

			// // 次の単語が助動詞（です）出ない場合は「です」を差し込む
			// if !equalsFeatures(data.Features, []string{"助動詞"}) {
			// 	return surface + "です"
			// }
		}
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
	if i+1 < len(tokens) {
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

			// 次の token は必ず感嘆符か疑問符のどちらかであることが確定しるため
			// -1 して数を調整している。
			for i := 0; i < e-1; i++ {
				suffix.WriteString(s)
			}

			src += suffix.String()
			break
		}
	}
	return src
}
