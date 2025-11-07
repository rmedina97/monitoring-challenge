package kibana

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type KibanaStatus struct {
	Status struct {
		Overall struct {
			Level string `json:"level"`
		} `json:"overall"`
	} `json:"status"`
}

func CheckKibana(endpoint, username, password string, timeout int, skipTLS bool) (int, string) {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	if skipTLS {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	req, err := http.NewRequest("GET", endpoint+"/api/status", nil)
	if err != nil {
		return 3, fmt.Sprintf("UNKNOWN - Errore creazione richiesta: %v", err)
	}

	//da verificare, api/status è pubblica e non richiede auth, in caso di auth sbagliato lo ignora
	req.SetBasicAuth(username, password)

	//ottimizza la ricezione del codice, 503 puo essere warning o critical, 200 può esser eok o warning se degradato
	resp, err := client.Do(req)
	if err != nil {
		return 3, fmt.Sprintf("UNKNOWN - Connessione fallita: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 2, fmt.Sprintf("CRITICAL - Kibana risponde con %d", resp.StatusCode)
	}

	var status KibanaStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return 3, fmt.Sprintf("UNKNOWN - Errore parsing JSON: %v", err)
	}

	switch status.Status.Overall.Level {
	case "available":
		return 0, "OK - Kibana è disponibile"
	case "degraded":
		return 1, "WARNING - Kibana è degradato"
	case "unavailable":
		return 2, "CRITICAL - Kibana non disponibile"
	default:
		return 3, fmt.Sprintf("UNKNOWN - Stato non riconosciuto: %s", status.Status.Overall.Level)
	}
}
