package model

type Owner struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Surname  string `db:"surname"`
	Email    string `db:"email"`
	Username string `db:"username"`
	Pass     string `db:"pass"`
	Phone    string `db:"phone"`
	Photo    string `db:"photo"`
	Gender   string `db:"gender"`
}
