# Sistemi Distribuiti e Cloud Computing - Alessandro Chillotti (mat. 0299824)
## Multicast totalmente e causalmente ordinato in Go
### Istanziazione dell'applicazione
Per istanziare l'applicazione utilizzare il seguente comando:
```[bash]
./start.sh ALGO NODES
```
dove NODES è il numero di peer appartenti al gruppo multicast e ALGO può assumere i seguenti valori:
- 1: multicast totalmente ordinato implementato in modo centralizzato tramite un sequencer
- 2: multicast totalmente ordinato implementato in modo decentralizzato tramite l’uso di clock logici
scalari
- 3: multicast causalmente ordinato implementato in modo decentralizzato tramite l’uso di clock
logici vettoriali
### Utilizzo dell'applicazione
Per l'utilizzazione dell'applicazione è stato sviluppato un semplice frontend che permette di interagire con ciascun container potendo effettuare:
- L'invio dei messaggi.
- Stampa della lista dei messaggi consegnati a livello applicativo.
Per eseguire il frontend, dopo aver instaziato l'applicazione, basta avviare l'applicazione in Go con il seguente comando:
```[bash]
go run frontend.go
```
Inoltre, è possibile specificare il flag verbose per avere più dettagli relativi ai messaggi consegnati al livello applicativo. In particolare:
```[bash]
go run frontend.go -V
```
### Rimozione container
Per rimuove i container creati dall'esecuzione dell'applicazione utilizzare il seguente comando:
```[bash]
./stop.sh 
```
