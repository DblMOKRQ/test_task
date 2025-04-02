package entity

type User struct {
	Name       string `json:"name" required:"true"`
	Surname    string `json:"surname" required:"true"`
	Patronymic string `json:"patronymic" required:"false"`
}

type UserRequest struct {
	Name        string
	Surname     string
	Patronymic  string
	Age         int
	Gender      string
	Nationality string
}
type FullUser struct {
	ID          int
	Name        string
	Surname     string
	Patronymic  string
	Age         int
	Gender      string
	Nationality string
}
