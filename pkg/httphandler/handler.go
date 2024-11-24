package httphandler

import (
	"context"
	"net/http"

	"github.com/L4B0MB4/EVTSRC/pkg/httphandler/controller"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type HttpHandler struct {
	httpServer      *http.Server
	router          *gin.Engine
	eventController *controller.EventController
}

func NewHttpHandler(c *controller.EventController) *HttpHandler {
	r := gin.Default()
	srv := &http.Server{
		Addr:    "0.0.0.0" + ":" + "5515",
		Handler: r,
	}
	handler := &HttpHandler{
		router:          r,
		httpServer:      srv,
		eventController: c,
	}

	handler.RegisterRoutes()

	return handler
}

func (h *HttpHandler) RegisterRoutes() {
	h.router.GET("aggregates/:aggregateId/events", h.eventController.GetEventsForAggregate)
	h.router.POST("aggregates/:aggregateId/events", h.eventController.AddEventToAggregate)
	h.router.GET("/events/:eventId/since", h.eventController.GetEventsSince)
}

func (h *HttpHandler) Start() error {
	return h.httpServer.ListenAndServe()
}

func (h *HttpHandler) Stop() {
	err := h.httpServer.Shutdown(context.Background())
	if err != nil {
		log.Warn().Err(err).Msg("Error during reading response body")
	}
}
