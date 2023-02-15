package bt

type BacktestDTO struct {
	StrategyName string `json:"strategy_name"`
	Symbol       string `json:"symbol"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
}
