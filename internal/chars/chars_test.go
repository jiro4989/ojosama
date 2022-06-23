package chars

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsExclQuesMark(t *testing.T) {
	tests := []struct {
		desc   string
		s      string
		wantOK bool
		wantEQ ExclQuesMark
	}{
		{
			desc:   "正常系: ！とはマッチいたしますわ",
			s:      "！",
			wantOK: true,
			wantEQ: newExcl("！", styleTypeFullWidth),
		},
		{
			desc:   "正常系: ❓とはマッチいたしますわ",
			s:      "❓",
			wantOK: true,
			wantEQ: newQues("❓", styleTypeEmoji),
		},
		{
			desc:   "正常系: 漆とはマッチいたしませんわ",
			s:      "漆",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got, got2 := IsExclQuesMark(tt.s)
			assert.Equal(tt.wantOK, got)
			if tt.wantOK {
				assert.Equal(&tt.wantEQ, got2)
			}
		})
	}
}

func TestSampleExclQuesByValue(t *testing.T) {
	tests := []struct {
		desc    string
		v       string
		t       *TestMode
		want    ExclQuesMark
		wantNil bool
	}{
		{
			desc: "正常系: ！とはマッチいたしますわ",
			v:    "！",
			t:    &TestMode{Pos: 0},
			want: newExcl("！", styleTypeFullWidth),
		},
		{
			desc: "正常系: ❓とはマッチいたしますわ",
			v:    "❓",
			t:    &TestMode{Pos: 2},
			want: newQues("❓", styleTypeEmoji),
		},
		{
			desc:    "正常系: 菫とはマッチいたしませんわ",
			v:       "菫",
			t:       &TestMode{Pos: 2},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := SampleExclQuesByValue(tt.v, tt.t)
			if tt.wantNil {
				assert.Nil(got)
				return
			}
			assert.Equal(&tt.want, got)
		})
	}
}

func TestFindExclQuesByStyleAndMeaning(t *testing.T) {
	tests := []struct {
		desc    string
		s       StyleType
		m       MeaningType
		want    ExclQuesMark
		wantNil bool
	}{
		{
			desc: "正常系: ❗を指定いたしますわ",
			s:    styleTypeEmoji,
			m:    meaningTypeExcl,
			want: newExcl("❗", styleTypeEmoji),
		},
		{
			desc: "正常系: ？を指定いたしますわ",
			s:    styleTypeFullWidth,
			m:    meaningTypeQues,
			want: newQues("？", styleTypeFullWidth),
		},
		{
			desc:    "正常系: 不明な要素の場合は何もお返しいたしませんわ",
			s:       styleTypeUnknown,
			m:       meaningTypeExcl,
			wantNil: true,
		},
		{
			desc:    "正常系: 不明な要素の場合は何もお返しいたしませんわ",
			s:       styleTypeFullWidth,
			m:       meaningTypeUnknown,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := FindExclQuesByStyleAndMeaning(tt.s, tt.m)
			if tt.wantNil {
				assert.Nil(got)
				return
			}
			assert.Equal(&tt.want, got)
		})
	}
}
