package elasticsearch

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func CheckElasticsearch(endpoint, username, password string, timeout int, skipTLS bool) (int, string) {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	if skipTLS {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	req, err := http.NewRequest("GET", endpoint+"/_cluster/health", nil)
	if err != nil {
		return 3, fmt.Sprintf("UNKNOWN - Errore creazione richiesta: %v", err)
	}

	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		return 3, fmt.Sprintf("UNKNOWN - Connessione fallita: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 2, fmt.Sprintf("CRITICAL - Elasticsearch risponde con %d", resp.StatusCode)
	}

	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return 3, fmt.Sprintf("UNKNOWN - Errore parsing JSON: %v", err)
	}

	switch health.Status {
	case "green":
		//log.Printf("health  %s", health.Status)
		return 0, "OK - Cluster Elasticsearch sano (green)"
	case "yellow":
		return 1, "WARNING - Cluster in stato yellow"
	case "red":
		return 2, "CRITICAL - Cluster in stato red"
	default:
		return 3, fmt.Sprintf("UNKNOWN - Stato non riconosciuto: %s", health.Status)
	}
}
