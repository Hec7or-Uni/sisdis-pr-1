# sisdis-pr-1
Aplicación distribuida de cálculo de número primos

## Getting Started

### Cliente

```bash
go run client.go <ip> <port>
```
### Servidor

#### Secuencial

```bash
go run ./server/server-draft.go -s <ip> <port>
```

#### Con una Goroutine por petición

```bash
go run ./server/server-draft.go -cspf <ip> <port>
```

#### Con un pool fijo de Goroutines

```bash
go run ./server/server-draft.go -cpf <ip> <port> <num gorutines>
```

#### Master Worker

variables de configuración en el codigo master.go para poder lanzar correctamente el script que ejecutara los workers.
`NIP`: identificador del alumno para logearse con el ssh
`SRC_PATH`: dirección del ejecutable -> "/home/NIP/Desktop/sisdis-pr-1/"

```bash
go run master.go <ip> <port> <num workers>
```