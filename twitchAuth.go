package twitchAuth

import (
	"net/http"
	"go/build"
	"path/filepath"
	"github.com/skratchdot/open-golang/open"
	"time"
	"context"
	"strings"
	"errors"
	"os"
	"io/ioutil"
)

func tokenReceived(tokenChannel chan string) func(http.ResponseWriter, *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
		token := r.URL.Query().Get("token")
		if token != ""{
			//log.Println("token get!")
			w.WriteHeader(http.StatusOK)
			tokenChannel <- token
		} else {
			tokenChannel <- "failed"
			w.WriteHeader(http.StatusFailedDependency)
		}
	}
}

func GetToken(clientid string, scopes []string)(token string, err error){
	importPath := "github.com/simplyserenity/twitchOAuth"

	p, err := build.Default.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}

	fs := http.FileServer(http.Dir(filepath.Join(p.Dir, "static")))

	confFile, oErr := os.OpenFile(p.Dir + "/config.dat", os.O_RDWR|os.O_CREATE, os.ModePerm)

	if oErr != nil {
		return "", oErr
	}

	content, fErr := ioutil.ReadFile(p.Dir + "/config.dat")

	if fErr != nil {
		return "", fErr
	} else if string(content) != "" {
		return string(content), nil
	}

	tokenChannel := make(chan string)
	handleToken := tokenReceived(tokenChannel)
	http.Handle("/", fs)
	http.HandleFunc("/token", handleToken)
	//log.Println("User sent to auth page.")

	srv := &http.Server{Addr: ":8080"}

	go func() {
		srv.ListenAndServe()
	}()

	//log.Println("Server started!")
	formattedScopes := strings.Join(scopes, "+")
	open.Run("https://api.twitch.tv/kraken/oauth2/authorize?client_id=" + clientid + "&redirect_uri=http://localhost:8080/authorize.html&response_type=token&scope="+formattedScopes)

	uToken := <- tokenChannel

	if uToken == "failed" {
		return "", errors.New("The user was redirected, but with no token. Maybe the queried manually for some reason?")
	}

	ctx, _ := context.WithTimeout(context.Background(), 1 * time.Second)

	srv.Shutdown(ctx)

	confFile.WriteString(uToken)

	confFile.Sync()

	if cErr := confFile.Close(); cErr != nil {
		return "", cErr
	}

	return uToken, nil
}