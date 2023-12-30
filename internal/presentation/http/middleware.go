package http

import (
	"context"
	"net/http"
)

func (p *presentation) GetUsernameFromSession(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sessionCookie, err := r.Cookie("session") // TODO(thomasmarlow): const or defines package (find a better name)
		if err != nil {
			wrappedHandler.ServeHTTP(w, r)
			return
		}

		username, err := p.session.GetUsername(sessionCookie.Value)
		if err == nil {
			ctx := context.WithValue(r.Context(), "username", username) // TODO(thomasmarlow): "should not use built-in type string as key for value; define your own type to avoid collisions"
			wrappedHandler.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		switch err.(type) {
		default:
			wrappedHandler.ServeHTTP(w, r)
		}
	})
}
