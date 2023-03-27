package cbs

import (
	"encoding/json"
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	url2 "net/url"
)

type observationsResponse struct {
	Context  string        `json:"@odata.context"`
	Value    []observation `json:"value"`
	NextLink string        `json:"@odata.nextLink"`
}

type Observation struct {
	Id             int
	Measure        string
	ValueAttribute string
	Value          float64
	Dimension      string
}

type observation map[string]json.RawMessage

type GetObservationsConfig struct {
	TableId       string
	DimensionName string
	Filter        *string
	Select        *string
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

	url := service.url(fmt.Sprintf("%s/Observations?%s", config.TableId, values.Encode()))

	for {
		observationsResponse := observationsResponse{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           url,
			ResponseModel: &observationsResponse,
		}
		_, _, e := service.httpService.HttpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		for _, observation := range observationsResponse.Value {
			var observation_ Observation

			b, _ := json.Marshal(observation)
			err := json.Unmarshal(b, &observation_)
			if err != nil {
				return nil, errortools.ErrorMessage(err)
			}

			if config.DimensionName != "" {
				d, ok := observation[config.DimensionName]
				if ok {
					err := json.Unmarshal(d, &observation_.Dimension)
					if err != nil {
						return nil, errortools.ErrorMessage(err)
					}
				}
			}

			observations = append(observations, observation_)
		}

		if observationsResponse.NextLink == "" {
			break
		}

		url = observationsResponse.NextLink
	}

	return &observations, nil
}
