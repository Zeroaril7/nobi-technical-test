package crypto

type CreateCrypto struct {
	Pair string `json:"pair" validate:"required"`
}
