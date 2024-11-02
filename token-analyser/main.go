package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

const discordAPI = "https://discord.com/api/v9"

type UserInfo struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Verified      bool   `json:"verified"`
	MFAEnabled    bool   `json:"mfa_enabled"`
	PremiumType   int    `json:"premium_type"`
	Locale        string `json:"locale"`
	Flags         int    `json:"flags"`
	Bio           string `json:"bio"`
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
		c := color.New(color.FgCyan)
		c.Println("Erreur lors de la création de la requête :", err)
		return
	}
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		c := color.New(color.FgCyan)
		c.Println("Erreur lors de l'envoi de la requête :", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		c := color.New(color.FgCyan)
		c.Println("Token invalide :", token)
		c.Printf("Statut de la réponse : %d, Corps : %s\n", resp.StatusCode, body)
		return
	}

	var userInfo UserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		c := color.New(color.FgCyan)
		c.Println("Erreur lors de l'analyse de la réponse JSON :", err)
		return
	}
	c := color.New(color.FgCyan)
	c.Printf("Informations utilisateur :\n- ID : %s\n- Nom : %s#%s\n- Email : %s\n- Téléphone : %s\n- Vérifié : %t\n- MFA activé : %t\n- Type de Nitro : %d\n- Locale : %s\n- Flags : %d\n- Bio : %s\n",
		userInfo.ID, userInfo.Username, userInfo.Discriminator, userInfo.Email, userInfo.Phone, userInfo.Verified, userInfo.MFAEnabled, userInfo.PremiumType, userInfo.Locale, userInfo.Flags, userInfo.Bio)

	countFriends(token)
}

func countFriends(token string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", discordAPI+"/users/@me/relationships", nil)
	if err != nil {
		c := color.New(color.FgCyan)
		c.Println("Erreur lors de la création de la requête pour les amis :", err)
		return
	}
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		c := color.New(color.FgCyan)
		c.Println("Erreur lors de l'envoi de la requête pour les amis :", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		c := color.New(color.FgCyan)
		c.Println("Impossible de récupérer la liste des amis, statut :", resp.StatusCode)
		c.Printf("Corps de la réponse : %s\n", body)
		return
	}

	var friends []Friend
	if err := json.Unmarshal(body, &friends); err != nil {
		c := color.New(color.FgCyan)
		c.Println("Erreur lors de l'analyse de la liste des amis :", err)
		return
	}

	totalFriends := len(friends)
	onlineFriends := 0

	for _, friend := range friends {
		if friend.Presence.Status == "online" {
			onlineFriends++
		}
	}
	d := color.New(color.FgCyan, color.Bold)

	d.Printf("Total d'amis : %d\n", totalFriends)
}

func readTokensAndCheck(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		color.New(color.FgRed).Println("Erreur lors de l'ouverture du fichier :", err)
		return
	}
	defer file.Close()

	var token string
	for {
		_, err := fmt.Fscanln(file, &token)
		if err != nil {
			break
		}
		checkToken(token)
		c := color.New(color.FgRed)
		c.Println("-----------------------------------------------------------")

		waitTime := rand.Intn(8) + 5
		time.Sleep(time.Duration(waitTime) * time.Second)
	}
}

func main() {
	if len(os.Args) < 2 {
		c := color.New(color.FgCyan)
		c.Println("Usage: go run main.go <fichier des tokens>")
		return
	}

	rand.Seed(time.Now().UnixNano())
	readTokensAndCheck(os.Args[1])
}
