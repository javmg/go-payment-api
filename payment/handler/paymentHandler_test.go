package handler

import (
	"net/http/httptest"
	"testing"
	"io/ioutil"
	"gitgub.com/javierjmgits/go-payment-api/payment/model"
	"time"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
)

//
// mock data

var now = time.Now()

var expectedPayment = model.Payment{
	Uid:           "myUid",
	AccountOrigin: "myAccountOrigin",
	AccountTarget: "myAccountTarget",
	Amount:        25,
	Date:          now,
	Processed:     true,
	ProcessedDate: &now,
}

//
// mocks

type paymentRepositoryImplMock struct {
}

func (mock *paymentRepositoryImplMock) GetAll() ([]model.Payment, error) {
	return []model.Payment{expectedPayment}, nil
}

func (mock *paymentRepositoryImplMock) GetByUid(uid string) (*model.Payment, error) {
	if uid == "unknown" {
		return nil, errors.New("no payment found")
	}
	return &expectedPayment, nil
}

func (mock *paymentRepositoryImplMock) Create(payment *model.Payment) (*model.Payment, error) {
	return payment, nil
}

func (mock *paymentRepositoryImplMock) Update(payment *model.Payment) (*model.Payment, error) {
	return payment, nil
}

func (mock *paymentRepositoryImplMock) Delete(payment *model.Payment) (error) {
	return nil
}

//
// testing classes

var router = mux.NewRouter()

func TestMain(m *testing.M) {

	mock := paymentRepositoryImplMock{}

	paymentHandler := NewPaymentHandler(&mock)

	paymentHandler.Register(router)

	m.Run()
}

//
// tests

func TestGetPayments(t *testing.T) {

	req := httptest.NewRequest("GET", "http://localhost:8080/api/v1/payments", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	var payments []PaymentView

	json.Unmarshal(body, &payments)

	// verify

	if len(payments) != 1 {
		t.Fatal("Expecting a single payment")
	}
}

func TestGetPaymentByUid(t *testing.T) {

	req := httptest.NewRequest("GET", "http://localhost:8080/api/v1/payments/uid/myUid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	var payment PaymentView

	json.Unmarshal(body, &payment)

	// verify

	if payment.Uid != expectedPayment.Uid {
		t.Fatalf("Expecting a payment with uid: %v", expectedPayment.Uid)
	}

	if payment.AccountOrigin != expectedPayment.AccountOrigin {
		t.Fatalf("Expecting a payment with account origin: %v", expectedPayment.AccountOrigin)
	}

	if payment.AccountTarget != expectedPayment.AccountTarget {
		t.Fatalf("Expecting a payment with account target: %v", expectedPayment.AccountTarget)
	}

	if payment.Amount != expectedPayment.Amount {
		t.Fatalf("Expecting a payment with amount: %v", expectedPayment.Amount)
	}

	if payment.Date != expectedPayment.Date {
		t.Fatalf("Expecting a payment with date: %v", expectedPayment.Date)
	}

	if payment.Processed != expectedPayment.Processed {
		t.Fatalf("Expecting a payment with processed: %v", expectedPayment.Processed)
	}

	if &payment.ProcessedDate == &expectedPayment.ProcessedDate {
		t.Fatalf("Expecting a payment with processed date: %v", &expectedPayment.ProcessedDate)
	}

}

func TestGetPaymentByUidKoNotFound(t *testing.T) {

	req := httptest.NewRequest("GET", "http://localhost:8080/api/v1/payments/uid/unknown", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != 404 {
		t.Fatal("Expecting a 404")
	}
}
