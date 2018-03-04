package reset

import (
	"fmt"
	"github.com/viant/endly"
	"github.com/viant/toolbox"
)

//ServiceID represents rest service id.
const ServiceID = "rest/runner"

type restService struct {
	*endly.AbstractService
}

func (s *restService) sendRequest(request *SendRequest) (*SendResponse, error) {
	var resetResponse = make(map[string]interface{})
	err := toolbox.RouteToService(request.Method, request.URL, request.Request, &resetResponse)
	if err != nil {
		return nil, err
	}
	return &SendResponse{
		Response: resetResponse,
	}, nil

}

const restSendExample = `
{
		"URL": "http://127.0.0.1:8085/v1/reporter/register/",
		"Method": "POST",
		"Request": {
			"Report": {
				"Columns": [
					{
						"Alias": "",
						"Name": "category"
					}
				],
				"From": "expenditure",
				"Groups": [
					"year"
				],
				"Name": "report1",
				"Values": [
					{
						"Column": "expenditure",
						"Function": "SUM"
					}
				]
			},
			"ReportType": "pivot"
		}
	}`

func (s *restService) registerRoutes() {
	s.Register(&endly.ServiceActionRoute{
		Action: "send",
		RequestInfo: &endly.ActionInfo{
			Description: "send REST request",
		},
		RequestProvider: func() interface{} {
			return &SendRequest{}
		},
		ResponseProvider: func() interface{} {
			return &SendResponse{}
		},
		Handler: func(context *endly.Context, request interface{}) (interface{}, error) {
			if req, ok := request.(*SendRequest); ok {
				return s.sendRequest(req)
			}
			return nil, fmt.Errorf("unsupported request type: %T", request)
		},
	})
}

//NewRestService creates a new reset service
func New() endly.Service {
	var result = &restService{
		AbstractService: endly.NewAbstractService(ServiceID),
	}
	result.AbstractService.Service = result
	result.registerRoutes()
	return result

}