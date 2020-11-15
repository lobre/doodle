package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"github.com/lobre/doodle/pkg/embeds/staticdir"
)

func (app *application) routes() http.Handler {
	// This chain is used for every request our application receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// This chain is used for all routes that are not static (css, js, ...).
	dynamicMiddleware := alice.New(app.session.Enable, app.injectCSRFCookie, app.authenticate)

	mux := pat.New()

	mux.Get("/ping", http.HandlerFunc(ping))
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))

	fileServer := http.FileServer(staticdir.FS)
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// Events
	mux.Get("/event/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createEventForm))
	mux.Post("/event/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createEvent))
	mux.Get("/event/:id", dynamicMiddleware.ThenFunc(app.showEvent))

	// Users
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	return standardMiddleware.Then(mux)
}
