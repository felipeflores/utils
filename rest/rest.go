package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/felipeflores/utils/ferrors"
)

// ReadJSON decode JSON from body to payload
func ReadJSON(ctx context.Context, r *http.Request, payload interface{}) error {
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		return ferrors.NewBadRequest(errors.New(fmt.Sprintf("Err: %v message: bad json format", err)))
	}

	return nil
}

// DeserializeJSON is a helper function to read a JSON from request body.
// PS.: payload needs to be a pointer.
func DeserializeJSON(r *http.Request, payload interface{}) error {
	return ReadJSON(context.Background(), r, payload)
}

func RequestBody[REQ Request](w http.ResponseWriter, r *http.Request) (REQ, error) {
	var body REQ

	err := DeserializeJSON(r, &body)
	if err != nil {
		// hpa.l.Error(fmt.Sprintf("[HandlerPersisterAdapter][Insert]Init %v", err))
		return body, err
	}

	if err := body.Validate(); err != nil {
		return body, ferrors.NewBadRequest(err)
	}

	return body, nil
}
