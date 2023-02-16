package core

func (b *Backtest) WriteReportToLog() {
	b.PrintfWithoutPrefix("Commission: %.0f", b.TotalCommission)
	b.PrintfWithoutPrefix("Orders: %d", len(b.OrderHistory))
	b.PrintfWithoutPrefix("Orders executed: %d", b.TotalOrdersExecuted)
	b.PrintfWithoutPrefix("Volume traded: %.0f", b.TotalVolumeTraded)
}
