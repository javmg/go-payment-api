package repository

import (
	"github.com/javierjmgits/go-payment-api/payment/model"
	"github.com/jinzhu/gorm"
)

type PaymentRepository interface {
	GetAll() ([]model.Payment, error)
	GetByUid(uid string) (*model.Payment, error)
	Create(*model.Payment) (*model.Payment, error)
	Update(*model.Payment) (*model.Payment, error)
	Delete(*model.Payment) error
}

type paymentRepositoryImpl struct {
	db *gorm.DB
}

func NewPaymentRepositoryImpl(db *gorm.DB) PaymentRepository {
	return &paymentRepositoryImpl{
		db: db,
	}
}

func (pri *paymentRepositoryImpl) GetAll() ([]model.Payment, error) {

	var payments []model.Payment
	errorDB := pri.db.Find(&payments).Error

	if errorDB != nil {
		return nil, errorDB
	}

	return payments, nil
}

func (pri *paymentRepositoryImpl) GetByUid(uid string) (*model.Payment, error) {

	var payment model.Payment
	errorFind := pri.db.Where("uid = ?", uid).First(&payment).Error

	if errorFind != nil {
		return nil, errorFind
	}

	return &payment, nil
}

func (pri *paymentRepositoryImpl) Create(payment *model.Payment) (*model.Payment, error) {

	errorDB := pri.db.Create(&payment).Error

	if errorDB != nil {
		return nil, errorDB
	}

	return payment, nil

}

func (pri *paymentRepositoryImpl) Update(payment *model.Payment) (*model.Payment, error) {

	errorDB := pri.db.Save(&payment).Error

	if errorDB != nil {
		return nil, errorDB
	}

	return payment, nil
}

func (pri *paymentRepositoryImpl) Delete(payment *model.Payment) error {

	errorDB := pri.db.Delete(&payment).Error

	if errorDB != nil {
		return errorDB
	}

	return nil
}
