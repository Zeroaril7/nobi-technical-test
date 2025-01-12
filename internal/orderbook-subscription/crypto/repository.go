package crypto

import (
	"context"

	"gorm.io/gorm"
)

type CryptoRepository interface {
	Add(data *CryptoEntity, ctx context.Context) error
	FindAll(ctx context.Context) ([]CryptoEntity, int64, error)
	Delete(id string, ctx context.Context) error
}

type cryptoRepositoryImpl struct {
	db *gorm.DB
}

func (r *cryptoRepositoryImpl) Add(data *CryptoEntity, ctx context.Context) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *cryptoRepositoryImpl) Delete(id string, ctx context.Context) error {
	return r.db.WithContext(ctx).Delete(&CryptoEntity{}, "id = ?", id).Error
}

func (r *cryptoRepositoryImpl) FindAll(ctx context.Context) (result []CryptoEntity, total int64, err error) {

	db := r.db.WithContext(ctx)

	if err = db.Model(&CryptoEntity{}).Count(&total).Error; err != nil {
		return
	}

	if err = db.Find(&result).Error; err != nil {
		return
	}

	return
}

func NewRepository(db *gorm.DB) CryptoRepository {
	return &cryptoRepositoryImpl{db: db}
}
