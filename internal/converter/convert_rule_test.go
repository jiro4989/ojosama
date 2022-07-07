package converter

import (
	"testing"

	"github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/stretchr/testify/assert"
)

func TestGetMeaningType(t *testing.T) {
	tests := []struct {
		desc    string
		typeMap map[MeaningType]ConvertConditions
		data    tokenizer.TokenData
		wantMT  MeaningType
		wantOK  bool
	}{
		{
			desc: "正常系: 一致したmeaningTypeを返却いたしますわ",
			typeMap: map[MeaningType]ConvertConditions{
				meaningTypeCoercion: {
					{
						Features: []string{"名詞"},
						Reading:  "a",
					},
					{
						Features: []string{"名詞"},
						Reading:  "b",
					},
				},
				meaningTypeHope: {
					{
						Features: []string{"名詞"},
						Reading:  "c",
					},
				},
			},
			data: tokenizer.TokenData{
				Features: []string{"名詞"},
				Reading:  "b",
			},
			wantMT: meaningTypeCoercion,
			wantOK: true,
		},
		{
			desc: "正常系: いずれとも一致しない場合はunknownですわ",
			typeMap: map[MeaningType]ConvertConditions{
				meaningTypeCoercion: {
					{
						Features: []string{"名詞"},
						Reading:  "z",
					},
					{
						Features: []string{"名詞"},
						Surface:  "z",
					},
				},
			},
			data: tokenizer.TokenData{
				Features: []string{"名詞"},
				Reading:  "b",
			},
			wantMT: meaningTypeUnknown,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			gotMT, gotOK := GetMeaningType(tt.typeMap, tt.data)
			assert.Equal(tt.wantMT, gotMT)
			assert.Equal(tt.wantOK, gotOK)
		})
	}
}
