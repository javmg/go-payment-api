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
	"strings"
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
		util.WriteError(w, http.StatusInternalServerError, errorDB.Error())
		return
	}

	util.WritePayload(w, http.StatusOK, newPaymentViews(payments))
}

func (ph *PaymentHandler) GetPaymentByUid(w http.ResponseWriter, r *http.Request) {

	payment, responseGenerated := ph.getAndCheckPaymentByUid(w, r, false)

	if responseGenerated {
		return
	}

	util.WritePayload(w, http.StatusOK, newPaymentView(payment))

}

func (ph *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {

	paymentCreate, responseGenerated := decodeAndValidatePaymentCreate(w, r)

	if responseGenerated {
		return
	}

	paymentToSave, errorUuid := newPayment(paymentCreate)

	if errorUuid != nil {
		util.WriteError(w, http.StatusInternalServerError, errorUuid.Error())
		return
	}

	paymentSaved, errorDB := ph.paymentRepository.Create(paymentToSave)

	if errorDB != nil {
		util.WriteError(w, http.StatusInternalServerError, errorDB.Error())
		return
	}

	util.WritePayload(w, http.StatusCreated, newPaymentView(paymentSaved))
}

func (ph *PaymentHandler) FlagPaymentAsProcessedByUid(w http.ResponseWriter, r *http.Request) {

	payment, responseGenerated := ph.getAndCheckPaymentByUid(w, r, true)

	if responseGenerated {
		return
	}

	now := time.Now().UTC().Truncate(time.Second)

	payment.Processed = true
	payment.ProcessedDate = &now

	payment, errorDB := ph.paymentRepository.Update(payment)

	if errorDB != nil {
		util.WriteError(w, http.StatusInternalServerError, errorDB.Error())
		return
	}

	util.WritePayload(w, http.StatusOK, newPaymentView(payment))

}

func (ph *PaymentHandler) DeletePaymentByUid(w http.ResponseWriter, r *http.Request) {

	payment, responseGenerated := ph.getAndCheckPaymentByUid(w, r, true)

	if responseGenerated {
		return
	}

	errorDB := ph.paymentRepository.Delete(payment)

	if errorDB != nil {
		util.WriteError(w, http.StatusInternalServerError, errorDB.Error())
		return
	}

	util.WritePayload(w, http.StatusNoContent, map[string]string{})
}

//
// private functions

func (ph *PaymentHandler) getAndCheckPaymentByUid(w http.ResponseWriter, r *http.Request, ensureNotProcessed bool) (payment *model.Payment, responseGenerated bool) {

	uid := mux.Vars(r)["uid"]

	payment, errorDB := ph.paymentRepository.GetByUid(uid)

	if errorDB != nil {

		if strings.Contains(errorDB.Error(), "not found") {
			util.WriteError(w, http.StatusNotFound, errorDB.Error())

		} else {
			util.WriteError(w, http.StatusInternalServerError, errorDB.Error())
		}

		return nil, true

	}

	if ensureNotProcessed && payment.Processed {
		util.WriteError(w, http.StatusConflict, "Payment already processed")
		return nil, true
	}

	return payment, false
}

func decodeAndValidatePaymentCreate(w http.ResponseWriter, r *http.Request) (paymentCreate *PaymentCreate, responseGenerated bool) {

	errorJson := json.NewDecoder(r.Body).Decode(&paymentCreate)

	if errorJson != nil {
		util.WriteError(w, http.StatusBadRequest, errorJson.Error())
		return nil, true
	}

	if paymentCreate.AccountOrigin == "" {
		util.WriteError(w, http.StatusBadRequest, "account origin is mandatory")
		return nil, true
	}

	if paymentCreate.AccountTarget == "" {
		util.WriteError(w, http.StatusBadRequest, "account target is mandatory")
		return nil, true
	}

	if paymentCreate.Amount <= 0 {
		util.WriteError(w, http.StatusBadRequest, "amount must be a positive number")
		return nil, true
	}

	return paymentCreate, false
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
