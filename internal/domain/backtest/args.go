package backtest

import "flag"

var BacktestID = flag.String("backtest_id", "", "")
var Balance = flag.Float64("balance", 10000, "")
var CommissionType = flag.String("commission_template", "binance", "available types: \"binance\"")
var Commission = flag.Float64("commission", 0, "")
var ReportLocation = flag.String("report_location", "", "")
