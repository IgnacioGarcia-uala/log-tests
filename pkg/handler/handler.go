package handler

import (
	"context"
	"log"

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
	log.Printf("<start> <handler> <handle> - Handling Request %+v", request)

	r, err := h.p.Process(request)
	if err != nil {
		log.Printf("<middle> <handler> <handle> - Error processing: %v", err.Error())
		return models.Response{}, err
	}

	log.Printf("<finish> <handler> <handle> - Response: %+v", r)
	return r, nil
}
