package service

import gt "github.com/bas24/googletranslatefree"

type Translate struct {
	SourceLan string
	TargetLan string
}

func NewTransletor(sourceLan, targetLan string) Translate {
	return Translate{
		SourceLan: sourceLan,
		TargetLan: targetLan,
	}
}

func (t *Translate) TransletorMSG(text string) (string, error) {
	result, err := gt.Translate(text, t.SourceLan, t.TargetLan)

	return result, err
}
