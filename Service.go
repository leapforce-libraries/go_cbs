package cbs

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
)

const (
	apiName string = "CBS"
	apiUrl  string = "https://datasets.cbs.nl/odata/v1/CBS"
)

type Service struct {
	httpService *go_http.Service
}

type ServiceConfig struct {
}

func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	accept := go_http.AcceptJson
	httpService, e := go_http.NewService(&go_http.ServiceConfig{
		Accept:     &accept,
		HttpClient: &http.Client{},
	})
	if e != nil {
		return nil, e
	}

	return &Service{httpService}, nil
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", apiUrl, path)
}

func (service Service) ApiName() string {
	return apiName
}

func (service Service) ApiKey() string {
	return ""
}

func (service Service) ApiCallCount() int64 {
	return service.httpService.RequestCount()
}

func (service *Service) ApiReset() {
	service.httpService.ResetRequestCount()
}
