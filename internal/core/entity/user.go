package entity

type User struct {
	ID       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
	Role     string `db:"role" json:"role"` // "student", "teacher", "parent", "guest"
}
