package handler

import (
	"context"

	"github.com/Bancar/uala-labssupport-monitoreofraude/internal/logger"
	"github.com/Bancar/uala-labssupport-monitoreofraude/pkg/models"
)

type Handler interface {
	Handle(ctx context.Context, request models.Request) (models.Response, error)
}

type Processor interface {
	Process(request models.Request) (models.Response, error)
}

type LogTestHandler struct {
	p Processor
}

func New(p Processor) Handler {
	return &LogTestHandler{
		p,
	}
}

func (h LogTestHandler) Handle(ctx context.Context, request models.Request) (models.Response, error) {
	log := logger.SetLogger(ctx, request.AccountId)

	log.Info("<start> <handler> <handle> - Handling Request", log.AnyField("request", request))

	r, err := h.p.Process(request)
	if err != nil {
		log.Info("<middle> <handler> <handle> - Error processing", log.AnyField("error", err.Error()))
		return models.Response{}, err
	}

	log.Info("<finish> <handler> <handle> - Response", log.AnyField("response", r))
	return r, nil
}
