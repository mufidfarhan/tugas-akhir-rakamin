package repository

import (
	"backend-evermos/internal/pkg/entity"
	"context"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type AddressesRepository interface {
	Transactor

	GetAddressesByUserID(ctx context.Context, userID string) (res []entity.Address, err error)
	GetAddressByID(ctx context.Context, addressID string) (res entity.Address, err error)
	CreateAddress(ctx context.Context, userID string, data entity.Address) (res uint, err error)
	UpdateAddressByID(ctx context.Context, userID string, addressID string, data entity.Address) (err error)
	DeleteAddressByID(ctx context.Context, userID string, addressID string) (err error)

	VerifyAddressAvailability(ctx context.Context, addressID string) (err error)
	VerifyAddressOwner(ctx context.Context, addressID string, userID string) (err error)
}

type AddressesRepositoryImpl struct {
	transactor
}

func NewAddressRepository(db *gorm.DB) AddressesRepository {
	return &AddressesRepositoryImpl{
		transactor: transactor{
			db: db,
		},
	}
}

func (r *AddressesRepositoryImpl) GetAddressesByUserID(ctx context.Context, userID string) (res []entity.Address, err error) {
	if err := r.tx(ctx).Where("user_id = ?", userID).Find(&res).Error; err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return res, nil
}

func (r *AddressesRepositoryImpl) GetAddressByID(ctx context.Context, addressID string) (res entity.Address, err error) {
	if err := r.tx(ctx).First(&res, addressID).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (r *AddressesRepositoryImpl) CreateAddress(ctx context.Context, userID string, data entity.Address) (res uint, err error) {
	uid64, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID: %w", err)
	}
	data.UserID = uint(uid64)

	result := r.tx(ctx).Create(&data)
	if result.Error != nil {
		return res, result.Error
	}

	return data.ID, nil
}

func (r *AddressesRepositoryImpl) UpdateAddressByID(ctx context.Context, userID string, addressID string, data entity.Address) (err error) {
	var addressData entity.Address
	if err := r.tx(ctx).Where("id = ? AND user_id = ?", addressID, userID).First(&addressData).Error; err != nil {
		return gorm.ErrRecordNotFound
	}

	if err := r.tx(ctx).Model(addressData).Updates(&data).Where("id = ?", addressID).Error; err != nil {
		return err
	}

	return nil
}

func (r *AddressesRepositoryImpl) DeleteAddressByID(ctx context.Context, userID string, addressID string) (err error) {
	var addressData entity.Address
	if err := r.tx(ctx).Where("id = ? AND user_id = ?", addressID, userID).First(&addressData).Error; err != nil {
		return gorm.ErrRecordNotFound
	}

	if err := r.tx(ctx).Model(addressData).Delete(&addressData).Error; err != nil {
		return err
	}

	return nil
}

func (r *AddressesRepositoryImpl) VerifyAddressAvailability(ctx context.Context, addressID string) (err error) {
	var address entity.Address
	if err := r.tx(ctx).Where("id = ? ", addressID).First(&address).Error; err != nil {
		return err
	}

	return nil
}

func (r *AddressesRepositoryImpl) VerifyAddressOwner(ctx context.Context, addressID string, userID string) (err error) {
	var address entity.Address
	if err := r.tx(ctx).Where("id = ? AND user_id = ?", addressID, userID).First(&address).Error; err != nil {
		return err
	}

	return nil
}
