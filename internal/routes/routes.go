package routes

import (
	"net/http"

	"app/internal/ui"
	"app/internal/views/pages"
	"app/internal/views/components/icons"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

func Register(router *chi.Mux) {
	var shell ui.ShellFunction = pages.App
	var wrapper ui.WrapperTemplate = pages.Index

	group := ui.RouteGroup{
		Title:    "App",
		Extras: []ui.NavItem{
			{ID: "docs", Label: "Documents", Icon: icons.Document("size-4"), OnMenubar: true, OnSidebar: true, Order: 1},
		},
		Routes: []ui.Route{
			{
				Method: "GET",
				Path:   "/",
				NavConfig: &ui.NavConfig{
					Label: "Home",
					Icon: icons.Home("size-4"),
					OnSidebar: true,
					OnMenubar: true,
					Order: 0,
				},
				Page: func(_ *http.Request) (templ.Component, error) {
					return pages.Home(), nil
				},
			},
			{
				Method: "GET",
				Path:   "/documents",
				NavConfig: &ui.NavConfig{
					Label: "Documents",
					Icon: icons.Document("size-4"),
					ParentID: "docs",
					ID: "DocsFirstChild",
					OnSidebar: true,
					OnMenubar: true,
					Order: 0,
				},
				Page: func(_ *http.Request) (templ.Component, error) {
					return pages.Documents(), nil
				},
			},
			{
				Method: "GET",
				Path:   "/example",
				NavConfig: &ui.NavConfig{
					Label: "Exampe1",
					Icon: icons.Document("size-4"),
					ParentID: "docs",
					OnSidebar: true,
					OnMenubar: true,
					Order: 1,
				},
				Page: func(_ *http.Request) (templ.Component, error) {
					return pages.Documents(), nil
				},
			},
			// {
			// 	Method: "GET",
			// 	Path: "/documents/{id}",
			// 	Page: func(req *http.Request) (templ.Component, error) {
			// 		id := chi.URLParam(req, "id")
			// 		return pages.Documents(id), nil
			// 	},
			// },
		},
	}

	router.Group(func(g chi.Router) {
		group.RegisterPage(g, shell, wrapper, pages.Error("Something went wrong"))
	})

}
