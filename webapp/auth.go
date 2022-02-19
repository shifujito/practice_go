package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	google2 "golang.org/x/oauth2/google"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// 未認証
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(err.Error())
	} else {
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	// provider := segs[3]
	switch action {
	case "login":
		config := oauth2.Config{
			ClientID:     "880683417554-lo76iaklpb2vif4j96opbs3fcphat66n.apps.googleusercontent.com",
			ClientSecret: "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",
			Endpoint:     google2.Endpoint,
			RedirectURL:  "http://localhost:8081/auth/callback/google",
			Scopes:       []string{"openid", "email", "profile"},
		}
		url := config.AuthCodeURL("SECURITY_KEY", oauth2.AccessTypeOffline, oauth2.ApprovalForce)

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sには非対応", action)
	}
}
