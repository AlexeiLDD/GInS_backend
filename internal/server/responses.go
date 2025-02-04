package server

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	StatusOk      = 200
	StatusCreated = 201

	StatusBadRequest   = 400
	StatusUnauthorized = 401
	StatusForbidden    = 403
	StatusNotFound     = 404
	StatusNotAllowed   = 405

	StatusInternalServerError = 500
)

const (
	ErrUserAlreadyExists  = "User with same email already exists"
	ErrWrongCredentials   = "Wrong credentials" //nolint:gosec
	ErrUnauthorized       = "User not authorized"
	ErrAuthorized         = "User already authorized"
	ErrDifferentPasswords = "Passwords do not match"
	ErrNotValidData       = "User data is not valid"

	ErrAdvertNotExist = "Advert does not exist"

	ErrInternalServer = "Server error"
	ErrBadRequest     = "Bad request"
	ErrNotAllowed     = "Method not allowed"
	ErrForbidden      = "User have no access to this content"
)

type ErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

type OkResponse struct {
	Code  int `json:"code"`
	Items any `json:"items"`
}

func sendResponse(writer http.ResponseWriter, serverResponse []byte) {
	_, err := writer.Write(serverResponse)
	if err != nil {
		log.Println(err)
		http.Error(writer, ErrInternalServer, StatusInternalServerError)

		return
	}
}

func SendOkResponse(writer http.ResponseWriter, response interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(StatusOk)

	serverResponse, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		http.Error(writer, ErrInternalServer, StatusInternalServerError)

		return
	}

	sendResponse(writer, serverResponse)
}

func SendErrResponse(request *http.Request, writer http.ResponseWriter, response *ErrResponse) {
	code := request.Context().Value("code").(*int)
	*code = response.Code

	serverResponse, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		http.Error(writer, ErrInternalServer, StatusInternalServerError)

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	sendResponse(writer, serverResponse)
}

func NewErrResponse(code int, status string) *ErrResponse {
	return &ErrResponse{
		Code:   code,
		Status: status,
	}
}

func NewOkResponse(items any) *OkResponse {
	return &OkResponse{
		Code:  StatusOk,
		Items: items,
	}
}
