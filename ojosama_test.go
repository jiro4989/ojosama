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
