package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"math/rand"
	"time"
	"log"
	"github.com/fatih/color"
)

const status = "dnd" 

var tokens = []string{
	"Token1", // Token 1
	"Token2", // Ajoute d'autres tokens ici
}

var statusText = []string{
	"👋 Open for any Query",
	"🧑🏻‍🎓 MoneySlayer",
	
}
var currentStatusIndex = 0

type customStatus struct {
	Text string `json:"text"`
}

type statusData struct {
	Status      string      `json:"status"`
	CustomStatus customStatus `json:"custom_status"`
}

func info(message string, withLoading bool) {
	gray := color.New(color.FgHiBlack).SprintFunc()
	purple := color.New(color.FgMagenta).SprintFunc()

	currentTime := time.Now().Format("15:04:05")
	if withLoading {
		fmt.Printf("%s %sINF %s> %s%s", gray(currentTime), purple("INF"), gray(">"), message, " ")
		for i := 0; i < 3; i++ {
			time.Sleep(200 * time.Millisecond)
			fmt.Print(".")
		}
		fmt.Println()
	} else {
		fmt.Printf("%s %sINF %s> %s\n", gray(currentTime), purple("INF"), gray(">"), message)
	}
}

func changeStatus(token, statusText string) {
	client := &http.Client{}
	req, err := json.Marshal(statusData{
		Status: status,
		CustomStatus: customStatus{
			Text: statusText,
		},
	})

	if err != nil {
		log.Fatal("Erreur lors de la création de la requête JSON:", err)
	}

	reqBody := bytes.NewBuffer(req)
	request, err := http.NewRequest("PATCH", "https://discord.com/api/v8/users/@me/settings", reqBody)
	if err != nil {
		log.Fatal("Erreur lors de la création de la requête HTTP:", err)
	}
	request.Header.Set("Authorization", token)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Println("Erreur lors de la requête HTTP:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		info(fmt.Sprintf("Statut changé avec succès en '%s' -> %d OK", statusText, response.StatusCode), false)
	} else if response.StatusCode == 429 {
		body, _ := ioutil.ReadAll(response.Body)
		var data map[string]interface{}
		json.Unmarshal(body, &data)
		retryAfter := data["retry_after"].(float64)
		info(fmt.Sprintf("Ratelimitée, attente de %.2f secondes", retryAfter/1000), false)
		time.Sleep(time.Duration(retryAfter) * time.Millisecond)
	} else {
		info(fmt.Sprintf("Échec du changement de statut -> %d", response.StatusCode), false)
	}
}

func main() {
	for {
		for _, token := range tokens {
			go changeStatus(token, statusText[currentStatusIndex])
		}

		currentStatusIndex = (currentStatusIndex + 1) % len(statusText)

		
		randomDelay := 30 + rand.Intn(60)
		time.Sleep(time.Duration(randomDelay) * time.Second)
	}
}
