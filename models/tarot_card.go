package models

type TarotCard struct {
	Idx        int64
	DeckIdx    int64
	Seq        int
	CardName   string
	IsMajor    bool
	ImgPath    string
	CreatedAt  int64
}