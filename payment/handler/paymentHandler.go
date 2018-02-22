package handler

import (
	"net/http"
	"gitgub.com/javierjmgits/go-payment-api/payment/model"
	"gitgub.com/javierjmgits/go-payment-api/base/util"
	"github.com/gorilla/mux"
	"encoding/json"
	"github.com/satori/go.uuid"
	"time"
	"gitgub.com/javierjmgits/go-payment-api/payment/repository"
)

type PaymentHandler struct {
	paymentRepository repository.PaymentRepository
}

type PaymentView struct {
	Uid           string     `json:"uid"`
	AccountOrigin string     `json:"accountOrigin"`
	AccountTarget string     `json:"accountTarget"`
	Amount        float64    `json:"amount"`
	Date          time.Time  `json:"date"`
	Processed     bool       `json:"processed"`
	ProcessedDate *time.Time `json:"processedDate"`
}

type PaymentCreate struct {
	AccountOrigin string    `json:"accountOrigin"`
	AccountTarget string    `json:"accountTarget"`
	Amount        float64   `json:"amount"`
	Date          time.Time `json:"date"`
}

func NewPaymentHandler(paymentRepository repository.PaymentRepository) *PaymentHandler {

	return &PaymentHandler{
		paymentRepository: paymentRepository,
	}
}

func (ph *PaymentHandler) Register(router *mux.Router) {
	router.HandleFunc("/api/v1/payments", ph.GetPayments).Methods("GET")
	router.HandleFunc("/api/v1/payments/uid/{uid}", ph.GetPaymentByUid).Methods("GET")
	router.HandleFunc("/api/v1/payments", ph.CreatePayment).Methods("POST")
	router.HandleFunc("/api/v1/payments/uid/{uid}/processed", ph.FlagPaymentAsProcessedByUid).Methods("PATCH")
	router.HandleFunc("/api/v1/payments/uid/{uid}", ph.DeletePaymentByUid).Methods("DELETE")
}

func (ph *PaymentHandler) GetPayments(w http.ResponseWriter, r *http.Request) {

	payments, errorDB := ph.paymentRepository.GetAll()

	if errorDB != nil {
		util.WriteError(w, 500, errorDB.Error())
		return
	}

	util.WritePayload(w, 200, newPaymentViews(payments))
}

func (ph *PaymentHandler) GetPaymentByUid(w http.ResponseWriter, r *http.Request) {

	payment, errorDB := ph.getPaymentByUid(r)

	if errorDB != nil {
		util.WriteError(w, 404, errorDB.Error())
		return
	}

	util.WritePayload(w, 200, newPaymentView(payment))

}

func (ph *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {

	var paymentCreate PaymentCreate
	errorJson := json.NewDecoder(r.Body).Decode(&paymentCreate)
	if errorJson != nil {
		util.WriteError(w, 400, errorJson.Error())
		return
	}

	payment, errorUuid := newPayment(&paymentCreate)

	if errorUuid != nil {
		util.WriteError(w, 500, errorUuid.Error())
		return
	}

	payment, errorDB := ph.paymentRepository.Create(payment)

	if errorDB != nil {
		util.WriteError(w, 500, errorDB.Error())
		return
	}

	util.WritePayload(w, 201, newPaymentView(payment))
}

func (ph *PaymentHandler) FlagPaymentAsProcessedByUid(w http.ResponseWriter, r *http.Request) {

	payment, errorFind := ph.getPaymentByUid(r)

	if errorFind != nil {
		util.WriteError(w, 404, errorFind.Error())
		return
	}

	if payment.Processed {
		util.WriteError(w, 406, "Payment already processed")
		return
	}

	now := time.Now().UTC().Truncate(time.Second)

	payment.Processed = true
	payment.ProcessedDate = &now

	payment, errorDB := ph.paymentRepository.Update(payment)

	if errorDB != nil {
		util.WriteError(w, 500, errorDB.Error())
		return
	}

	util.WritePayload(w, 200, newPaymentView(payment))

}

func (ph *PaymentHandler) DeletePaymentByUid(w http.ResponseWriter, r *http.Request) {

	payment, errorFind := ph.getPaymentByUid(r)

	if errorFind != nil {
		util.WriteError(w, 404, errorFind.Error())
		return
	}

	if payment.Processed {
		util.WriteError(w, 406, "Payment already processed")
		return
	}

	errorDB := ph.paymentRepository.Delete(payment)

	if errorDB != nil {
		util.WriteError(w, 500, errorDB.Error())
		return
	}

	util.WritePayload(w, 204, map[string]string{})
}

//
// private functions

func (ph *PaymentHandler) getPaymentByUid(r *http.Request) (*model.Payment, error) {

	uid := mux.Vars(r)["uid"]

	payment, errorDB := ph.paymentRepository.GetByUid(uid)

	if errorDB != nil {
		return nil, errorDB
	}

	return payment, nil

}

func newPayment(paymentCreate *PaymentCreate) (*model.Payment, error) {

	uuidResult, errorUuid := uuid.NewV4()

	if errorUuid != nil {
		return nil, errorUuid
	}

	return &model.Payment{
		Uid:           uuidResult.String(),
		AccountOrigin: paymentCreate.AccountOrigin,
		AccountTarget: paymentCreate.AccountTarget,
		Amount:        paymentCreate.Amount,
		Date:          paymentCreate.Date,
		Processed:     false,
	}, nil

}

func newPaymentView(payment *model.Payment) *PaymentView {

	return &PaymentView{
		Uid:           payment.Uid,
		AccountOrigin: payment.AccountOrigin,
		AccountTarget: payment.AccountTarget,
		Amount:        payment.Amount,
		Date:          payment.Date,
		Processed:     payment.Processed,
		ProcessedDate: payment.ProcessedDate,
	}
}

func newPaymentViews(payments []model.Payment) []PaymentView {

	var results []PaymentView

	for _, item := range payments {

		payment := newPaymentView(&item)

		results = append(results, *payment)
	}

	return results
}
