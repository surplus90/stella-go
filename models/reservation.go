package models

type Reservation struct {
	Idx           int64
	UserName      string
	DeckIdx       int
	AmountCards   int
	SelectedCards int
	WayToArray    int
	ReservationAt int64
	SetcardsAt    int64
	CreatedAt     int64
	EncKey        string
}