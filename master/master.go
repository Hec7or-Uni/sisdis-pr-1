package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/exec"
	"sisdis-pr-1/com"
	"strings"
	"time"
)

var connections = make(chan net.Conn)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func getParam(id int, key string, dfl string) (string) {
	value, defined := os.LookupEnv(key)
	if defined { return value }
	if len(os.Args) >= id + 1 && os.Args[id] != "" { return os.Args[id] }
	return dfl
}

func handleRequests(workerId int, endpoint string) {
	workerDown := 0
	for {
		if checkWorkerStatus(endpoint) {
			workerDown = 0
			conn := <- connections

			fmt.Println("Worker ", workerId, " handling request")

			// encoder & decoder para comunicarse con el cliente
			enc := gob.NewEncoder(conn)
			dec := gob.NewDecoder(conn)

			// enviar request a worker
			sendRequest(endpoint, enc, dec, conn.RemoteAddr().String()[:15])
		} else {
			fmt.Println("Worker ", workerId, " is down")
			time.Sleep(time.Duration(5000) * time.Millisecond)	// 5 segundos
			workerDown++
			if workerDown == 3 {
				fmt.Println("Worker ", workerId, " is not responding")
				break
			}
		}
	}
}

// Manda una request a un worker
func sendRequest(endpoint string, enc *gob.Encoder, dec *gob.Decoder, ipClient string){
	timeStart := time.Now()
	// Datos de la request
	var req com.Request
	txonStart1 := time.Now()
	dec.Decode(&req)
	txonEnd1 := time.Now()

	// Conectamos con el worker
	tcpAddr, err := net.ResolveTCPAddr("tcp", endpoint)
	checkError(err)

	conn2w, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	enc2w := gob.NewEncoder(conn2w)
	dec2w := gob.NewDecoder(conn2w)
	
	txonStartw1 := time.Now()
	err = enc2w.Encode(req.Interval)
	checkError(err)
	txonEndw1 := time.Now()

	txonStartw2 := time.Now()
	reply := receiveReply(dec2w, conn2w)
	txonEndw2 := time.Now()
	primos_reply := com.Reply{Id: req.Id, Primes: reply.Primes}
	
	txonStart2 := time.Now()
	enc.Encode(primos_reply)
	checkError(err)
	txonEnd2 := time.Now()

	txon := txonEnd1.Sub(txonStart1) + txonEnd2.Sub(txonStart2) + txonEndw1.Sub(txonStartw1) + txonEndw2.Sub(txonStartw2)	// tiempo de transmisión
	tex := reply.T														// tiempo de ejecución
	to := time.Since(timeStart) - txon - tex	// tiempo de espera (overhead)
	fmt.Println(ipClient, "\t", req.Id, "\t", txon, "\t", tex, "\t", to)
}

// Gorutine que recibe mensajes de los workers con los resultados
// de las operaciones
func receiveReply(dec2w *gob.Decoder, conn2w net.Conn) com.CustomReply {
	defer conn2w.Close()
	var reply com.CustomReply
	err := dec2w.Decode(&reply)
	checkError(err)
	return reply
}

//----------------------------------------------------------------------
// Master
//----------------------------------------------------------------------

func createPool(endpoints [3]string) {
	for index, endpoint := range endpoints {
		go handleRequests(index + 1, endpoint)	// crear goroutine para leer peticiones
	}
}

func checkWorkerStatus(endpoint string) bool {
	out, err := exec.Command("ping", endpoint[:15], "-c 1").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "0% packet loss")
}

func main() {
	CONN_TYPE := getParam(1, "TYPE", "tcp")
	CONN_HOST := getParam(2, "HOST", "127.0.0.1")
	CONN_PORT := getParam(3, "PORT", "5000")
	// información de los parametros
	fmt.Printf("Listening in: %s:%s\n", CONN_HOST, CONN_PORT)
	
	listener, err := net.Listen(CONN_TYPE, CONN_HOST + ":" + CONN_PORT)
	checkError(err)
	
	endpoints := [3]string{
		"155.210.154.197:5001",
		"155.210.154.199:5002",
		"155.210.154.200:5003", // maquina apagada 3/10/2022 <- fuerza fallo
	}

	// crear pool de handle request
	createPool(endpoints)

	for {
		conn, err := listener.Accept()
		checkError(err)
		connections <- conn
	}
}