# monitoring-challenge

## Descrizione

Questo progetto riguarda la creazione di uno **script in linguaggio Go** per il controllo dello **stato di salute di un cluster** attraverso l'utilizzo di **Elasticsearch, Logstash e Kibana**.

## Struttura del progetto

- **cmd/**: contiene il punto d'ingresso principale dello script (`main.go`). In base ai parametri passati da linea di comando, il main chiamerà le funzioni interne corrispondenti.

- **internal/**: contiene le funzioni interne, suddivise per infrastruttura. Ad esempio:
  - **logstash/**: include gli script relativi a Logstash.
  - **elasticsearch/**: include gli script relativi a Elasticsearch.
  - **kibana/**: include gli script relativi a Kibana.

- **quickstart/**: Contiene gli ambienti di test utilizzati per testare gli script
