package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"sisdis-pr-1/com"
	"strconv"
	"strings"
	"time"
)

const NIP = 798095
const SRC_PATH = "/home/a798095/Desktop/sisdis-pr-1/"

var connections = make(chan net.Conn)

func getParam(id int, dfl string) (string) {
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
	com.CheckError(err)

	conn2w, err := net.DialTCP("tcp", nil, tcpAddr)
	com.CheckError(err)

	enc2w := gob.NewEncoder(conn2w)
	dec2w := gob.NewDecoder(conn2w)
	
	txonStartw1 := time.Now()
	err = enc2w.Encode(req.Interval)
	com.CheckError(err)
	txonEndw1 := time.Now()

	txonStartw2 := time.Now()
	reply := receiveReply(dec2w, conn2w)
	txonEndw2 := time.Now()
	primos_reply := com.Reply{Id: req.Id, Primes: reply.Primes}
	
	txonStart2 := time.Now()
	enc.Encode(primos_reply)
	com.CheckError(err)
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
	com.CheckError(err)
	return reply
}

//----------------------------------------------------------------------
// Master
//----------------------------------------------------------------------

func createPool(CONN_PORT string, MAX_WORKERS int) {
	// lanza los workers
	_, err := exec.Command(SRC_PATH + "launcher.sh", SRC_PATH + "lab102_machines.txt", CONN_PORT, strconv.Itoa(MAX_WORKERS), strconv.Itoa(NIP)).Output()
	if err != nil {
		fmt.Println("Error al lanzar los workers", err)
		os.Exit(1)
	}

	// leer archivo donde se encuentran las direcciones de los workers
	readFile, err := os.Open("./endpoints.txt")

	if err != nil { log.Fatal(err) }
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	readFile.Close()

	// crear pool de adminstracion para los workers
	for index, endpoint := range lines {
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
	CONN_HOST := getParam(1, "127.0.0.1")
	CONN_PORT := getParam(2, "5000")
	MAX_WORKERS, err := strconv.Atoi(getParam(3, "1"))
	if err != nil {
		fmt.Println("Error al convertir el numero de workers", err)
		os.Exit(1)
	}
	// información de los parametros
	fmt.Printf("Listening in: %s:%s\n", CONN_HOST, CONN_PORT)
	
	listener, err := net.Listen("tcp", CONN_HOST + ":" + CONN_PORT)
	com.CheckError(err)

	// crear pool de handle request
	createPool(CONN_PORT, MAX_WORKERS)

	for {
		conn, err := listener.Accept()
		com.CheckError(err)
		connections <- conn
	}
}