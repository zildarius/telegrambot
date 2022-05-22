package jokes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Joke struct {
	ID   int    `json:"id"`
	Joke string `json:"joke"`
}

type JokeResponse struct {
	Type  string `json:"type"`
	Value Joke   `json:"value"`
}

func ReturnNextJoke(jokeNumber string) string {
	httpClient := http.Client{}
	jokeURL := "http://api.icndb.com/jokes/random"
	if jokeNumber != "" {
		jokeURL = "http://api.icndb.com/jokes/" + string(jokeNumber)
	}

	log.Println("JOKE NUMBER: ", jokeNumber)
	log.Println("JOKE NUMBER2: ", string(jokeNumber))
	log.Println("JOKE URL: ", jokeURL)

	resp, err := httpClient.Get(jokeURL) ////?limitTo=[nerdy]")
	if err != nil {
		return "Jokes API not responding"
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	joke := JokeResponse{}

	log.Println("LOG: ", string(body))

	err = json.Unmarshal(body, &joke)
	if err != nil {
		return "Похоже шутки с таким номером не нашлось. Попробуйте ввести другой номер."
	}

	jokeString := strings.ReplaceAll(joke.Value.Joke, "&quot;", "\"")

	return strconv.Itoa(joke.Value.ID) + ": " + jokeString
}
