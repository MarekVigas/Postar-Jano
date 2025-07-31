package model

type Day struct {
	ID          int    `json:"id" db:"id"`
	Description string `json:"description" db:"description"`
	Capacity    int    `json:"capacity" db:"capacity"`
	LimitBoys   *int   `json:"limit_boys" db:"limit_boys"`
	LimitGirls  *int   `json:"limit_girls" db:"limit_girls"`
	Price       int    `json:"price" db:"price"`
	EventID     int    `json:"-" db:"event_id"`
}
