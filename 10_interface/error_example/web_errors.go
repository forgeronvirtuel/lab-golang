package errorexample

import (
	"fmt"
	"net/http"
)

type WebServiceError struct {
	Status     int    `json:"status"`
	StatusText string `json:"statustext"`
}

func (n WebServiceError) Error() string {
	return fmt.Sprintf("%d : %s", n.Status, n.StatusText)
}

func NewWebServiceError(status int) *WebServiceError {
	return &WebServiceError{
		Status:     status,
		StatusText: http.StatusText(status),
	}
}

type NotFoundError struct {
	Resource string `json:"ressource"`
	*WebServiceError
}

func (n NotFoundError) Error() string {
	return fmt.Sprintf("%s, %s", n.WebServiceError.Error(), n.Resource)
}

func NewNotFoundError(resource string) *NotFoundError {
	return &NotFoundError{
		Resource:        resource,
		WebServiceError: NewWebServiceError(http.StatusNotFound),
	}
}

func DoStuff(path string) error {
	if path == "/" {
		// Do very important thing
		return nil
	}
	return NewNotFoundError(path)
}
