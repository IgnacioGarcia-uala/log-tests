package main

import (
	"github.com/Bancar/uala-labssupport-monitoreofraude/internal/processor"
	"github.com/Bancar/uala-labssupport-monitoreofraude/pkg/handler"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	p := processor.New()
	h := handler.New(p)

	lambda.Start(h.Handle)
}
