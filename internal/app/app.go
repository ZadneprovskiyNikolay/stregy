package app

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"stregy/internal/composites"
	"stregy/internal/config"
	"stregy/internal/domain/backtest"
	"stregy/pkg/logging"
	"time"

	"github.com/julienschmidt/httprouter"
)

func Run(cfg *config.Config) {
	logger := logging.GetLogger()

	pgormComposite, err := composites.NewPGormComposite(cfg.PosgreSQL.Host, cfg.PosgreSQL.Port, cfg.PosgreSQL.Username, cfg.PosgreSQL.Password, cfg.PosgreSQL.Database)
	if err != nil {
		logger.Fatal("pgorm composite failed")
	}

	userComposite, err := composites.NewUserComposite(pgormComposite)
	if err != nil {
		logger.Fatal("user composite failed")
	}

	quoteComposite, err := composites.NewQuoteComposite(pgormComposite)
	if err != nil {
		logger.Fatal("quote composite failed")
	}

	tickComposite, err := composites.NewTickComposite(pgormComposite)
	if err != nil {
		logger.Fatal("tick composite failed")
	}

	strategyComposite, err := composites.NewStrategyComposite(userComposite.Service)
	if err != nil {
		logger.Fatal("strategy composite failed")
	}

	exgAccountComposite, err := composites.NewExchangeAccountComposite(pgormComposite, userComposite.Service)
	if err != nil {
		logger.Fatal("exchange account composite failed")
	}

	symbolComposite, err := composites.NewSymbolComposite(pgormComposite)
	if err != nil {
		logger.Fatal("symbol composite failed")
	}

	_, err = composites.NewOrderComposite()
	if err != nil {
		logger.Fatal("order composite failed")
	}

	backtestComposite, err := composites.NewBacktestComposite(
		pgormComposite, exgAccountComposite.Service,
		userComposite.Service, tickComposite.Service,
		quoteComposite.Service, symbolComposite.Service)
	if err != nil {
		logger.Fatal("backtest composite failed")
	}

	if *backtest.BacktestID == "" {
		StartServer(userComposite, strategyComposite, exgAccountComposite, backtestComposite)
	} else {
		err = backtestComposite.Service.Run()
		if err != nil {
			logger.Fatal(err.Error())
		}
	}
}

func StartServer(
	userComposite *composites.UserComposite,
	strategyComposite *composites.StrategyComposite,
	exgAccountComposite *composites.ExchangeAccountComposite,
	backtestComposite *composites.BacktestComposite) {
	logger := logging.GetLogger()

	router := httprouter.New()

	userComposite.Handler.Register(router)
	strategyComposite.Handler.Register(router)
	exgAccountComposite.Handler.Register(router)
	backtestComposite.Handler.Register(router)

	cfg := config.GetConfig()
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", cfg.Listen.BindIP, cfg.Listen.Port))
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Server started")
	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatalln(server.Serve(listener))
}
