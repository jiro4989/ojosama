package ojosama

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		desc    string
		src     string
		opt     *ConvertOption
		want    string
		wantErr bool
	}{
		{
			desc:    "正常系: 名詞の手前には「お」をお付けいたしますわ",
			src:     "これはハーブです",
			want:    "こちらはおハーブですわ",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: 「名詞,一般」の後に「動詞,自立」が来た時は2Tokenで1つの動詞として解釈いたしますので、「お」を付けませんわ",
			src:     "プレイする",
			want:    "プレイいたしますわ",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: 「プレイする」の例文ですわ",
			src:     "〇〇をプレイする",
			want:    "〇〇をプレイいたしますわ",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: あなたは貴方に変換いたしますわ",
			src:     "あなた",
			want:    "貴方",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: わたしはわたくしに変換いたしますわ",
			src:     "わたし",
			want:    "わたくし",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: 「ですか」のときは「ですの」に変換いたしますわ",
			src:     "ビデオテープはどこで使うんですか",
			want:    "おビデオテープはどちらで使うんですの",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: 「汚い|きたない」のときは「きったねぇ」に変換いたしますわ",
			src:     "汚いです",
			want:    "きったねぇですわ",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: 「汚い|きたない」のときは「きったねぇ」に変換いたしますわ",
			src:     "きたないです",
			want:    "きったねぇですわ",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: 「臭い|くさい」のときは「くっせぇ」に変換いたしますわ",
			src:     "臭いです",
			want:    "くっせぇですわ",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: 「臭い|くさい」のときは「くっせぇ」に変換いたしますわ",
			src:     "くさいです",
			want:    "くっせぇですわ",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: 複数の文も処理できますわ",
			src:     "〇〇をプレイする。きたないわ。",
			want:    "〇〇をプレイいたしますわ。きったねぇですわ。",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: 「ました」は「おりました」ですの",
			src:     "わたしも使ってました",
			want:    "わたくしも使っておりましたわ",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: アルファベット単語の場合は「お」をつけませんの",
			src:     "これはgrassです。あれはabcdefg12345です",
			want:    "こちらはgrassですわ。あちらはabcdefg12345ですわ",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: 「します」は「いたしますわ」に変換しますわ",
			src:     "使用します",
			want:    "使用いたしますわ",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: 名詞が連続する場合は最初の1つ目にだけ「お」を付けますわ",
			src:     "一般女性。経年劣化。トップシークレット",
			want:    "お一般女性。お経年劣化。おトップシークレット",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: お嬢様は変更しませんわ",
			src:     "お嬢様",
			want:    "お嬢様",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: サロメお嬢様は固有名詞ですわ",
			src:     "わたしは壱百満天原サロメです",
			want:    "わたくしは壱百満天原サロメですわ",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: 文の途中の「する」の場合は変換しませんわ",
			src:     "テキストファイルをまるごと変換する場合は、以下のように実行します。",
			want:    "おテキストファイルをまるごと変換する場合は、以下のように実行いたしますわ。",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: 文の途中の「する」の場合は変換しませんわ",
			src:     "以下のお嬢様を可能な限り再現するアルゴリズムを目指してます。",
			want:    "以下のお嬢様を可能な限り再現するおアルゴリズムを目指してますわ。",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: すでに直前に「お」が付与されている場合は、「お」を付与しませんわ",
			src:     "お寿司。おハーブ",
			want:    "お寿司。おハーブ",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: 「ください」は「くださいまし」に変換いたしますわ",
			src:     "お使いください",
			want:    "お使いくださいまし",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: 「ください」は「くださいまし」に変換いたしますわ",
			src:     "お使いください！",
			want:    "お使いくださいまし！",
			opt:     nil,
			wantErr: false,
		},
		{
			desc: "正常系: 波線のばしを付与する場合がありますわ",
			src:  "これはハーブです！",
			want: "こちらはおハーブですわ～～！！！",
			opt: &ConvertOption{
				forceAppendLongNote: forceAppendLongNote{
					enable:               true,
					wavyLineCount:        2,
					exclamationMarkCount: 3,
				},
			},
			wantErr: false,
		},
		{
			desc: "正常系: 波線のばしを付与する場合がありますわ",
			src:  "プレイします！",
			want: "プレイいたしますわ～～！！！",
			opt: &ConvertOption{
				forceAppendLongNote: forceAppendLongNote{
					enable:               true,
					wavyLineCount:        2,
					exclamationMarkCount: 3,
				},
			},
			wantErr: false,
		},
		{
			desc: "正常系: 波線のばしを付与する場合がありますわ",
			src:  "プレイする！",
			want: "プレイいたしますわ～～！！！",
			opt: &ConvertOption{
				forceAppendLongNote: forceAppendLongNote{
					enable:               true,
					wavyLineCount:        2,
					exclamationMarkCount: 3,
				},
			},
			wantErr: false,
		},

		{
			desc:    "正常系: （動詞）ないはそのままですわ",
			src:     "限らない、飾らない、数えない、くだらない",
			want:    "限らない、飾らない、数えない、くだらない",
			opt:     nil,
			wantErr: false,
		},

		{
			desc: "正常系: ありましたはありましたのままですわ",
			src:  "ハーブがありました！",
			want: "おハーブがありましたわ～～！！！",
			opt: &ConvertOption{
				forceAppendLongNote: forceAppendLongNote{
					enable:               true,
					wavyLineCount:        2,
					exclamationMarkCount: 3,
				},
			},
			wantErr: false,
		},

		{
			desc:    "正常系: 一人称はすべて「私」に変換いたしますわ",
			src:     "俺は。オレは。おれは。僕は。ボクは。ぼくは。あたしは。わたしは。",
			want:    "私は。ワタクシは。わたくしは。私は。ワタクシは。わたくしは。わたくしは。わたくしは。",
			opt:     nil,
			wantErr: false,
		},

		// FIXME: 話しますわね、にしたいけれど「話す」で1単語と判定されているの
		// で変換が難しい
		{
			desc:    "正常系: 小説の一文をテストしますわ",
			src:     "俺の体験した怖い話、聞いてくれるか？ありがとう。じゃぁ、話すよ。",
			want:    "私の体験した怖い話、聞いてくれますの？ありがとうございますわ。それでは、話すよ。",
			opt:     nil,
			wantErr: false,
		},

		{
			desc:    "正常系: 小説の一文をテストしますわ",
			src:     "俺の転職前の会社の頃の話だから、もう3年も前の話になる。",
			want:    "私の転職前のお会社の頃の話ですので、もう3年も前の話になりますわ。",
			opt:     nil,
			wantErr: false,
		},

		// FIXME: 横断お歩道になってしまう
		{
			desc:    "正常系: 小説の一文をテストしますわ",
			src:     "通勤時に通る横断歩道の話なんだ。よくあるだろ、事故が多い横断歩道の話だよ。",
			want:    "通勤時に通る横断お歩道の話なんですの。よくありますでしょう、お事故が多い横断お歩道の話ですわ。",
			opt:     nil,
			wantErr: false,
		},

		{
			desc:    "正常系: 小説の一文をテストしますわ",
			src:     "ただ内容は、偶然事故が多いって話じゃないから安心してくれ。実際に俺が体験した話だ。",
			want:    "ただお内容は、偶然お事故が多いって話ではありませんので安心してくださいまし。実際に私が体験した話ですわ。",
			opt:     nil,
			wantErr: false,
		},

		// FIXME: 長いですので、にしたい
		{
			desc:    "正常系: 小説の一文をテストしますわ",
			src:     "信じてもいいし、信じなくても良い。まぁ、ゆったり聞いてくれ。夜は長いからな。",
			want:    "信じてもいいですし、信じなくても良いですわ。まぁ、ゆったり聞いてくださいまし。夜は長いのでね。",
			opt:     nil,
			wantErr: false,
		},

		// FIXME: 戻る、の部分の処理が難しい
		// {
		// 	desc:    "正常系: 小説の一文をテストしますわ",
		// 	src:     "それで、横断歩道の話に戻るぞ。通勤時に使う横断歩道だからな、当然毎日通るだろ？",
		// 	want:    "それで、横断歩道の話に戻りますわね。通勤時に使う横断歩道ですのでね、当然毎日通るでしょう？",
		// 	opt:     nil,
		// 	wantErr: false,
		// },

		// FIXME: ですわけれどになってる
		{
			desc:    "正常系: 小説の一文をテストしますわ",
			src:     "入社直後は道も周囲の店も何もわからなくて新鮮に感じたもんだけどさ、一番楽に通勤できるルートが確立したら、すぐに通勤が退屈になるわけだ。",
			want:    "入社直後はお道もお周囲のお店も何もわからなくて新鮮に感じたものですわけれど、一番楽に通勤できるおルートが確立したら、すぐに通勤が退屈になるわけですわ。",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: お嬢様らしく短く笑いますわ",
			src:     "うふ",
			want:    "おほ",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: お嬢様らしく笑いますわ",
			src:     "うふふ",
			want:    "おほほ",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: お嬢様らしくそこそこ笑いますわ",
			src:     "うふふふふ",
			want:    "おほほほほ",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: お嬢様らしく長く笑いますわ",
			src:     "うふふふふふ",
			want:    "おほほほほほ",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: 切れ目があってもお嬢様らしく笑いますわ",
			src:     "うふうふふふふふふふふ",
			want:    "おほおほほほほほほほほ",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: パパはおパパ上、ママはおママ上とお呼びいたしますわ",
			src:     "パパ、ママ、父様、母様、パパ上、ママ上",
			want:    "おパパ上、おママ上、お父様、お母様、おパパ上、おママ上",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: 罵倒には「お」を付けませんのよ",
			src:     "カス",
			want:    "カス",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: こそあど言葉にも対応しておりましてよ",
			src:     "これ、この、ここ、こちら、こう、こんな。",
			want:    "こちら、こちらの、こちら、こちら、こう、このような。",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: こそあど言葉にも対応しておりましてよ",
			src:     "それ、その、そこ、そちら、そう、そんな。",
			want:    "そちら、そちらの、そちら、そちら、そう、そのような。",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: こそあど言葉にも対応しておりましてよ",
			src:     "あれは、あの、あそこ、あちら、ああ、あんな。",
			want:    "あちらは、あちらの、あちら、あちら、ああ、あのような。",
			opt:     nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got, err := Convert(tt.src, tt.opt)
			if tt.wantErr {
				assert.Error(err)
				assert.Empty(got)
				return
			}

			assert.NoError(err)
			assert.Equal(tt.want, got)
		})
	}
}
