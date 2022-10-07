# sisdis-pr-1
Aplicación distribuida de cálculo de número primos

## Requirements

Para la correcta ejecución de los scripts se debe tener acceso por clave publica a los laboratorios que se desee acceder.

## Getting Started

### Cliente

```bash
./client <ip> <port>
```
### Servidor

#### Secuencial

```bash
./server/server -s <ip> <port>
```

#### Con una Goroutine por petición

```bash
./server/server -cspf <ip> <port>
```

#### Con un pool fijo de Goroutines

```bash
./server/server -cpf <ip> <port> <num gorutines>
```

#### Master Worker

variables de configuración en el codigo master.go para poder lanzar correctamente el script que ejecutara los workers.
`NIP`: identificador del alumno para logearse con el ssh
`SRC_PATH`: dirección del ejecutable -> "/home/NIP/Desktop/sisdis-pr-1/"

```bash
./master <ip> <port> <num workers>
```