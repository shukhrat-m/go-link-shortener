package link

import (
	"go/adv-demo/configs"
	"go/adv-demo/pkg/event"
	"go/adv-demo/pkg/middleware"
	"go/adv-demo/pkg/request"
	"go/adv-demo/pkg/response"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type LinkHandlerDeps struct {
	LinkRepository *LinkRepository
	*configs.Config
	EventBus *event.EventBus
}

type LinkHandler struct {
	LinkRepository *LinkRepository
	EventBus       *event.EventBus
}

func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{LinkRepository: deps.LinkRepository, EventBus: deps.EventBus}

	router.HandleFunc("POST /link", handler.Create())
	router.Handle("PATCH /link/{id}", middleware.IsAuthed(handler.Update(), deps.Config))
	router.HandleFunc("DELETE /link/{id}", handler.Delete())
	router.HandleFunc("GET /{hash}", handler.GoTo())
	router.Handle("GET /link", middleware.IsAuthed(handler.GetAll(), deps.Config))
}

func (handler *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[LinkCreateRequest](&w, r)

		if err != nil {
			return
		}

		defer r.Body.Close()
		link := NewLink(body.Url)

		var existingLink *Link

		for {
			existingLink, err = handler.LinkRepository.GetByHash(link.Hash)

			if err != nil || existingLink == nil {
				break
			}

			link.Hash = RandStringRunes(6)
		}

		createdLink, err := handler.LinkRepository.Create(link)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response.JSON(w, 201, createdLink)
	}
}

func (handler *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[LinkUpdateRequest](&w, r)

		if err != nil {
			return
		}

		defer r.Body.Close()

		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 64)

		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		link, err := handler.LinkRepository.Update(&Link{
			Model: gorm.Model{ID: uint(id)},
			Url:   body.Url,
			Hash:  body.Hash,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response.JSON(w, 200, link)
	}
}

func (handler *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")

		id, err := strconv.ParseUint(idString, 10, 64)

		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		_, err = handler.LinkRepository.GetById(uint(id))

		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		err = handler.LinkRepository.Delete(uint(id))

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (handler *LinkHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		link, err := handler.LinkRepository.GetByHash(hash)

		if err != nil {
			http.Error(w, "Link not found", http.StatusNotFound)
			return
		}

		handler.EventBus.Publish(event.Event{Type: event.LinkVisitedEvent, Data: link.ID})
		http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect)
	}
}

func (handler *LinkHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}

		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))

		if err != nil {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}

		links := handler.LinkRepository.GetAll(uint(limit), uint(offset))
		count := handler.LinkRepository.GetCount()

		response.JSON(w, 200, &GetAllLinksResponse{Links: links, Count: count})
	}
}
