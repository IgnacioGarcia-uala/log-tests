package models

type Limits struct {
	Type   string   `json:"type"`
	Config []Config `json:"config"`
}

type Config struct {
	Amount float64 `json:"amount"`
	Period string  `json:"period"`
}

func GetLimits() Limits {
	return Limits{
		Type: "P2P",
		Config: []Config{
			{
				Amount: 1000,
				Period: "PT24H",
			},
			{
				Amount: 10000,
				Period: "P28D",
			},
		},
	}
}
