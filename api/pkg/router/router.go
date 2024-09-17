package router

import (
	_ "github.com/bovinxx/code-processor/api/docs"
	"github.com/bovinxx/code-processor/api/pkg/handlers"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Router struct {
	Router *chi.Mux
}

func NewRouter(handler handlers.Handler) *Router {
	router := chi.NewRouter()
	router.Get("/swagger/*", httpSwagger.WrapHandler)
	router.Post("/register", handler.Register)
	router.Post("/login", handler.Login)
	router.Post("/task", handler.CreateTask)
	router.Get("/status/{task_id}", handler.GetStatusTask)
	router.Get("/result/{task_id}", handler.GetResultTask)
	router.Handle("/metrics", promhttp.Handler())
	return &Router{router}
}
