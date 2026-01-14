package stat

import (
	"go/adv-demo/configs"
	"go/adv-demo/pkg/event"
	"go/adv-demo/pkg/middleware"
	"go/adv-demo/pkg/response"
	"net/http"
	"time"
)

const (
	GroupByDay   = "day"
	GroupByMonth = "month"
)

type StatHandlerDeps struct {
	StatRepository *StatRepository
	*configs.Config
	EventBus *event.EventBus
}

type StatHandler struct {
	LinkRepository *StatRepository
	EventBus       *event.EventBus
	StatRepository *StatRepository
}

func NewLinkHandler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := &StatHandler{LinkRepository: deps.StatRepository, EventBus: deps.EventBus, StatRepository: deps.StatRepository}

	router.Handle("GET /stat", middleware.IsAuthed(handler.GetStat(), deps.Config))
}

func (handler *StatHandler) GetStat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		from, err := time.Parse("2006-01-02", r.URL.Query().Get("from"))
		if err != nil {
			http.Error(w, "Invalid 'from' date", http.StatusBadRequest)
			return
		}

		to, err := time.Parse("2006-01-02", r.URL.Query().Get("to"))
		if err != nil {
			http.Error(w, "Invalid 'to' date", http.StatusBadRequest)
			return
		}

		by := r.URL.Query().Get("by")
		if by == "" {
			http.Error(w, "Invalid 'by' range", http.StatusBadRequest)
			return
		}

		if by != GroupByDay && by != GroupByMonth {
			http.Error(w, "Invalid 'by' range. Must be 'day', 'month', or 'year'", http.StatusBadRequest)
			return
		}

		stats := handler.StatRepository.GetStats(by, from, to)

		response.JSON(w, 200, stats)
	}
}
