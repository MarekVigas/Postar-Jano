package model

type Stat struct {
	EventID    int  `json:"event_id" db:"event_id"`
	DayID      int  `json:"day_id" db:"day_id"`
	Capacity   int  `json:"capacity" db:"capacity"`
	LimitBoys  *int `json:"limit_boys" db:"limit_boys"`
	LimitGirls *int `json:"limit_girls" db:"limit_girls"`
	BoysCount  int  `json:"boys_count" db:"boys_count"`
	GirlsCount int  `json:"girls_count" db:"girls_count"`
}
