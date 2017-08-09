window.onload = function () {
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
}