package handler

import "github.com/gorilla/mux"

type BaseHandler interface {
	Register(router *mux.Router)
}
