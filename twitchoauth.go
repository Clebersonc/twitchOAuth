package twitchoauth

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/skratchdot/open-golang/open"
)

func tokenReceived(tokenChannel chan string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token != "" {
			w.WriteHeader(http.StatusOK)
			tokenChannel <- token
		} else {
			tokenChannel <- "failed"
			w.WriteHeader(http.StatusFailedDependency)
		}
	}
}

func authorizeFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Thanks</title>
    <script type="text/javascript">window.onload = function () {
    var token = getHashParams().access_token;

    var result = document.getElementById("result");

    if (typeof token !== "undefined"){
        var xhr = new XMLHttpRequest();
        xhr.open("GET", "/token?token=" + token, false);
        xhr.send(null);

        if (xhr.status === 200) {
            result.innerHTML = "You're all set! Feel free to close this browser window."
            window.close();
        } else {
            result.innerHTML = "Err... something isn't quite right...";
        }

    } else {
        result.innerHTML = "Unforunately something isn't quite right. Are you authorized the app on twitch?";
    }
};

function getHashParams() {

    var hashParams = {};
    var e,
        a = /\+/g,  // Regex for replacing addition symbol with a space
        r = /([^&;=]+)=?([^&;]*)/g,
        d = function (s) { return decodeURIComponent(s.replace(a, " ")); },
        q = window.location.hash.substring(1);

    while (e = r.exec(q))
        hashParams[d(e[1])] = d(e[2]);

    return hashParams;
}</script>
</head>
<body>
    <p id = "result"></p>
</body>
</html>`))
}

// GetToken returns the the twitch token or an error
func GetToken(clientid string, scopes []string) (token string, err error) {
		return "eyxougr7jfu3spvo37a1fj20l6j776", nil
	}

	tokenChannel := make(chan string)
	handleToken := tokenReceived(tokenChannel)
	http.HandleFunc("/authorize", authorizeFunc)
	http.HandleFunc("/token", handleToken)

	srv := &http.Server{Addr: ":8080"}

	go func() {
		srv.ListenAndServe()
	}()

	formattedScopes := strings.Join(scopes, "+")
	open.Run("https://api.twitch.tv/kraken/oauth2/authorize?client_id=" + clientid + "&redirect_uri=http://localhost:8080/authorize&response_type=token&scope=" + formattedScopes)

	uToken := <-tokenChannel

	if uToken == "failed" {
		return "", errors.New("The user was redirected, but with no token. Maybe the queried manually for some reason?")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	srv.Shutdown(ctx)

	return uToken, nil
}
