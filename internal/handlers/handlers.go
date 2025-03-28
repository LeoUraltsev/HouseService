package handlers

import (
	"net/http"

	"github.com/LeoUraltsev/HauseService/internal/gen"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

// GetDummyLogin implements gen.ServerInterface.
func (h *Handler) GetDummyLogin(w http.ResponseWriter, r *http.Request, params gen.GetDummyLoginParams) {
	panic("unimplemented")
}

// GetHouseId implements gen.ServerInterface.
func (h *Handler) GetHouseId(w http.ResponseWriter, r *http.Request, id gen.HouseId) {
	panic("unimplemented")
}

// PostFlatCreate implements gen.ServerInterface.
func (h *Handler) PostFlatCreate(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostFlatUpdate implements gen.ServerInterface.
func (h *Handler) PostFlatUpdate(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostHouseCreate implements gen.ServerInterface.
func (h *Handler) PostHouseCreate(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostHouseIdSubscribe implements gen.ServerInterface.
func (h *Handler) PostHouseIdSubscribe(w http.ResponseWriter, r *http.Request, id gen.HouseId) {
	panic("unimplemented")
}

// PostLogin implements gen.ServerInterface.
func (h *Handler) PostLogin(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostRegister implements gen.ServerInterface.
func (h *Handler) PostRegister(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
