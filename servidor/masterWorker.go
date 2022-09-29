/*
* AUTOR: Rafael Tolosana Calasanz
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2021
* FICHERO: server.go
* DESCRIPCIÓN: contiene la funcionalidad esencial para realizar los servidores
*				correspondientes a la práctica 1
 */
package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sisdis-pr-1/com"
)

var connections = make(chan com.AUX1)
var res = make(chan com.AUX2)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

// PRE: verdad
// POST: IsPrime devuelve verdad si n es primo y falso en caso contrario
func IsPrime(n int) (foundDivisor bool) {
	foundDivisor = false
	for i := 2; (i < n) && !foundDivisor; i++ {
		foundDivisor = (n%i == 0)
	}
	return !foundDivisor
}

// PRE: interval.A < interval.B
// POST: FindPrimes devuelve todos los números primos comprendidos en el
// 		intervalo [interval.A, interval.B]
func FindPrimes(interval com.TPInterval) (primes []int) {
	for i := interval.A; i <= interval.B; i++ {
		if IsPrime(i) {
			primes = append(primes, i)
		}
	}
	return primes
}

func getParam(id int, key string, dfl string) (string) {
	value, defined := os.LookupEnv(key)
	if defined { return value }
	if len(os.Args) >= id + 1 && os.Args[id] != "" { return os.Args[id] }
	return dfl
}

//----------------------------------------------------------------------
// Worker
//----------------------------------------------------------------------

// Tratamiento para generar el resultado de la operación
func handler() {
	for {
		data := <- connections
		// Obtener los primos del intervalo
		primos := FindPrimes(data.Request.Interval)
		primos_reply := com.Reply{Id: data.Request.Id, Primes: primos}
		connReply := com.AUX2{C: data.C, Reply: primos_reply}
		res <- connReply
	}
}

// Crea un pool de allocators que se encargan de recibir las peticiones
func createHandlersPool(num int) {
	for i := 0; i < num; i++ {
		go handler()	// crear goroutine para leer peticiones
	}
}

//----------------------------------------------------------------------
// Master
//----------------------------------------------------------------------

// Gorutine que recibe petciones de los clientes y reparte
// el trabajo entre los workers
func allocate(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		checkError(err)

		dec := gob.NewDecoder(conn)
		var req com.Request
		dec.Decode(&req)
		request := com.Request{Id: req.Id, Interval: req.Interval}
		connection := com.AUX1{C: conn, Request: request}
		connections <- connection
	}
}

// Crea un pool de allocators que se encargan de recibir las peticiones
func createAllocatorPool(num int, listener net.Listener) {
	for i := 0; i < num; i++ {
		go allocate(listener)	// crear goroutine para leer peticiones
	}
}

// Gorutine que recibe mensajes de los workers con los resultados
// de las operaciones y los envía al cliente correspondiente
func response() {
	for {
		msg := <- res
		enc := gob.NewEncoder(msg.C)
		enc.Encode(msg.Reply)
		msg.C.Close()
	}
}

// Crea un pool de delivery workers que se encargan de devolver el resultado al cliente
func createDeliverPool(num int) {
	for i := 0; i < num; i++ {
		go response()	// crear goroutine para enviar resultados
	}
}

func main() {
	CONN_TYPE := getParam(2, "TYPE", "tcp")
	CONN_HOST := getParam(3, "HOST", "127.0.0.1")
	CONN_PORT := getParam(4, "PORT", "5000")
	// información de los parametros
	fmt.Printf("Listening in: %s:%s\n", CONN_HOST, CONN_PORT)
	
	listener, err := net.Listen(CONN_TYPE, CONN_HOST + ":" + CONN_PORT)
	checkError(err)
	
	// crear pool de Allocators
	createDeliverPool(6)
	createHandlersPool(6)
	createAllocatorPool(6, listener)

	for {
		// esperar a que se cierre el programa
	}
}