package ojosama

import (
	"testing"

	"github.com/jiro4989/ojosama/internal/chars"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	opt := &ConvertOption{
		DisableKutenToExclamation: true,
	}

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
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 「名詞,一般」の後に「動詞,自立」が来た時は2Tokenで1つの動詞として解釈いたしますので、「お」を付けませんわ",
			src:     "プレイする",
			want:    "プレイいたしますわ",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 「プレイする」の例文ですわ",
			src:     "〇〇をプレイする",
			want:    "〇〇をプレイいたしますわ",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 二人称はすべて「貴方」に変換いたしますわ",
			src:     "あなたは。あんたは。おまえは。お前は。てめぇは。てめえは。貴様は。君は。",
			want:    "貴方は。貴方は。貴方は。貴方は。貴方は。貴方は。貴方は。貴方は。",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: わたしはわたくしに変換いたしますわ",
			src:     "わたし",
			want:    "わたくし",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 「ですか」のときは「ですの」に変換いたしますわ",
			src:     "ビデオテープはどこで使うんですか",
			want:    "おビデオテープはどちらで使うんですの",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 「汚い|きたない」のときは「きったねぇ」に変換いたしますわ",
			src:     "汚いです",
			want:    "きったねぇですわ",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 「汚い|きたない」のときは「きったねぇ」に変換いたしますわ",
			src:     "きたないです",
			want:    "きったねぇですわ",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 「臭い|くさい」のときは「くっせぇ」に変換いたしますわ",
			src:     "臭いです",
			want:    "くっせぇですわ",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 「臭い|くさい」のときは「くっせぇ」に変換いたしますわ",
			src:     "くさいです",
			want:    "くっせぇですわ",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 複数の文も処理できますわ",
			src:     "〇〇をプレイする。きたないわ。",
			want:    "〇〇をプレイいたしますわ。きったねぇですわ。",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 「ました」は「おりました」ですの",
			src:     "わたしも使ってました",
			want:    "わたくしも使っておりましたわ",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: アルファベット単語の場合は「お」をつけませんの",
			src:     "これはgrassです。あれはabcdefg12345です",
			want:    "こちらはgrassですわ。あちらはabcdefg12345ですわ",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 「します」は「いたしますわ」に変換しますわ",
			src:     "使用します",
			want:    "使用いたしますわ",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 名詞が連続する場合は最初の1つ目にだけ「お」を付けますわ",
			src:     "一般女性。経年劣化。トップシークレット",
			want:    "お一般女性。お経年劣化。おトップシークレット",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: お嬢様は変更しませんわ",
			src:     "お嬢様",
			want:    "お嬢様",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: サロメお嬢様は固有名詞ですわ",
			src:     "わたしは壱百満天原サロメです",
			want:    "わたくしは壱百満天原サロメですわ",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 名字も1単語として認識いたしますわ",
			src:     "わたしの名字は壱百満天原です",
			want:    "わたくしのお名字は壱百満天原ですわ",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 壱百満点の笑顔を皆様方にお届けするのがわたくしの使命ですわ",
			src:     "壱百満点",
			want:    "壱百満点",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 部分一致の場合は処理しませんわ",
			src:     "壱",
			want:    "お壱",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 文の途中の「する」の場合は変換しませんわ",
			src:     "テキストファイルをまるごと変換する場合は、以下のように実行します。",
			want:    "おテキストファイルをまるごと変換する場合は、以下のように実行いたしますわ。",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 文の途中の「する」の場合は変換しませんわ",
			src:     "以下のお嬢様を可能な限り再現するアルゴリズムを目指してます。",
			want:    "以下のお嬢様を可能な限り再現するおアルゴリズムを目指してますわ。",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: すでに直前に「お」が付与されている場合は、「お」を付与しませんわ",
			src:     "お寿司。おハーブ",
			want:    "お寿司。おハーブ",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 「ください」は「くださいまし」に変換いたしますわ",
			src:     "お使いください",
			want:    "お使いくださいまし",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 「ください」は「くださいまし」に変換いたしますわ",
			src:     "お使いください！",
			want:    "お使いくださいまし！",
			opt:     opt,
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
				forceCharsTestMode: &chars.TestMode{
					Pos: 0,
				},
			},
			wantErr: false,
		},
		{
			desc: "正常系: 半角感嘆符にも反応しますわ！",
			src:  "これはハーブです!これもハーブです?",
			want: "こちらはおハーブですわ～～!!!こちらもおハーブですわ～～???",
			opt: &ConvertOption{
				forceAppendLongNote: forceAppendLongNote{
					enable:               true,
					wavyLineCount:        2,
					exclamationMarkCount: 3,
				},
				forceCharsTestMode: &chars.TestMode{
					Pos: 1,
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
				forceCharsTestMode: &chars.TestMode{
					Pos: 0,
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
				forceCharsTestMode: &chars.TestMode{
					Pos: 0,
				},
			},
			wantErr: false,
		},

		{
			desc:    "正常系: （動詞）ないはそのままですわ",
			src:     "限らない、飾らない、数えない",
			want:    "限らない、飾らない、数えない",
			opt:     opt,
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
				forceCharsTestMode: &chars.TestMode{
					Pos: 0,
				},
			},
			wantErr: false,
		},

		{
			desc:    "正常系: 一人称はすべて「私」に変換いたしますわ",
			src:     "俺は。オレは。おれは。僕は。ボクは。ぼくは。あたしは。わたしは。",
			want:    "私は。ワタクシは。わたくしは。私は。ワタクシは。わたくしは。わたくしは。わたくしは。",
			opt:     opt,
			wantErr: false,
		},

		// FIXME: 話しますわね、にしたいけれど「話す」で1単語と判定されているの
		// で変換が難しい
		{
			desc:    "正常系: 小説の一文をテストしますわ",
			src:     "俺の体験した怖い話、聞いてくれるか？ありがとう。じゃぁ、話すよ。",
			want:    "私の体験した怖い話、聞いてくれますの？ありがとうございますわ。それでは、話すよ。",
			opt:     opt,
			wantErr: false,
		},

		{
			desc:    "正常系: 小説の一文をテストしますわ",
			src:     "俺の転職前の会社の頃の話だから、もう3年も前の話になる。",
			want:    "私の転職前のお会社の頃の話ですので、もう3年も前の話になりますわ。",
			opt:     opt,
			wantErr: false,
		},

		{
			desc:    "正常系: 小説の一文をテストしますわ",
			src:     "通勤時に通る横断歩道の話なんだ。よくあるだろ、事故が多い横断歩道の話だよ。",
			want:    "通勤時に通る横断歩道の話なんですの。よくありますでしょう、お事故が多い横断歩道の話ですわ。",
			opt:     opt,
			wantErr: false,
		},

		{
			desc:    "正常系: 小説の一文をテストしますわ",
			src:     "ただ内容は、偶然事故が多いって話じゃないから安心してくれ。実際に俺が体験した話だ。",
			want:    "ただお内容は、偶然お事故が多いって話ではありませんので安心してくださいまし。実際に私が体験した話ですわ。",
			opt:     opt,
			wantErr: false,
		},

		// FIXME: 長いですので、にしたい
		{
			desc:    "正常系: 小説の一文をテストしますわ",
			src:     "信じてもいいし、信じなくても良い。まぁ、ゆったり聞いてくれ。夜は長いからな。",
			want:    "信じてもいいですし、信じなくても良いですわ。まぁ、ゆったり聞いてくださいまし。夜は長いのでね。",
			opt:     opt,
			wantErr: false,
		},

		// FIXME: 戻る、の部分の処理が難しい
		// {
		// 	desc:    "正常系: 小説の一文をテストしますわ",
		// 	src:     "それで、横断歩道の話に戻るぞ。通勤時に使う横断歩道だからな、当然毎日通るだろ？",
		// 	want:    "それで、横断歩道の話に戻りますわね。通勤時に使う横断歩道ですのでね、当然毎日通るでしょう？",
		// 	opt:     opt,
		// 	wantErr: false,
		// },

		// FIXME: ですわけれどになってる
		{
			desc:    "正常系: 小説の一文をテストしますわ",
			src:     "入社直後は道も周囲の店も何もわからなくて新鮮に感じたもんだけどさ、一番楽に通勤できるルートが確立したら、すぐに通勤が退屈になるわけだ。",
			want:    "入社直後はお道もお周囲のお店も何もわからなくて新鮮に感じたものですわけれど、一番楽に通勤できるおルートが確立したら、すぐに通勤が退屈になるわけですわ。",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: お嬢様らしく短く笑いますわ",
			src:     "うふ",
			want:    "おほ",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: お嬢様らしく笑いますわ",
			src:     "うふふ",
			want:    "おほほ",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: お嬢様らしくそこそこ笑いますわ",
			src:     "うふふふふ",
			want:    "おほほほほ",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: お嬢様らしく長く笑いますわ",
			src:     "うふふふふふ",
			want:    "おほほほほほ",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 切れ目があってもお嬢様らしく笑いますわ",
			src:     "うふうふふふふふふふふ",
			want:    "おほおほほほほほほほほ",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: パパはおパパ上、ママはおママ上とお呼びいたしますわ",
			src:     "パパ、ママ、父様、母様、パパ上、ママ上",
			want:    "おパパ上、おママ上、お父様、お母様、おパパ上、おママ上",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: パパはおパパ上、ママはおママ上とお呼びいたしますわ",
			src:     "皆、皆様",
			want:    "皆様方、皆様方",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 罵倒には「お」を付けませんのよ",
			src:     "カス",
			want:    "カス",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: こそあど言葉（こ）にも対応しておりましてよ",
			src:     "これ、この、ここ、こちら、こう、こんな。",
			want:    "こちら、こちらの、こちら、こちら、こう、このような。",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: こそあど言葉（そ）にも対応しておりましてよ",
			src:     "それ、その、そこ、そちら、そう、そんな。",
			want:    "そちら、そちらの、そちら、そちら、そう、そのような。",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: こそあど言葉（あ）にも対応しておりましてよ",
			src:     "あれは、あの、あそこ、あちら、ああ、あんな。",
			want:    "あちらは、あちらの、あちら、あちら、ああ、あのような。",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: こそあど言葉（ど）にも対応しておりましてよ",
			src:     "どれ、どの、どこ、どちら、どう、どんな。",
			want:    "どちら、どちらの、どちら、どちら、どう、どのような。",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 二重丁寧語にはしませんの",
			src:     "お嬢様、おにぎり、お腹、お寿司、おませさん、お利口さん",
			want:    "お嬢様、おにぎり、お腹、お寿司、おませさん、お利口さん",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 読みが「お」で始まる言葉には「お」を付けませんの",
			src:     "追い剥ぎ、大型、大分、おろし金",
			want:    "追い剥ぎ、大型、大分、おろし金",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: サ辺接続が手前に来る単語には「お」を付けませんわ",
			src:     "横断歩道、解体新書、回覧板、回転寿司、駐車場",
			want:    "横断歩道、解体新書、回覧板、回転寿司、駐車場",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 名詞＋動詞＋終助詞の組み合わせも変換いたしますわ～～！！この処理すっごく大変でしたの！！",
			src:     "野球しようぜ。サッカーやろうよ。バスケやるか。柔道やるな。陸上すんな。テニスするぞ。卓球やるべ。ゲームするの。",
			want:    "お野球をいたしませんこと。おサッカーをいたしませんこと。おバスケをいたしますわ。お柔道をしてはいけませんわ。お陸上をしてはいけませんわ。おテニスをいたしますわよ。お卓球をいたしませんこと。おゲームをいたしますわよ。",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 名詞＋した（過去形）＋終助詞も変換いたしますわ！",
			src:     "野球したぜ。サッカーやったよ。バスケしたで。柔道やったわ。陸上したぞ。",
			want:    "お野球をいたしましたわ。おサッカーをいたしましたわ。おバスケをいたしましたわ。お柔道をいたしましたわ。お陸上をいたしましたわ。",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 単純に「名詞＋した」だけでは変換いたしませんわ。まだ文の終わりではありませんので",
			src:     "野球した後。",
			want:    "野球した後。",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 名詞＋動詞＋助動詞のみで終助詞がない場合でもエラーにはなりませんのよ",
			src:     "流鏑馬やろう",
			want:    "流鏑馬やろう",
			opt:     opt,
			wantErr: false,
		},
		{
			desc: "正常系: すべて全角文字に変換しますわ",
			src:  "です！？!?❗❓",
			want: "ですわ～！？！？！？",
			opt: &ConvertOption{
				forceAppendLongNote: forceAppendLongNote{
					enable:               true,
					wavyLineCount:        1,
					exclamationMarkCount: 1,
				},
				forceCharsTestMode: &chars.TestMode{
					Pos: 0,
				},
			},
			wantErr: false,
		},
		{
			desc: "正常系: すべて絵文字に変換しますわ",
			src:  "です！？!?❗❓",
			want: "ですわ～❗❓❗❓❗❓",
			opt: &ConvertOption{
				forceAppendLongNote: forceAppendLongNote{
					enable:               true,
					wavyLineCount:        1,
					exclamationMarkCount: 1,
				},
				forceCharsTestMode: &chars.TestMode{
					Pos: 2,
				},
			},
			wantErr: false,
		},
		{
			desc: "正常系: 絵文字を連続して付与もできますわ",
			src:  "です！寿司",
			want: "ですわ～～❗❗❗お寿司",
			opt: &ConvertOption{
				forceAppendLongNote: forceAppendLongNote{
					enable:               true,
					wavyLineCount:        2,
					exclamationMarkCount: 3,
				},
				forceCharsTestMode: &chars.TestMode{
					Pos: 2,
				},
			},
			wantErr: false,
		},
		{
			desc:    "正常系: 意味のない文章のテストですわ",
			src:     "あ！い❓❗う",
			want:    "あ！い❓❗う",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 「ー」の前に「お」は付きませんの",
			src:     "しようぜー。しようぜーー。しようぜ～。しようぜ～～。しようぜー～。",
			want:    "しようぜー。しようぜーー。しようぜ～。しようぜ～～。しようぜー～。",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: なん(じゃ|や|だ)は、いずれも「なんですの」に変換いたしますわ",
			src:     "なんじゃこれ。なんだこれ。なんやこれ。",
			want:    "なんですのこちら。なんですのこちら。なんですのこちら。",
			opt:     opt,
			wantErr: false,
		},
		{
			desc:    "正常系: 名詞＋(じゃ|だ|ですの)ですの",
			src:     "アホだ。カラスや。バナナじゃ。これだ。それや。あれじゃ。",
			want:    "おアホですの。おカラスですの。おバナナですの。これですの。それですの。あれですの。",
			opt:     opt,
			wantErr: false,
		},
		{
			// FIXME: 「幽霊が怖い」で終わるべきところを
			// 「幽霊が怖」で終わった場合も変換されてしまう
			desc: "正常系: 形容詞＋自立の後に「ですわ」を付与する場合、「。」をランダムに！に変換いたしますわ。変換記号の@1を誤って変換したりはしませんわ",
			src:  "@1悲しいことはとても悲しい。幽霊が怖い。幽霊が怖。",
			want: "@1悲しいことはとても悲しいですわ❗お幽霊が怖いですわ❗お幽霊が怖ですわ❗",
			opt: &ConvertOption{
				forceKutenToExclamation: true,
			},
			wantErr: false,
		},
		{
			desc: "正常系: 「。」で終わる時に「！」に変換いたしますわ！",
			src:  "プレイする。ショットガンだ。",
			want: "プレイいたしますわ❗おショットガンですの❗",
			opt: &ConvertOption{
				forceKutenToExclamation: true,
			},
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
