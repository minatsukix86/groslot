package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

const discordAPI = "https://discord.com/api/v9"

type NitroGift struct {
	ID           string `json:"id"`
	Code         string `json:"code"`
	Type         int    `json:"type"`
	Subscription string `json:"subscription_plan"`
	Duration     int    `json:"duration"`
	ExpiresAt    string `json:"expires_at"`
}

func checkNitroGift(token string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", discordAPI+"/gift-codes/"+token, nil)
	if err != nil {
		color.Red("Erreur lors de la création de la requête : %s", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		color.Red("Erreur lors de l'envoi de la requête : %s", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		color.Red("Le lien de cadeau Nitro est invalide ou a expiré. Statut : %d", resp.StatusCode)
		return
	}

	var nitroGift NitroGift
	if err := json.NewDecoder(resp.Body).Decode(&nitroGift); err != nil {
		color.Red("Erreur lors de l'analyse de la réponse JSON : %s", err)
		return
	}

	color.Green("Informations sur le cadeau Nitro :")
	color.Green("- ID : %s", nitroGift.ID)
	color.Green("- Code : %s", nitroGift.Code)
	color.Green("- Type : %s", nitroGift.Subscription)
	color.Green("- Durée : %d mois", nitroGift.Duration)
	color.Green("- Expiration : %s", nitroGift.ExpiresAt)

	expirationTime, err := time.Parse(time.RFC3339, nitroGift.ExpiresAt)
	if err != nil {
		color.Red("Erreur lors de la conversion de la date d'expiration : %s", err)
		return
	}
	timeRemaining := time.Until(expirationTime)

	color.Green("Temps restant avant expiration : %s", timeRemaining)
}

func main() {
	if len(os.Args) < 2 {
		color.Red("Usage : go run main.go <chemin_vers_le_fichier>")
		return
	}

	filePath := os.Args[1]

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		color.Red("Erreur lors de la lecture du fichier : %s", err)
		return
	}

	codes := string(data)
	for _, code := range splitLines(codes) {
		checkNitroGift(code)
		color.Green("-----------------------------------------------")

		time.Sleep(time.Duration(5+rand.Intn(8)) * time.Second)
	}
}

func splitLines(data string) []string {
	lines := strings.Split(strings.TrimSpace(data), "\n")
	return lines
}
