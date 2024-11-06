package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	apiURL    = "https://discord.com/api/v10/users/@me/settings"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
)


func setStatus(token, status string) error {
	
	payload := map[string]string{"status": status}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Erreur de création du payload JSON : %v", err)
	}
	req, err := http.NewRequest("PATCH", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Erreur de création de la requête : %v", err)
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)


	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Erreur de requête HTTP : %v", err)
	}
	defer resp.Body.Close()

	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Échec de mise à jour du statut, code HTTP: %d", resp.StatusCode)
	}
	return nil
}


func loadTokens(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de l'ouverture du fichier %s : %v", filename, err)
	}
	defer file.Close()

	var tokens []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tokens = append(tokens, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Erreur lors de la lecture du fichier %s : %v", filename, err)
	}

	return tokens, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./status <online|dnd|idle|invisible>")
		return
	}

	
	status := os.Args[1]
	validStatuses := map[string]bool{"online": true, "dnd": true, "idle": true, "invisible": true}

	if !validStatuses[status] {
		fmt.Println("Statut invalide. Utilisez: online, dnd, idle, ou invisible.")
		return
	}

	tokens, err := loadTokens("token.txt")
	if err != nil {
		fmt.Println("Erreur de chargement des tokens:", err)
		return
	}

	for {
		
		for _, token := range tokens {
			err := setStatus(token, status)
			if err != nil {
				fmt.Printf("Erreur lors de la mise à jour du statut pour le token %s : %v\n", token, err)
			} else {
				fmt.Printf("Statut mis à jour avec succès pour le token %s sur '%s'\n", token, status)
			}
		}

	
		time.Sleep(5 * time.Minute)
	}
}
