package response

type User struct {
	ID    string `json:"userId"`
	Email string `json:"email,omitempty"`
	Name  string `json:"name"`
}

type Register struct {
	Email       string `json:"email,omitempty"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

type Login struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}
