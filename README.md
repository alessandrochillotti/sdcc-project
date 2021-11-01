# Sistemi Distribuiti e Cloud Computing - Alessandro Chillotti (mat. 0299824)
## Multicast totalmente e causalmente ordinato in Go
### Struttura del progetto
All'interno di questa repository è possibile trovare:
- La directory `app`, contenente l'applicativo.
- Il report, contenente la spiegazione delle scelte progettuali.

Per seguire i prossimi passi, sarà necessario posizionarsi nella directory `app`.
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

Per eseguire il frontend, si necessita di posizionarsi all'interno della directory `frontend`. All'interno è presente un `Makefile` che contiene le seguenti regole:
- `frontend`: viene compilato il file relativo al frontend.
- `test`: viene compilato il file relativo al test.
- `clean`: vengono rimossi gli eseguibili creati.

In questo caso, si otterrà l'eseguibile dell'applicazione con il seguente comando.
```[bash]
make frontend
```
Verrà generato il file `frontend` eseguibile ed è possibile lanciarlo nel seguente modo.
```[bash]
./frontend
```
Inoltre, è possibile specificare il flag verbose per avere più dettagli relativi ai messaggi consegnati al livello applicativo. In particolare:
```[bash]
./frontend -V
```
### Testing dell'applicazione
Per testare l'applicazione, dopo averla istanziata ed essersi posizionati all'interno della directory `frontend`, bisogna compilare il file con il seguente comando. 
```[bash]
make test
```
Successivamente, si lancia l'eseguibile appena generato con il seguente comando.
```[bash]
./test
```
Si aprirà un menù che permettera di scegliere quale dei due test, per quello specifico algoritmo, si desidera eseguire. In particolare, per ogni algoritmo sono stati sviluppate due tipologie di test:
- Un solo peer effettua l'evento di invio del messaggio in multicast.
- Più peer, in modo concorrente, effettuano l'evento di invio del messaggio in multicast.

A fine esecuzione, sarò mostrato un output che permetterà di capire l'esito del test effettuato.
### Rimozione container
Per rimuove i container creati dall'esecuzione dell'applicazione utilizzare il seguente comando:
```[bash]
./stop.sh 
```
