package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/javierjmgits/go-payment-api/payment/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

//
// mock data

var now = time.Now()

//
// mocks

type paymentRepositoryImplMock struct {
	mock.Mock
}

func (mock *paymentRepositoryImplMock) GetAll() ([]model.Payment, error) {
	args := mock.Mock.Called(nil)

	results := args.Get(0)

	if results != nil {
		return results.([]model.Payment), nil
	}

	return nil, args.Get(1).(error)
}

func (mock *paymentRepositoryImplMock) GetByUid(uid string) (*model.Payment, error) {

	args := mock.Mock.Called(uid)

	result := args.Get(0)

	if result != nil {
		return result.(*model.Payment), nil
	}

	return nil, args.Get(1).(error)
}

func (mock *paymentRepositoryImplMock) Create(payment *model.Payment) (*model.Payment, error) {

	args := mock.Mock.Called(payment)

	result := args.Get(0)

	if result != nil {
		return result.(*model.Payment), nil
	}

	return nil, args.Get(1).(error)
}

func (mock *paymentRepositoryImplMock) Update(payment *model.Payment) (*model.Payment, error) {

	args := mock.Mock.Called(payment)

	result := args.Get(0)

	if result != nil {
		return result.(*model.Payment), nil
	}

	return nil, args.Get(1).(error)

}

func (mock *paymentRepositoryImplMock) Delete(payment *model.Payment) error {

	mock.Mock.Called(payment)

	return nil
}

//
// tests

func TestGetPaymentsKoError(t *testing.T) {

	router, mockRepository := setUp()
	mockRepository.On("GetAll", nil).Return(nil, errors.New("DB error"))

	req := httptest.NewRequest("GET", "http://localhost:8080/api/v1/payments", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	mockRepository.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestGetPayments(t *testing.T) {

	router, mockRepository := setUp()
	expectedPayment1 := expectedPayment("myUid", false)
	mockRepository.On("GetAll", nil).Return([]model.Payment{*expectedPayment1}, nil)

	req := httptest.NewRequest("GET", "http://localhost:8080/api/v1/payments", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	mockRepository.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := ioutil.ReadAll(resp.Body)

	var payments []PaymentView

	json.Unmarshal(body, &payments)

	// verify

	assert.Len(t, payments, 1)
}

func TestGetPaymentByUidKoNotFound(t *testing.T) {

	router, mockRepository := setUp()
	mockRepository.On("GetByUid", "unknown").Return(nil, errors.New("record not found"))

	req := httptest.NewRequest("GET", "http://localhost:8080/api/v1/payments/uid/unknown", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	mockRepository.AssertExpectations(t)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGetPaymentByUidKoError(t *testing.T) {

	router, mockRepository := setUp()
	mockRepository.On("GetByUid", "myUid").Return(nil, errors.New("DB error"))

	req := httptest.NewRequest("GET", "http://localhost:8080/api/v1/payments/uid/myUid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	mockRepository.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestGetPaymentByUid(t *testing.T) {

	router, mockRepository := setUp()
	expectedPayment := expectedPayment("myUid", false)
	mockRepository.On("GetByUid", "myUid").Return(expectedPayment, nil)

	req := httptest.NewRequest("GET", "http://localhost:8080/api/v1/payments/uid/myUid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	mockRepository.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := ioutil.ReadAll(resp.Body)

	var payment PaymentView

	json.Unmarshal(body, &payment)

	// verify

	assert.Equal(t, expectedPayment.Uid, payment.Uid)
	assert.Equal(t, expectedPayment.AccountOrigin, payment.AccountOrigin)
	assert.Equal(t, expectedPayment.AccountTarget, payment.AccountTarget)
	assert.Equal(t, expectedPayment.Amount, payment.Amount)
	assert.Equal(t, expectedPayment.Date, payment.Date)
	assert.Equal(t, expectedPayment.Processed, payment.Processed)
	assert.Equal(t, expectedPayment.ProcessedDate, payment.ProcessedDate)
}

func TestCreatePaymentKoMissingMandatoryFields(t *testing.T) {

	router, mockRepository := setUp()
	var paymentCreate PaymentCreate
	paymentCreateAsBytes, _ := json.Marshal(paymentCreate)

	req := httptest.NewRequest("POST", "http://localhost:8080/api/v1/payments", bytes.NewReader(paymentCreateAsBytes))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	mockRepository.AssertExpectations(t)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreatePayment(t *testing.T) {

	router, mockRepository := setUp()
	expectedPayment := expectedPayment("myUid", false)
	paymentCreate := PaymentCreate{
		AccountOrigin: expectedPayment.AccountOrigin,
		AccountTarget: expectedPayment.AccountTarget,
		Amount:        expectedPayment.Amount,
		Date:          expectedPayment.Date,
	}

	mockRepository.On("Create", mock.MatchedBy(func(passed *model.Payment) bool {

		// uid is generated
		if passed.Uid == "" || passed.Uid == expectedPayment.Uid {
			return false
		}

		if passed.AccountOrigin != expectedPayment.AccountOrigin {
			return false
		}

		if passed.AccountTarget != expectedPayment.AccountTarget {
			return false
		}

		if passed.Date != expectedPayment.Date {
			return false
		}

		return true
	})).Return(expectedPayment, nil)

	mockRepository.On("Create", mock.Anything).Return(expectedPayment, nil)

	paymentCreateAsBytes, _ := json.Marshal(paymentCreate)

	req := httptest.NewRequest("POST", "http://localhost:8080/api/v1/payments", bytes.NewReader(paymentCreateAsBytes))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	mockRepository.AssertExpectations(t)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	body, _ := ioutil.ReadAll(resp.Body)

	var payment PaymentView

	json.Unmarshal(body, &payment)

	// verify

	assert.NotEmpty(t, payment.Uid)
	assert.Equal(t, paymentCreate.AccountOrigin, payment.AccountOrigin)
	assert.Equal(t, paymentCreate.AccountTarget, payment.AccountTarget)
	assert.Equal(t, paymentCreate.Amount, payment.Amount)
	assert.Equal(t, paymentCreate.Date, payment.Date)
	assert.False(t, payment.Processed)
	assert.Empty(t, payment.ProcessedDate)
}

func TestFlagPaymentAsProcessedByUidKoNotFound(t *testing.T) {

	router, mockRepository := setUp()
	mockRepository.On("GetByUid", "unknown").Return(nil, errors.New("not found"))

	req := httptest.NewRequest("PATCH", "http://localhost:8080/api/v1/payments/uid/unknown/processed", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	mockRepository.AssertExpectations(t)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestFlagPaymentAsProcessedByUidKoAlreadyProcessed(t *testing.T) {

	router, mockRepository := setUp()
	expectedPayment := expectedPayment("myUid", true)

	mockRepository.On("GetByUid", "myUid").Return(expectedPayment, nil)

	req := httptest.NewRequest("PATCH", "http://localhost:8080/api/v1/payments/uid/myUid/processed", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	mockRepository.AssertExpectations(t)
	assert.Equal(t, http.StatusConflict, resp.StatusCode)
}

func TestFlagPaymentAsProcessedByUid(t *testing.T) {

	router, mockRepository := setUp()

	expectedPayment := expectedPayment("myUid", false)

	mockRepository.On("GetByUid", "myUid").Return(expectedPayment, nil)
	mockRepository.On("Update", mock.MatchedBy(func(passed *model.Payment) bool {

		// uid is generated
		if passed.Processed == false {
			return false
		}

		if passed.ProcessedDate == nil {
			return false
		}

		return true
	})).Return(expectedPayment, nil)

	req := httptest.NewRequest("PATCH", "http://localhost:8080/api/v1/payments/uid/myUid/processed", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	mockRepository.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := ioutil.ReadAll(resp.Body)

	var payment PaymentView

	json.Unmarshal(body, &payment)

	// verify

	assert.True(t, payment.Processed)
	assert.Equal(t, expectedPayment.ProcessedDate, payment.ProcessedDate)
}

func TestDeletePaymentByUidKoNotFound(t *testing.T) {

	router, mockRepository := setUp()
	mockRepository.On("GetByUid", "unknown").Return(nil, errors.New("not found"))

	req := httptest.NewRequest("DELETE", "http://localhost:8080/api/v1/payments/uid/unknown", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	mockRepository.AssertExpectations(t)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestDeletePaymentByUidKoAlreadyProcessed(t *testing.T) {

	router, mockRepository := setUp()
	expectedPayment := expectedPayment("myUid", true)

	mockRepository.On("GetByUid", "myUid").Return(expectedPayment, nil)

	req := httptest.NewRequest("DELETE", "http://localhost:8080/api/v1/payments/uid/myUid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	mockRepository.AssertExpectations(t)
	assert.Equal(t, http.StatusConflict, resp.StatusCode)
}

func TestDeletePaymentByUid(t *testing.T) {

	router, mockRepository := setUp()
	expectedPayment := expectedPayment("myUid", false)

	mockRepository.On("GetByUid", "myUid").Return(expectedPayment, nil)
	mockRepository.On("Delete", expectedPayment).Return(nil)

	req := httptest.NewRequest("DELETE", "http://localhost:8080/api/v1/payments/uid/myUid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	mockRepository.AssertExpectations(t)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

//
// private functions

func setUp() (*mux.Router, *paymentRepositoryImplMock) {

	var router = mux.NewRouter()
	var mockRepository paymentRepositoryImplMock

	NewPaymentHandler(&mockRepository).Register(router)

	return router, &mockRepository

}

func expectedPayment(uid string, processed bool) *model.Payment {

	payment := model.Payment{
		Uid:           uid,
		AccountOrigin: "myAccountOrigin",
		AccountTarget: "myAccountTarget",
		Amount:        25,
		Date:          now,
		Processed:     processed,
	}

	if processed {
		payment.ProcessedDate = &now
	}

	return &payment

}
