package dto

type (
	AuthInput struct {
		UserAddress string `json:"user_address" db:"user_address"`
		Password    string `json:"password" db:"password"`
	}
	KeyImput struct {
		UserAddress   string `json:"user_address" db:"user_address"`
		Password      string `json:"password" db:"password"`
		EncryptedKey  []byte `json:"encrypted_key" db:"encrypted_key"`
		EncryptedData []byte `json:"encrypted_data" db:"encrypted_data"`
		KeyIV         []byte `json:"key_iv" db:"key_iv"`
		DataIV        []byte `json:"data_iv" db:"data_iv"`
	}
	KeyOutput struct {
		ID            string `json:"id" db:"id"`
		EncryptedKey  []byte `json:"encrypted_key" db:"encrypted_key"`
		EncryptedData []byte `json:"encrypted_data" db:"encrypted_data"`
		KeyIV         []byte `json:"key_iv" db:"key_iv"`
		DataIV        []byte `json:"data_iv" db:"data_iv"`
	}
	DeleteKeyInput struct {
		ID          string `json:"id" db:"id"`
		UserAddress string `json:"user_address" db:"user_address"`
		Password    string `json:"password" db:"password"`
	}
)
