package main

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string      `json:"message"`
	Error   bool        `json:"error"`
	Data    interface{} `json:"data"`
}

type ApiHandler struct {
	service *APiService
}

func NewHandler(svc *APiService) *ApiHandler {
	return &ApiHandler{service: svc}
}

func (handler *ApiHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var createUserRequest struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&createUserRequest); err != nil {
		writeResponse(w, 400, &Response{
			Message: "invalid json body" + err.Error(),
			Error:   true,
			Data:    nil,
		})
		return
	}
	user, err := handler.service.CreateUser(createUserRequest.Name)
	if err != nil {
		writeResponse(w, 500, &Response{
			Message: "failed to create user " + err.Error(),
			Error:   true,
			Data:    nil,
		})
		return
	}
	ok(w, &Response{
		Message: "user created",
		Error:   false,
		Data:    user,
	})
}

func (handler *ApiHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var createTxRequest struct {
		UserId     int64 `json:"user_id"`
		Amount     int64 `json:"amount"`
		ReceiverId int64 `json:"receiver_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&createTxRequest); err != nil {
		writeResponse(w, 400, &Response{
			Message: "invalid json body" + err.Error(),
			Error:   true,
			Data:    nil,
		})
		return
	}
	tx, err := handler.service.CreateTransaction(createTxRequest.UserId, createTxRequest.ReceiverId, createTxRequest.Amount)
	if err != nil {
		writeResponse(w, 500, &Response{
			Message: "failed to create transaction " + err.Error(),
			Error:   true,
			Data:    nil,
		})
		return
	}
	ok(w, &Response{
		Message: "transaction created",
		Error:   false,
		Data:    tx,
	})
}

func (handler *ApiHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	data, err := handler.service.GetUsers()
	if err != nil {
		writeResponse(w, 500, &Response{
			Message: "failed to fetch users " + err.Error(),
			Error:   true,
			Data:    nil,
		})
		return
	}
	ok(w, &Response{
		Message: "users",
		Error:   false,
		Data:    data,
	})
}

func writeResponse(w http.ResponseWriter, statusCode int, data *Response) {
	j, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(j)
}

func ok(w http.ResponseWriter, data *Response) {
	writeResponse(w, 200, data)
}
