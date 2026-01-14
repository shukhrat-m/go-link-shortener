package main

import (
	"go/adv-demo/configs"
	"go/adv-demo/internal/auth"
	"go/adv-demo/internal/link"
	"go/adv-demo/internal/stat"
	"go/adv-demo/internal/user"
	"go/adv-demo/pkg/db"
	"go/adv-demo/pkg/event"
	"go/adv-demo/pkg/middleware"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	eventBus := event.NewEventBus()

	router := http.NewServeMux()

	//Repositories
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	statRepository := stat.NewStatRepository(db)

	//Services
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{EventBus: eventBus, StatRepository: statRepository})

	// Initialize Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{Config: conf, AuthService: authService})
	link.NewLinkHandler(router, link.LinkHandlerDeps{LinkRepository: linkRepository, EventBus: eventBus, Config: conf})
	stat.NewLinkHandler(router, stat.StatHandlerDeps{StatRepository: statRepository, Config: conf, EventBus: eventBus})

	go statService.AddClick()

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	server := &http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	server.ListenAndServe()
}
