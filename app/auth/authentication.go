package auth

import (
	"github.com/gorilla/sessions"
	"math/rand"
	"net/http"
)

var cookieStore *sessions.CookieStore

var Authenticated string = "authenticated"
var UserID string = "user_id"

var session_cookie string = "vb3d_v0_middleware"

func init() {
	key := make([]byte, 32)
	rand.Read(key)
	cookieStore = sessions.NewCookieStore(key)
}

func GetCookieStore() *sessions.CookieStore {
	return cookieStore
}

func GetSessionCookie() string {
	return session_cookie
}

func Auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		session, _ := cookieStore.Get(r, session_cookie)
		if auth, ok := session.Values[Authenticated].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			h(w, r)
		}
	}
}
