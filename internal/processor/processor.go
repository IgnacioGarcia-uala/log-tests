package processor

import (
	"log"

	"github.com/Bancar/uala-labssupport-monitoreofraude/pkg/handler"
	"github.com/Bancar/uala-labssupport-monitoreofraude/pkg/models"
)

type Processor struct{}

func New() handler.Processor {
	return &Processor{}
}

func (p Processor) Process(request models.Request) (models.Response, error) {
	log.Printf("<start> <processor> <process> - Processing Log Test. Request: %+v", request)
	var response models.Response

	log.Printf("<middle> <processor> <process> - Getting limit config for %v", request.Type)
	limits := models.GetLimits()
	log.Printf("<middle> <processor> <process> - Limit Config Obtained %+v", limits)

	for _, limit := range limits.Config {
		log.Printf("<middle> <processor> <process> - Evaluating limit config %+v", limit)
		if request.Amount > limit.Amount {
			log.Printf("<middle> <processor> <process> - The Amount %v Excess the limit: %v", request.Amount, limit.Amount)
			return models.RejectedResponse("PT24H", request.Amount, request.Amount-limit.Amount), nil
		}
	}

	response = models.ApproveResponse()
	log.Printf("<end> <processor> <process> - Successfully Processed. Response %+v", response)
	return response, nil
}
