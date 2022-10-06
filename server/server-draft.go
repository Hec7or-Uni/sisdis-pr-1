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
	"time"
)

const MAX_GORUTINES = 6

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

func getParam(id int, dfl string) (string) {
	if len(os.Args) >= id + 1 && os.Args[id] != "" { return os.Args[id] }
	return dfl
}

func handleSimple(conn net.Conn) {
	timeStart := time.Now()
	defer conn.Close()

	// encoder & decoder
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	// Recibimos el intervalo
	var req com.Request
	txonStart1 := time.Now()
	dec.Decode(&req)
	txonEnd1 := time.Now()

	// Obtener los primos del intervalo
	texStart := time.Now()
	primos := FindPrimes(req.Interval)
	texEnd := time.Now()
	primos_reply := com.Reply{Id: req.Id, Primes: primos}

	txonStart2 := time.Now()
	err := enc.Encode(primos_reply)
	com.CheckError(err)
	txonEnd2 := time.Now()

	txon := txonEnd1.Sub(txonStart1) + txonEnd2.Sub(txonStart2) // tiempo de transmisión
	tex := texEnd.Sub(texStart)																	// tiempo de ejecución
	to := time.Since(timeStart) - txon - tex										// tiempo de espera (overhead)
	fmt.Println(conn.RemoteAddr().String()[:15], "\t", req.Id, "\t", txon, "\t", tex, "\t", to)
}

func handleCPF(ch chan net.Conn) {
	for {
		timeStart := time.Now()
		var conn net.Conn
		conn = <-ch
		defer conn.Close()

		// encoder & decoder
		enc := gob.NewEncoder(conn)
		dec := gob.NewDecoder(conn)

		// Recibimos el intervalo
		var req com.Request
		txonStart1 := time.Now()
		dec.Decode(&req)
		txonEnd1 := time.Now()

		// Obtener los primos del intervalo
		texStart := time.Now()
		primos := FindPrimes(req.Interval)
		texEnd := time.Now()
		primos_reply := com.Reply{Id: req.Id, Primes: primos}

		txonStart2 := time.Now()
		err := enc.Encode(primos_reply)
		com.CheckError(err)
		txonEnd2 := time.Now()

		txon := txonEnd1.Sub(txonStart1) + txonEnd2.Sub(txonStart2) // tiempo de transmisión
		tex := texEnd.Sub(texStart)																	// tiempo de ejecución
		to := time.Since(timeStart) - txon - tex										// tiempo de espera (overhead)
		fmt.Println(conn.RemoteAddr().String()[:15], "\t", req.Id, "\t", txon, "\t", tex, "\t", to)
	}
}

func main() {
	ALG := getParam(1, "-s")
	CONN_HOST := getParam(2, "127.0.0.1")
	CONN_PORT := getParam(3, "5000")
	// información de los parametros
	fmt.Printf("Listening in: %s:%s\n", CONN_HOST, CONN_PORT)

	listener, err := net.Listen("tcp", CONN_HOST+":"+CONN_PORT)
	com.CheckError(err)

	switch ALG {
	case "-s":
		for {
			conn, err := listener.Accept()
			com.CheckError(err)
			handleSimple(conn)
		}
	case "-cspf":
		for {
			conn, err := listener.Accept()
			com.CheckError(err)
			go handleSimple(conn)
		}
	case "-cpf":
		ch := make(chan net.Conn)
		for i := 0; i < MAX_GORUTINES; i++ {
			go handleCPF(ch)
		}

		for {
			conn, err := listener.Accept()
			com.CheckError(err)
			ch <- conn
		}
	default:
		fmt.Println("Undefined")
	}
}
