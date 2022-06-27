package ojosama

import "github.com/ikawaha/kagome/v2/tokenizer"

type TokenCtl struct {
	tokens []tokenizer.Token
	pos    int
}

func NewTokenCtl(tokens []tokenizer.Token) *TokenCtl {
	return &TokenCtl{
		tokens: tokens,
	}
}

func (t *TokenCtl) Runnable() bool {
	return t.pos < len(t.tokens)
}

func (t *TokenCtl) Next() {
	t.pos++
}

func (t *TokenCtl) TokenData() *tokenizer.TokenData {
	data := tokenizer.NewTokenData(t.tokens[t.pos])
	return &data
}

func (t *TokenCtl) PrevTokenData() *tokenizer.TokenData {
	if t.availablePrevToken() {
		data := tokenizer.NewTokenData(t.tokens[t.pos-1])
		return &data
	}
	return nil
}

func (t *TokenCtl) NextTokenData() *tokenizer.TokenData {
	if t.availableNextToken() {
		data := tokenizer.NewTokenData(t.tokens[t.pos+1])
		return &data
	}
	return nil
}

func (t *TokenCtl) availablePrevToken() bool {
	return 0 < t.pos
}

func (t *TokenCtl) availableNextToken() bool {
	return len(t.tokens) <= t.pos+1
}
