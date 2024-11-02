package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"
)

type FormData struct {
	Webhook  string
	Message  string
	Interval int
	ImageURL string
}

var (
	stopChan   chan bool
	mu         sync.Mutex
	isSpamming bool
)

func main() {
	stopChan = make(chan bool)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("form.html"))
		
		tmpl.Execute(w, map[string]bool{"isSpamming": isSpamming})
	})

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		r.ParseForm()
		data := FormData{
			Webhook:  r.FormValue("webhook"),
			Message:  r.FormValue("message"),
			Interval: parseInterval(r.FormValue("interval")),
			ImageURL: r.FormValue("imageURL"),
		}

		mu.Lock()
		if isSpamming {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			mu.Unlock()
			return 
		}

		stopChan = make(chan bool) 
		isSpamming = true         
		mu.Unlock()

		go func(data FormData) {
			embed := map[string]interface{}{
				"embeds": []map[string]interface{}{
					{
						"title":       "Nouveau Message",
						"description": data.Message,
						"color":       3447003, 
						"timestamp":   time.Now().Format(time.RFC3339),
						"footer": map[string]interface{}{
							"text": "Envoyé via Go",
						},
					},
				},
			}

			if data.ImageURL != "" {
				embed["embeds"].([]map[string]interface{})[0]["image"] = map[string]interface{}{
					"url": data.ImageURL,
				}
			}

			toBytes, err := json.Marshal(embed)
			if err != nil {
				fmt.Println("Erreur lors de la création de la charge JSON :", err)
				return
			}

			for {
				select {
				case <-stopChan:
					fmt.Println("Envoi arrêté.")
					mu.Lock()
					isSpamming = false 
					mu.Unlock()
					return
				default:
					resp, err := http.Post(data.Webhook, "application/json", bytes.NewBuffer(toBytes))
					if err != nil {
						fmt.Println("Erreur lors de l'envoi du message :", err)
						time.Sleep(time.Duration(data.Interval) * time.Millisecond)
						continue
					}
					resp.Body.Close()
					time.Sleep(time.Duration(data.Interval) * time.Millisecond)
				}
			}
		}(data)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		if isSpamming {
			close(stopChan)  
			isSpamming = false
		}
		mu.Unlock()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	fmt.Println("Le serveur est en cours d'exécution sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func parseInterval(intervalStr string) int {
	interval, err := time.ParseDuration(intervalStr + "ms")
	if err != nil {
		return 1000
	}
	return int(interval.Milliseconds())
}
