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
			desc:    "正常系: 草はハーブに変換してお嬢様に変換するんですの",
			src:     "草です",
			want:    "おハーブですわ",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: 変換対象以外のTokenはそのまま残すのですわ",
			src:     "これは草です",
			want:    "これはおハーブですわ",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: 名詞の手前には「お」をお付けいたしますわ",
			src:     "バイオ",
			want:    "おバイオ",
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
			want:    "これはgrassですわ。あれはabcdefg12345ですわ",
			opt:     nil,
			wantErr: false,
		},
		{
			desc:    "正常系: 「ます」は「ますわ」に変換しますわ",
			src:     "使用します",
			want:    "使用しますわ",
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
