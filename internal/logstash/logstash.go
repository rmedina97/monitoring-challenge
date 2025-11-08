package logstash

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// HealthResponse rappresenta la struttura minima della risposta _health_report
type HealthResponse struct {
	Status string `json:"status"`
}

func CheckLogstash(endpoint, username, password string, timeout int, skipTLS bool) (int, string) {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	if skipTLS {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	req, err := http.NewRequest("GET", endpoint+"/_health_report", nil)
	if err != nil {
		return 3, fmt.Sprintf("UNKNOWN - Errore creazione richiesta: %v", err)
	}

	// Set Basic Auth se username e password forniti
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 3, fmt.Sprintf("UNKNOWN - Connessione fallita: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return 3, "UNKNOWN - Unauthorized (credenziali mancanti o errate)"
	}

	if resp.StatusCode != 200 {
		return 2, fmt.Sprintf("CRITICAL - Logstash risponde con HTTP %d", resp.StatusCode)
	}
	log.Printf("response is %d", resp.StatusCode)

	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return 3, fmt.Sprintf("UNKNOWN - Errore parsing JSON: %v", err)
	}

	switch health.Status {
	case "green":
		return 0, "OK - Logstash sano (green)"
	case "yellow":
		return 1, "WARNING - Logstash parzialmente degradato (yellow)"
	case "red":
		return 2, "CRITICAL - Logstash in stato red"
	default:
		return 3, fmt.Sprintf("UNKNOWN - Stato non riconosciuto: %s", health.Status)
	}
}
