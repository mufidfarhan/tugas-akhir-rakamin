package repository

import (
	"backend-evermos/internal/pkg/entity"
	"context"
	"errors"

	"gorm.io/gorm"
)

type UsersRepository interface {
	Transactor

	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (res entity.User, err error)
	CreateUser(ctx context.Context, data entity.User) (res uint, err error)
	GetUserByID(ctx context.Context, userID string) (res entity.User, err error)
	UpdateUserByID(ctx context.Context, userID string, data entity.User) (err error)

	VerifyEmail(ctx context.Context, email string) (err error)
	VerifyPhoneNumber(ctx context.Context, phoneNumber string) (err error)
}

type UsersRepositoryImpl struct {
	transactor
}

func NewUsersRepository(db *gorm.DB) UsersRepository {
	return &UsersRepositoryImpl{
		transactor: transactor{
			db: db,
		},
	}
}

func (r *UsersRepositoryImpl) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (res entity.User, err error) {
	if err := r.tx(ctx).Where("phone_number = ?", phoneNumber).First(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (r *UsersRepositoryImpl) CreateUser(ctx context.Context, data entity.User) (res uint, err error) {
	result := r.tx(ctx).Create(&data)
	if result.Error != nil {
		return res, result.Error
	}

	return data.ID, nil
}

func (r *UsersRepositoryImpl) GetUserByID(ctx context.Context, userID string) (res entity.User, err error) {
	if err := r.tx(ctx).Where("id = ?", userID).First(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (r *UsersRepositoryImpl) UpdateUserByID(ctx context.Context, userID string, data entity.User) (err error) {
	var userData entity.User
	if err := r.tx(ctx).Where("id = ?", userID).First(&userData).Error; err != nil {
		return gorm.ErrRecordNotFound
	}

	if err := r.tx(ctx).Model(userData).Updates(&data).Where("id = ?", userID).Error; err != nil {
		return err
	}

	return nil
}

func (r *UsersRepositoryImpl) VerifyEmail(ctx context.Context, email string) (err error) {
	var count int64
	if err := r.tx(ctx).Model(&entity.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("email sudah digunakan")
	}

	return nil
}

func (r *UsersRepositoryImpl) VerifyPhoneNumber(ctx context.Context, phoneNumber string) (err error) {
	var count int64
	if err := r.tx(ctx).Model(&entity.User{}).Where("phone_number = ?", phoneNumber).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("nomor hp sudah digunakan")
	}

	return nil
}
