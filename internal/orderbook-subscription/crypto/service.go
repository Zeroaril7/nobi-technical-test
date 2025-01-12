package crypto

import "context"

type CryptoService interface {
	Add(data *CryptoEntity, ctx context.Context) error
	FindAll(ctx context.Context) ([]CryptoEntity, int64, error)
	Delete(id string, ctx context.Context) error
}

type cryptoServiceImpl struct {
	repository CryptoRepository
}

func (s *cryptoServiceImpl) Add(data *CryptoEntity, ctx context.Context) error {
	return s.repository.Add(data, ctx)
}

func (s *cryptoServiceImpl) Delete(id string, ctx context.Context) error {
	return s.repository.Delete(id, ctx)
}

func (s *cryptoServiceImpl) FindAll(ctx context.Context) ([]CryptoEntity, int64, error) {
	return s.repository.FindAll(ctx)
}

func NewService(repository CryptoRepository) CryptoService {
	return &cryptoServiceImpl{repository: repository}
}
