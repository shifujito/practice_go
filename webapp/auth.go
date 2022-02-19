package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

func GetConnect() *oauth2.Config {
	loadEnv()
	config := &oauth2.Config{
		ClientID:     os.Getenv("ClientID"),
		ClientSecret: os.Getenv("ClientSecret"),
		// ここにはなにをいれるのか？
		Endpoint:    google2.Endpoint,
		RedirectURL: "http://localhost:8081/auth/callback/google",
		Scopes:      []string{"profile"},
	}
	return config
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	// provider := segs[3]
	switch action {
	case "login":
		config := GetConnect()
		// ここはなにをいれるべきなのか
		url := config.AuthCodeURL("SECURITY_KEY", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
		// url := config.AuthCodeURL("SECURITY_KEY")
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		// oauthコネクト
		config := GetConnect()
		// 認可コード
		code := r.URL.Query().Get("code")
		// contextはよくわからないから調べる
		ctx := context.Background()
		// アクセストークン
		tok, err := config.Exchange(ctx, code)
		if err != nil {
			panic(err)
		}
		// HTTPクライアント
		client := config.Client(ctx, tok)
		// get user info
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(body))
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sには非対応", action)
	}
}
