package dto

type (
	User struct {
		ID          string `json:"id" db:"id"`
		UserAddress string `json:"user_address" db:"user_address"`
	}

	UserInput struct {
		UserAddress string `json:"user_address" db:"user_address"`
		Password    string `json:"password" db:"password"`
	}

	UpdatePasswordInput struct {
		UserAddress string `json:"user_address" db:"user_address"`
		Password    string `json:"password" db:"password"`
		NewPassword string `json:"new_password" db:"new_password"`
	}
)
