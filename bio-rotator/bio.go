package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var token = "token[usr]"
var bioTexts = []string{"ta bio1", "ta bio2", "ta bio3"}
var currentBioIndex = 0

func updateBio() {
	url := "https://discord.com/api/v8/users/@me/profile"

	data := map[string]string{
		"bio": bioTexts[currentBioIndex],
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Erreur lors de la conversion des données en JSON :", err)
		return
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Erreur lors de la création de la requête :", err)
		return
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de la requête :", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Biographie changée avec succès en '%s' -> %d OK\n", bioTexts[currentBioIndex], resp.StatusCode)
	} else if resp.StatusCode == 429 {
		body, _ := ioutil.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		retryAfter := result["retry_after"].(float64) / 1000
		fmt.Printf("Ratelimitée, attente de %.2f secondes\n", retryAfter)
		time.Sleep(time.Duration(retryAfter) * time.Second)
	} else {
		fmt.Printf("Échec du changement de biographie -> %d\n", resp.StatusCode)
	}

	currentBioIndex = (currentBioIndex + 1) % len(bioTexts)
}

func main() {
	for {
		updateBio()
		time.Sleep(30 * time.Second)
	}
}
