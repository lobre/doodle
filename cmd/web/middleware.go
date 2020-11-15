package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/lobre/doodle/pkg/models"
)

// secureHeaders will inject headers in the response
// to prevent XSS and Clickjacking attacks.
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

// logRequest will use the information logger to log all requests.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

// recoverPanic will gracefully handle any panic that happens in the current go routine.
// By default, panics don't shut the entire application (only the current go routine),
// but if one arise, the server will return an empty response. This middleware is taking
// care of recovering the panic and sending a regular 500 server error.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// deferred function are always run in the event of a panic as
		// Go unwinds the stack.
		defer func() {
			if err := recover(); err != nil {
				// setting this header will make the http.Server
				// automatically close the current connection.
				w.Header().Set("Connection", "close")

				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// authenticate fetches the user's ID from their session data, checks the
// database to see if the ID is valid and for an active user, and then updates
// the request context to include this information.
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exists := app.session.Exists(r, "authenticatedUserID")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		user, err := app.userStore.Get(app.session.GetInt(r, "authenticatedUserID"))
		if errors.Is(err, models.ErrNoRecord) || !user.Active {
			// session exists but user has been removed or disabled from db
			app.session.Remove(r, "authenticatedUserID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyIsAuthenticated, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// requireAuthentication will redirect the user to the login page if they
// are not authenticated.
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Remove browser cache, as the targetted page is dynamic
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

// injectCSRFCookie will injects a customized CSRF token in a cookie (which is encrypted). That same token
// will be used as a hidden field in forms (from nosurf.Token()). On the form submission, the server
// will check that these two values match. It makes it impossible for an attacker
// to guess what the value of the token is. So directly trying to post a request to
// our secured endpoint without this parameter would fail.
// The only way to submit the form is from our frontend.
// This prevents an attacker from abusing the authenticated token, and have access
// to restricted resources when posting a request to our endpoint from a third party website.
func (app *application) injectCSRFCookie(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	cookie := http.Cookie{
		HttpOnly: true,
		Path:     "/",
	}

	if app.isHTTPS {
		cookie.Secure = true
	}

	csrfHandler.SetBaseCookie(cookie)
	return csrfHandler
}
