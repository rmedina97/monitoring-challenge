package main

import (
	"flag"
	"fmt"
	"os"

	elasticsearch "github.com/rmedina97/monitoring-challenge/internal/elasticsearch"
)

func main() {
	check := flag.String("check", "", "Servizio da controllare: elasticsearch|kibana|logstash, obbligatorio")
	host := flag.String("host", "", "Endpoint del servizio, obbligatorio")
	user := flag.String("user", "", "Username, obbligatorio")
	pass := flag.String("password", "", "Password, obbligatorio")
	timeout := flag.Int("timeout", 5, "Timeout in secondi")
	skipTLS := flag.Bool("skip-tls", true, "Ignora verifica TLS")

	flag.Parse()

	missing := false

	if *check == "" {
		fmt.Println("Errore: il parametro -check è obbligatorio.")
		missing = true
	}
	if *user == "" {
		fmt.Println("Errore: il parametro -user è obbligatorio.")
		missing = true
	}
	if *pass == "" {
		fmt.Println("Errore: il parametro -password è obbligatorio.")
		missing = true
	}
	if *host == "" {
		fmt.Println("Errore: il parametro -password è obbligatorio.")
		missing = true
	}

	if missing {
		//fmt.Println()
		flag.Usage() // mostra la descrizione di ogni flag (descrizione è il terzo parametro)
		os.Exit(1)
	}

	switch *check {
	case "elasticsearch":
		code, msg := elasticsearch.CheckElasticsearch(*host, *user, *pass, *timeout, *skipTLS)
		fmt.Println(msg)
		os.Exit(code)
	default:
		fmt.Println("UNKNOWN - parametro --check non valido (usa elasticsearch|kibana|logstash)")
		os.Exit(3)
	}
}
