package bt

import (
	"net/http"
	"stregy/internal/adapters/api"
	userapi "stregy/internal/adapters/api/user"
	"stregy/internal/config"
	"stregy/internal/domain/backtest"
	"stregy/internal/domain/exgaccount"
	"stregy/internal/domain/user"
	"stregy/pkg/handlers"
	"stregy/pkg/logging"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mitchellh/mapstructure"
)

const (
	backtestURL = "/api/backtest"
)

type handler struct {
	backtestService backtest.Service
	exgAccService   exgaccount.Service
	userService     user.Service
}

func NewHandler(
	backtestService backtest.Service,
	userService user.Service,
) api.Handler {
	return &handler{
		backtestService: backtestService,
		userService:     userService,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	createSEHandler := handlers.JsonHandler(h.backtestHandler, &BacktestDTO{})
	userHandler := userapi.AuthenticationHandler(createSEHandler, h.userService)
	router.POST(backtestURL, handlers.ToSimpleHandler(userHandler))
}

func (h *handler) backtestHandler(
	w http.ResponseWriter,
	r *http.Request,
	params httprouter.Params,
	args map[string]interface{},
) {
	cfg := config.GetConfig()
	logger := logging.GetLogger()
	// Parse and validate request.
	user := user.User{}
	mapstructure.Decode(args["user"], &user)
	apiDTO := BacktestDTO{}
	mapstructure.Decode(args["json"], &apiDTO)

	// Create db record.
	startDate, _ := time.Parse("2006-01-02 15:04:05", apiDTO.StartDate)
	endDate, _ := time.Parse("2006-01-02 15:04:05", apiDTO.EndDate)
	domainDTO := backtest.BacktestDTO{
		StrategyName: apiDTO.StrategyName,
		SymbolName:   apiDTO.Symbol,
		StartDate:    startDate,
		EndDate:      endDate,
	}
	btDomain, err := h.backtestService.Create(domainDTO)
	if err != nil {
		logger.Error(err.Error())
		handlers.ReturnError(w, http.StatusInternalServerError, "")
		return
	}

	// Execute.
	if !*cfg.IsDebug {
		err = h.backtestService.Launch(btDomain)
		if err != nil {
			logger.Error(err.Error())
			handlers.ReturnError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	handlers.JsonResponseWriter(w, map[string]string{"backtest_id": btDomain.ID})
}
