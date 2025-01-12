package crypto

type CryptoEntity struct {
	ID   int64  `json:"id" gorm:"primaryKey"`
	Pair string `json:"pair"`
}

func (CryptoEntity) TableName() string {
	return "crypto"
}
