# sisdis-pr-1
Aplicación distribuida de cálculo de número primos

## Run

Cliente

```bash
go run client.go 155.210.154.198 5000
```

Master

```bash
go run master.go tcp 155.210.154.198 5000
```

Workers

```bash
go run worker.go tcp 155.210.154.199 5001
go run worker.go tcp 155.210.154.199 5002
go run worker.go tcp 155.210.154.199 5003
go run worker.go tcp 155.210.154.199 5004
go run worker.go tcp 155.210.154.199 5005
go run worker.go tcp 155.210.154.199 5006
```

## cliente-servidor secuencial

```mermaid
sequenceDiagram
  Client ->> Server: Request
  activate Server
Note right of Server: getIntervals(start int, end int)
  Server ->> Client: Response
  deactivate Server

    Client_2 ->> Server: Request
  activate Server
Note left of Server: getIntervals(start int, end int)
  Server ->> Client_2: Response
  deactivate Server
```

## cliente servidor concurrente

### Sin pool fijo

```mermaid
sequenceDiagram
  Client_{1..N}->>Server: Request
  activate Server
  Client_{1..N}->>Server: Request
  activate Server
  Note right of Server: getIntervals(start int, end int)
  Server->>Client_{1..N}: Response
  deactivate Server
  Server->>Client_{1..N}: Response
  deactivate Server
```

### Con pool fijo

```mermaid
sequenceDiagram
  Client_{1..N}->>Server: Request
  activate Server
  Client_{1..N}->>Server: Request
  activate Server
  Note right of Server: getIntervals(start int, end int)
  Server->>Client_{1..N}: Response
  deactivate Server
  Server->>Client_{1..N}: Response
  deactivate Server
Note over Server: Max N Goroutines
```

## master-worker
```mermaid
sequenceDiagram
  Client ->> Server: Request
  activate Server
  Server->>Worker_1: run
  activate Worker_1
  Server->>Worker_2: run
  activate Worker_2
  Note left of Server: getIntervals(start int, end int)
  Worker_1->>Server: stop
  deactivate Worker_1
  Worker_2->>Server: stop
  deactivate Worker_2
  Server ->> Client: Response
  deactivate Server
```
