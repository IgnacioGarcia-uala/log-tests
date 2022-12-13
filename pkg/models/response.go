package models

type Response struct {
	Status      string  `json:"status"`
	LimitPeriod string  `json:"limitPeriod,omitempty"`
	LimitAmount float64 `json:"limitAmount,omitempty"`
	Excess      float64 `json:"excess,omitempty"`
}

func ApproveResponse() Response {
	return Response{
		Status: "APPROVED",
	}
}

func RejectedResponse(period string, amount, excess float64) Response {
	return Response{
		Status:      "REJECTED",
		LimitPeriod: period,
		LimitAmount: amount,
		Excess:      excess,
	}
}
