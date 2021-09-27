# Sistemi Distribuiti e Cloud Computing - Alessandro Chillotti (mat. 0299824)
## Multicast totalmente e causalmente ordinato in Go
### Istanziazione dell'applicazione
Per istanziare l'applicazione utilizzare il seguente comando:
```[bash]
ALGO=X docker-compose up --scale node=N -d
```
dove N è il numero di partecipanti al gruppo di multicast e X assume i seguenti valori:
- 1: multicast totalmente ordinato implementato in modo centralizzato tramite un sequencer
- 2: multicast totalmente ordinato implementato in modo decentralizzato tramite l’uso di clock logici
scalari
- 3: multicast causalmente ordinato implementato in modo decentralizzato tramite l’uso di clock
logici vettoriali
### Utilizzo dell'applicazione
Per l'interazione con i container istanziati in precedenza utilizzare il seguente comando:
```[bash]
docker attach app_node_Y
```
dove Y è il numero del nodo con cui si è interessati interagire. Ovviamente, Y è un numero che prende valori da 1 ad N.