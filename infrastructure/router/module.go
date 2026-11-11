package router

import "github.com/gorilla/mux"

// Module is a self-contained set of HTTP endpoints for one feature of the
// application (auth, courses, ...). It owns every route and handler related
// to that feature.
//
// The application bootstrap collects all modules and calls Setup on each one
// during startup. All route registration must happen inside Setup.
type Module interface {
	// Path returns the URL prefix of this module (for example "/auth").
	//
	// An empty string means the module adds no prefix. Otherwise the value
	// must start with "/" and must not end with one. It must not change
	// during the lifetime of the module.
	Path() string

	// Setup registers the module's endpoints.
	//
	// root is the application's top-level router, for routes that must live
	// outside this module's prefix. Most implementations ignore it.
	//
	// router is the router the module registers its routes on. Setup is
	// given the authenticated API router for protected modules; its middleware
	// supplies the persisted user through session.CurrentUser(r.Context()).
	// Public authentication routes use a separate router.
	// Setup is
	// responsible for applying Path() to it, typically with
	// router.PathPrefix(m.Path()).Subrouter(), so route patterns are relative
	// to that prefix.
	//
	// Setup is called once during startup, before the server accepts requests.
	Setup(root *mux.Router, router *mux.Router)
}
