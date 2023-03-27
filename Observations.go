package cbs

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	url2 "net/url"
)

type ObservationsResponse struct {
	Context  string        `json:"@odata.context"`
	Value    []Observation `json:"value"`
	NextLink string        `json:"@odata.nextLink"`
}

type Observation struct {
	Id             int     `json:"Id"`
	Measure        string  `json:"Measure"`
	ValueAttribute string  `json:"ValueAttribute"`
	Value          float64 `json:"Value"`
	Dimension      string  `json:"-"`
}

type GetObservationsConfig struct {
	Filter *string
	Select *string
}

func (service *Service) GetObservations(config *GetObservationsConfig) (*[]Observation, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	var values = url2.Values{}
	if config != nil {
		if config.Filter != nil {
			values.Set("$filter", *config.Filter)
		}
		if config.Select != nil {
			values.Set("$select", *config.Select)
		}
	}

	var observations []Observation

	url := service.url(fmt.Sprintf("%s/Observations?", values.Encode()))

	for {
		observationsResponse := ObservationsResponse{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           url,
			ResponseModel: &observationsResponse,
		}
		_, _, e := service.httpService.HttpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		observations = append(observations, observationsResponse.Value...)

		if observationsResponse.NextLink == "" {
			break
		}

		url = observationsResponse.NextLink
	}

	return &observations, nil
}
