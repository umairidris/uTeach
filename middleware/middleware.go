// Package middleware provides app specific middleware handlers.
package middleware

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/context"
	"github.com/umairidris/uTeach/httperror"
	"github.com/umairidris/uTeach/models"
	"github.com/umairidris/uTeach/session"
)

type Middleware struct {
	App *application.App
}

func (m *Middleware) SetThreadIDVar(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		threadID, err := strconv.ParseInt(vars["threadID"], 10, 64)
		if err != nil {
			httperror.HandleError(w, httperror.StatusError{http.StatusBadRequest, err})
			return
		}

		context.SetThreadID(r, threadID)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// MustLogin ensures that the next handler is only accessible by users that are logged in.
func (m *Middleware) MustLogin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		usm := session.NewUserSessionManager(m.App.CookieStore)
		if _, ok := usm.SessionUser(r); !ok {
			httperror.HandleError(w, httperror.StatusError{http.StatusForbidden, nil})
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (m *Middleware) isAdmin(r *http.Request) bool {
	usm := session.NewUserSessionManager(m.App.CookieStore)
	user, _ := usm.SessionUser(r)
	return user.IsAdmin
}

func (m *Middleware) isThreadCreator(r *http.Request) bool {
	tm := models.NewThreadModel(m.App.DB)
	threadID := context.ThreadID(r)
	thread, err := tm.GetThreadByID(nil, threadID)
	if err != nil {
		return false
	}
	usm := session.NewUserSessionManager(m.App.CookieStore)
	user, _ := usm.SessionUser(r)
	return thread.Creator == user
}

func (m *Middleware) MustBeAdmin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !m.isAdmin(r) {
			httperror.HandleError(w, httperror.StatusError{http.StatusForbidden, nil})
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (m *Middleware) MustBeAdminOrThreadCreator(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !m.isThreadCreator(r) && !m.isAdmin(r) {
			httperror.HandleError(w, httperror.StatusError{http.StatusForbidden, nil})
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
