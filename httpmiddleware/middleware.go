package httpmiddleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Middleware struct{}

func New() *Middleware {
	return &Middleware{}
}

func (m *Middleware) HandlerError(h func(resp http.ResponseWriter, req *http.Request) error) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		err := h(resp, req.WithContext(ctx))
		if err != nil {
			httpStatus := httpStatusCode(err)
			message := err.Error()

			fmt.Println(httpStatus, message)
			resp.WriteHeader(httpStatus)

			errorResponse := ErrorResponse{
				Timestamp: time.Now(),
				Message:   message,
			}
			if httpStatus == 400 {
				switch e := err.(type) {
				case badrequest:
					errorResponse.Fields = make([]Field, 0)
					for key, v := range e.GetFields() {
						f := Field{Name: key, Message: v.Error()}
						errorResponse.Fields = append(errorResponse.Fields, f)
					}

				}
			}
			SendJSON(resp, errorResponse)
		}
	})
}

func SendJSON(resp http.ResponseWriter, payload interface{}) error {
	resp.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(resp).Encode(payload); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return errors.New("Error to decode")
	}
	return nil
}
