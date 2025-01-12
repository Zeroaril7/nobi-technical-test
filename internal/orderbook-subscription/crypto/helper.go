package crypto

func (m *CreateCrypto) ToCryptoEntity(e *CryptoEntity) *CryptoEntity {
	return &CryptoEntity{
		Pair: m.Pair,
	}
}
