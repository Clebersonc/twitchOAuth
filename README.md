# twitchOAuth
It's a little hacky so far, so there is certainly room for improvements, but it works!


## Usage


It's fairly simple to use in the current state all you really need is your client id


```go

package main

import "github.com/simplyserenity/twitchOAuth"

func main(){
  token := twitchAuth.GetToken(<clientID>);
}

```


And that's it really, the user's browser will open and they'll be asked to verify, and the token should make its way back to you.
