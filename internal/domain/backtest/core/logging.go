package core

func (b *Backtest) Print(s string) {
	b.logger.Print(s)
}

func (b *Backtest) Printf(format string, v ...interface{}) {
	b.logger.Printf(format, v...)
}

func (b *Backtest) PrintfWithoutPrefix(format string, v ...interface{}) {
	b.logger.Config.OmmitTimePrefix = true
	b.logger.Printf(format, v...)
	b.logger.Config.OmmitTimePrefix = false
}
