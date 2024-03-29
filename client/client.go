/*
* AUTOR: Rafael Tolosana Calasanz
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2021
* FICHERO: client.go
* DESCRIPCIÓN: cliente completo para los cuatro escenarios de la práctica 1
 */
package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sisdis-pr-1/com"
	"time"
)

// sendRequest envía una petición (id, interval) al servidor. Una petición es un par id
// (el identificador único de la petición) e interval, el intervalo en el cual se desea que el servidor encuentre los
// números primos. La petición se serializa utilizando el encoder y una vez enviada la petición
// se almacena en una estructura de datos, junto con una estampilla
// temporal. Para evitar condiciones de carrera, la estructura de datos compartida se almacena en una Goroutine
// (handleRequests) y que controla los accesos a través de canales síncronos. En este caso, se añade una nueva
// petición a la estructura de datos mediante el canal addChan
func sendRequest(endpoint string, id int, interval com.TPInterval, addChan chan com.TimeRequest, delChan chan com.TimeReply){
    tcpAddr, err := net.ResolveTCPAddr("tcp", endpoint)
    com.CheckError(err)

    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    com.CheckError(err)

    encoder := gob.NewEncoder(conn)
    decoder := gob.NewDecoder(conn)
    request := com.Request{Id: id, Interval: interval}
    timeReq := com.TimeRequest{Id: id, T: time.Now()}
    err = encoder.Encode(request)
    com.CheckError(err)
    addChan <- timeReq
    go receiveReply(decoder, delChan, conn)
}

// handleRequests es una Goroutine que garantiza el acceso en exclusión mutua a la tabla de peticiones. La tabla de peticiones
// almacena todas las peticiones activas que se han realizado al servidor y cuándo se han realizado. El objetivo es que el cliente
// pueda calcular, para cada petición, cuál es el tiempo total desde que se envía hasta que se recibe.
// Las peticiones le llegan a la goroutine a través del canal addChan. Por el canal delChan se
// indica que ha llegado una respuesta de una petición. En la respuesta, se obtiene también el timestamp de la recepción.
// Antes de eliminar una petición se imprime por la salida estándar el id de una petición y el tiempo transcurrido, observado
// por el cliente (tiempo de transmisión + tiempo de overheads + tiempo de ejecución efectivo)
func handleRequests(addChan chan com.TimeRequest, delChan chan com.TimeReply, done chan bool, MAX_REQ int) {
    requests := make(map[int]time.Time)
    i := 0
    for {
        select {
            case request := <- addChan:
                requests[request.Id] = request.T
            case reply := <- delChan:
                fmt.Println(reply.Id, " ", reply.T.Sub(requests[reply.Id]))
                delete(requests, reply.Id)
                i++;
                if i == MAX_REQ {
                    done <- true
                }
        }
    }
}

// receiveReply recibe las respuestas (id, primos) del servidor. Respuestas que corresponden con peticiones previamente
// realizadas. 
// el encoder y una vez enviada la petición se almacena en una estructura de datos, junto con una estampilla
// temporal. Para evitar condiciones de carrera, la estructura de datos compartida se almacena en una Goroutine
// (handleRequests) y que controla los accesos a través de canales síncronos. En este caso, se añade una nueva
// petición a la estructura de datos mediante el canal addChan
func receiveReply(decoder *gob.Decoder, delChan chan com.TimeReply, conn net.Conn){
    var reply com.Reply
    err := decoder.Decode(&reply)
    com.CheckError(err)
    timeReply := com.TimeReply{Id: reply.Id, T: time.Now()}
    delChan <- timeReply 
	conn.Close()
}

func getParam(id int, dfl string) (string) {
	if len(os.Args) >= id + 1 && os.Args[id] != "" { return os.Args[id] }
	return dfl
}

func main() {
	CONN_HOST := getParam(1, "127.0.0.1")
	CONN_PORT := getParam(2, "5000")

    endpoint := CONN_HOST + ":" + CONN_PORT
    numIt := 10
    requestTmp := 6
    interval := com.TPInterval{A: 1000, B: 70000}
    tts := 3000 // time to sleep between consecutive requests

    addChan := make(chan com.TimeRequest)
    delChan := make(chan com.TimeReply)
	done := make(chan bool)

    go handleRequests(addChan, delChan, done, numIt * requestTmp)
    
    for i := 0; i < numIt; i++ {
        for t := 1; t <= requestTmp; t++{
            sendRequest(endpoint, i * requestTmp + t, interval, addChan, delChan)
        }
        time.Sleep(time.Duration(tts) * time.Millisecond)
    }

    <- done
}
