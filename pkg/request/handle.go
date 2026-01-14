package request

import (
	"go/adv-demo/pkg/response"
	"net/http"
)

func HandleBody[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)

	if err != nil {
		response.JSON(*w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return nil, err
	}

	err = IsValid(body)

	if err != nil {
		validationErrMsg := "Validation failed. " + err.Error()
		response.JSON(*w, http.StatusBadRequest, map[string]string{"error": validationErrMsg})
		return nil, err
	}

	return body, nil
}
