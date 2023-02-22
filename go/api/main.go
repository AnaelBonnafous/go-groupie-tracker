// package main

// import (
// 	"encoding/base64"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/joho/godotenv"
// )

// func main() {
// 	auth_string := goDotEnvVariable("CLIENT_ID") + ":" + goDotEnvVariable("CLIENT_SECRET")
// 	auth_byte := []byte(auth_string)
// 	auth_base64 := base64.StdEncoding.EncodeToString(auth_byte)
// 	urll := "https://accounts.spotify.com/api/token"
// 	req, _ := http.NewRequest("POST", urll, nil)
// 	client := &http.Client{}
// 	req.Header = http.Header{
// 		"Authorization": {"Basic " + auth_base64},
// 		"Content-Type":  {"application/x-www-form-urlencoded"},
// 	}
// 	res, _ := client.Do(req)
// 	fmt.Println(res)

// }
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

type Response struct {
	Artists struct {
		Href  string `json:"href"`
		Items []struct {
			External_urls struct {
				Spotify string
			}
		}
	}
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func main() {
	authConfig := &clientcredentials.Config{
		ClientID:     goDotEnvVariable("CLIENT_ID"),
		ClientSecret: goDotEnvVariable("CLIENT_SECRET"),
		TokenURL:     spotify.TokenURL,
	}
	accessToken, err := authConfig.Token(context.Background())
	if err != nil {
		log.Fatalf("error retrieve access token: %v", err)
	}
	url := "https://api.spotify.com/v1/search"
	name := "damso"
	query := "?q=" + name + "&type=artist&limit=1"
	queryurl := url + query
	req, _ := http.NewRequest("GET", queryurl, nil)
	req.Header = http.Header{
		"Authorization": {"Bearer " + accessToken.AccessToken},
	}
	client := &http.Client{}
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	var responce Response
	json.Unmarshal(body, &responce)
	fmt.Println(responce)
}
