package handlers

import (
	"github.com/go-chi/chi/v5"
	_ "github.com/s-turchinskiy/EffectiveMobile/internal/handlers/swagger"
	"github.com/s-turchinskiy/EffectiveMobile/internal/middleware/logger"
	"github.com/s-turchinskiy/EffectiveMobile/internal/middleware/read_validate_json"
	httpswagger "github.com/swaggo/http-swagger"
)

// @Title Subscriptions API
// @Description Сервис хранения подписок.
// @Version 1.0

// @Contact.email s.turchinskiy@yandex.ru

// @BasePath /
// @Host nohost.io:8080

// @Tag.name Subscription
// @Tag.description "Взаимодействие с подпиской"

// @Tag.name Subscriptions
// @Tag.description "Отчеты по подпискам"

func Router(h *Handler) chi.Router {

	router := chi.NewRouter()
	router.Use(logger.Logger)
	//router.Use(gzip.GzipMiddleware)
	router.Use(readvalidatejson.Middleware)
	router.Route("/api", func(r chi.Router) {
		r.Route("/subscription", func(r chi.Router) {
			r.Post("/create", h.CreateSubscription)
			r.Get("/read", h.ReadSubscriptions)
			r.Put("/update", h.UpdateSubscription)
			r.Delete("/delete", h.DeleteSubscription)
		})

		r.Route("/subscriptions", func(r chi.Router) {
			r.Post("/sum", h.SumSubscriptions)
		})
	})

	router.Mount("/swagger", httpswagger.WrapHandler)

	return router

}
