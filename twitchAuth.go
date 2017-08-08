package twitchAuth

import (
	"log"
	"net/http"
	"go/build"
	"path/filepath"
	"github.com/skratchdot/open-golang/open"
	"time"
)

func tokenReceived(tokenChannel chan string) func(http.ResponseWriter, *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
		token := r.URL.Query().Get("token")
		if token != ""{
			log.Println("token get!")
			w.WriteHeader(http.StatusOK)
			tokenChannel <- token
		} else {
			w.WriteHeader(http.StatusFailedDependency)
		}
	}
}

func GetToken(clientid string)(token string){
	importPath := "github.com/simplyserenity/twitchOAuth"

	p, err := build.Default.Import(importPath, "", build.FindOnly)
	if err != nil {
		panic(err)
	}

	fs := http.FileServer(http.Dir(filepath.Join(p.Dir, "static")))

	tokenChannel := make(chan string)
	handleToken := tokenReceived(tokenChannel)
	http.Handle("/", fs)
	http.HandleFunc("/token", handleToken)
	log.Println("User sent to auth page.")

	srv := &http.Server{Addr: ":8080"}

	go func() {
		srv.ListenAndServe()
	}()

	log.Println("Server started!")

	open.Run("https://api.twitch.tv/kraken/oauth2/authorize?client_id=" + clientid + "&redirect_uri=http://localhost:8080/authorize.html&response_type=token&scope=chat_login+user_read")

	uToken := <- tokenChannel

	time.Sleep(100 * time.Millisecond)

	srv.Close()		//was going to use shutdown but it would throw nil pointer errors every time

	return uToken
}