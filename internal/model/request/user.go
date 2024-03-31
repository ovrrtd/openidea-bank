package request

type Register struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}
