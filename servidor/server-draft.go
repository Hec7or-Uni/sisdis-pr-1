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

func handleSequential(conn net.Conn) {
	defer conn.Close()
	
	// encoder & decoder
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	
	// Recibimos el intervalo
	var req com.Request
	dec.Decode(&req)
	
	// Obtener los primos del intervalo
	primos := FindPrimes(req.Interval)
	primos_reply := com.Reply{Id: req.Id, Primes: primos}
	err := enc.Encode(primos_reply)
	checkError(err)
}

func handleCSPF() {
	fmt.Println("Not implemented")
}

func handleCPF(ch chan net.Conn) {
	for {
		var conn net.Conn
		conn = <- ch
		defer conn.Close()
			
		// encoder & decoder
		enc := gob.NewEncoder(conn)
		dec := gob.NewDecoder(conn)
		
		// Recibimos el intervalo
		var req com.Request
		dec.Decode(&req)
		
		// Obtener los primos del intervalo
		primos := FindPrimes(req.Interval)
		primos_reply := com.Reply{Id: req.Id, Primes: primos}
		err := enc.Encode(primos_reply)
		checkError(err)
	}
}

func handleMW() {
	fmt.Println("Not implemented")
}

func main() {
	ALG := getParam(1, "ALG", "-s")
	CONN_TYPE := getParam(2, "TYPE", "tcp")
	CONN_HOST := getParam(3, "HOST", "127.0.0.1")
	CONN_PORT := getParam(4, "PORT", "5000")
	// información de los parametros
	fmt.Printf("Listening in: %s:%s\n", CONN_HOST, CONN_PORT)

	listener, err := net.Listen(CONN_TYPE, CONN_HOST + ":" + CONN_PORT)
	checkError(err)

	switch ALG {
		case "-s":
			for {
				conn, err := listener.Accept()
				checkError(err)
				handleSequential(conn)
			}
		case "-cspf":
			fmt.Println("Concurrente sin pool fijo de gorutines")
		case "-cpf":
			var ch chan net.Conn
			ch = make(chan net.Conn)
			for i := 0; i < 6; i++ {
				go handleCPF(ch)
			}

			for {
				conn, err := listener.Accept()
				checkError(err)
				ch <- conn
			}
		case "-mw":
			fmt.Println("Master Worker")
		default:
			fmt.Println("Undefined")
	}
}