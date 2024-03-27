package response

type User struct {
	ID        int64  `json:"userId"`
	Email     string `json:"email,omitempty"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
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
