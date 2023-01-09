package handler

import (
	"context"

	"github.com/Bancar/uala-labssupport-monitoreofraude/internal/processor"
	"github.com/Bancar/uala-labssupport-monitoreofraude/pkg/handler"
	"github.com/Bancar/uala-labssupport-monitoreofraude/pkg/models"
	"github.com/aws/aws-lambda-go/lambdacontext"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	ct = &lambdacontext.LambdaContext{
		AwsRequestID:       "awsRequestId1234",
		InvokedFunctionArn: "arn:aws:lambda:xxx",
		Identity:           lambdacontext.CognitoIdentity{},
		ClientContext:      lambdacontext.ClientContext{},
	}
	ctx = lambdacontext.NewContext(context.TODO(), ct)
)

var _ = Describe("Handler", func() {
	Context("Test", happyPath)
})

func happyPath() {
	var (
		p = processor.New()
		h = handler.New(p)
	)

	It("Main", func() {
		req := models.Request{
			AccountId: "abc-123",
			Type:      "CASH_OUT_CVU",
			Amount:    125,
		}

		_, err := h.Handle(ctx, req)
		Ω(err).ShouldNot(HaveOccurred())
	})
}
