package models

type PickCards struct {
	Idx            int64
	ReservationIdx int64
	EncKey         string
	Cards          string
    CreatedAt      int64
}