package ui

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

type ShellFunction func(body templ.Component) templ.Component

type WrapperTemplate func (
	title string,
	shell templ.Component,
	navItems []NavItem,
	currentPath string,
) templ.Component

func (routeGroup RouteGroup) RegisterPage(router chi.Router, ShellFunction ShellFunction, wrapperTemplate WrapperTemplate, errorTemplate templ.Component) {
	for _, route := range routeGroup.Routes {
		var wrappedHandler http.Handler = http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
			pageBody, pageErr := route.Page(request)
			if pageErr != nil {
				templ.Handler(errorTemplate).ServeHTTP(responseWriter, request)
				return
			}
			page := wrapperTemplate(
				routeGroup.Title,
				ShellFunction(pageBody),
				routeGroup.GroupRoutes(request),
				request.URL.Path,
			)
			templ.Handler(page).ServeHTTP(responseWriter, request)
		})

		for i := len(route.Middlewares) - 1; i >= 0; i-- {
			wrappedHandler = route.Middlewares[i](wrappedHandler)
		}

		router.Method(route.Method, route.Path, wrappedHandler)
	}
}
