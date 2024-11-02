package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const discordAPI = "https://discord.com/api/v9"

type UserInfo struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	PremiumType int    `json:"premium_type"`
}

type Friend struct {
	User struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	} `json:"user"`
	Presence struct {
		Status string `json:"status"`
	} `json:"presence"`
}

func checkToken(token string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", discordAPI+"/users/@me", nil)
	if err != nil {
		fmt.Println("Erreur lors de la création de la requête :", err)
		return
	}
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de la requête :", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Token invalide :", token)
		return
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		fmt.Println("Erreur lors de l'analyse de la réponse JSON :", err)
		return
	}

	fmt.Printf("Token valide pour l'utilisateur : %s#%s\n", userInfo.Username, userInfo.ID)

	switch userInfo.PremiumType {
	case 1:
		fmt.Println("L'utilisateur possède un Nitro Basic.")
	case 2:
		fmt.Println("L'utilisateur possède un Nitro Boost.")
	default:
		fmt.Println("L'utilisateur n'a pas de Nitro.")
	}

	countFriends(token)
}

func countFriends(token string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", discordAPI+"/users/@me/relationships", nil)
	if err != nil {
		fmt.Println("Erreur lors de la création de la requête pour les amis :", err)
		return
	}
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de la requête pour les amis :", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Impossible de récupérer la liste des amis, statut :", resp.StatusCode)
		return
	}

	var friends []Friend
	if err := json.NewDecoder(resp.Body).Decode(&friends); err != nil {
		fmt.Println("Erreur lors de l'analyse de la liste des amis :", err)
		return
	}

	totalFriends := len(friends)
	onlineFriends := 0

	for _, friend := range friends {
		if friend.Presence.Status == "online" {
			onlineFriends++
		}
	}

	fmt.Printf("Total d'amis : %d\n", totalFriends)
}

func checkTokensFromFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier :", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		token := scanner.Text()
		checkToken(token)

		fmt.Println("-----------------------------------------------------------")

		delay := time.Duration(5+rand.Intn(6)) * time.Second
		time.Sleep(delay)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erreur lors de la lecture du fichier :", err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <fichier_tokens>")
		return
	}

	filename := os.Args[1]
	checkTokensFromFile(filename)
}
