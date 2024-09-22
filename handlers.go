package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/kumersun/bnovo/entity"
	"github.com/kumersun/bnovo/validator"
)

type Response struct {
	Data  any    `json:"data"`
	Error string `json:"error"`
}

func NewResponse() *Response {
	return &Response{}
}

func sendResponse(response *Response, w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Debug-Time", getDebugTime())
	w.Header().Set("X-Debug-Memory", getDebugMemory())
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func sendSuccessResponse(w http.ResponseWriter, data any) {
	response := NewResponse()
	response.Data = data
	sendResponse(response, w, http.StatusOK)
}

func sendErrorResponse(w http.ResponseWriter, err error, code int) {
	response := NewResponse()
	response.Error = err.Error()
	sendResponse(response, w, code)
}

func decodeGuest(w http.ResponseWriter, r *http.Request) (*entity.Guest, error) {
	var guest entity.Guest

	if err := json.NewDecoder(r.Body).Decode(&guest); err != nil {
		return nil, err
	}

	prepareGuest(&guest)
	guestValidator := validator.NewGuestValidator(&guest)
	err := guestValidator.Validate()
	if err != nil {
		return nil, err
	}

	return &guest, nil
}

func prepareGuest(guest *entity.Guest) {
	guest.Name = strings.TrimSpace(guest.Name)
	guest.Surname = strings.TrimSpace(guest.Surname)
	guest.Phone = strings.TrimSpace(guest.Phone)
	guest.Email = strings.TrimSpace(guest.Email)
	guest.Country = strings.TrimSpace(guest.Country)
}

func createGuest(
	w http.ResponseWriter,
	r *http.Request,
) {
	newGuest, err := decodeGuest(w, r)
	if err != nil {
		sendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := guestRepo.CreateGuest(newGuest); err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, newGuest)
}

func getGuests(
	w http.ResponseWriter,
	r *http.Request,
) {
	guests, err := guestRepo.GetGuests(r.Context())
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, guests)
}

func getGuest(
	w http.ResponseWriter,
	r *http.Request,
) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendErrorResponse(w, fmt.Errorf("Invalid guest ID: %d", id), http.StatusBadRequest)
		return
	}

	guest, err := guestRepo.GetGuest(id)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if guest == nil {
		sendErrorResponse(w, errors.New("Guest not found"), http.StatusBadRequest)
		return
	}

	sendSuccessResponse(w, guest)
}

func updateGuest(
	w http.ResponseWriter,
	r *http.Request,
) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendErrorResponse(w, errors.New("Invalid guest ID"), http.StatusBadRequest)
		return
	}

	updatedGuest, err := decodeGuest(w, r)
	if err != nil {
		sendErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	updatedGuest.ID = id

	if err := guestRepo.UpdateGuest(updatedGuest); err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, updatedGuest)
}

func deleteGuest(
	w http.ResponseWriter,
	r *http.Request,
) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendErrorResponse(w, errors.New("Invalid guest ID"), http.StatusBadRequest)
		return
	}

	if err := guestRepo.DeleteGuest(id); err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	var response [1]int = [1]int{id}
	sendSuccessResponse(w, response)
}
