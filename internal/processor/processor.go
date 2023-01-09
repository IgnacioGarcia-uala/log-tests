package processor

import (
	"context"
	"fmt"

	"github.com/Bancar/uala-labssupport-monitoreofraude/internal/logger"
	"github.com/Bancar/uala-labssupport-monitoreofraude/pkg/handler"
	"github.com/Bancar/uala-labssupport-monitoreofraude/pkg/models"
)

type Processor struct{}

func New() handler.Processor {
	return &Processor{}
}

func (p Processor) Process(request models.Request) (models.Response, error) {
	log := logger.GetLogger(context.TODO())

	log.Info("<start> <processor> <process> - Processing Log Test", log.AnyField("request", request))
	var response models.Response

	log.Info("<middle> <processor> <process> - Getting limit config")
	limits := models.GetLimits()
	log.Info("<middle> <processor> <process> - Limit Config Obtained", log.AnyField("limits", limits))

	for _, limit := range limits.Config {
		log.Info("<middle> <processor> <process> - Evaluating limit config", log.AnyField("limitConfiguration", limit))
		if request.Amount > limit.Amount {
			log.Info(fmt.Sprintf("<middle> <processor> <process> - The Amount %v Excess the limit: %v", request.Amount, limit.Amount))
			return models.RejectedResponse("PT24H", request.Amount, request.Amount-limit.Amount), nil
		}
	}

	response = models.ApproveResponse()
	log.Info("<end> <processor> <process> - Successfully Processed. Response.", log.AnyField("response", response))
	return response, nil
}
