package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sisdis-pr-1/com"
	"time"
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

//----------------------------------------------------------------------
// Worker
//----------------------------------------------------------------------

// Tratamiento para generar el resultado de la operación
func handler(conn net.Conn) {
	defer conn.Close()
	
	// encoder & decoder
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	
	// Recibimos el intervalo
	var req com.TPInterval
	dec.Decode(&req)
	
	// Obtener los primos del intervalo
	texStart := time.Now()
	primos := FindPrimes(req)
	texEnd := time.Now()
	reply:= com.CustomReply{Primes: primos, T: texEnd.Sub(texStart)}
	err := enc.Encode(reply)
	checkError(err)
}

func main() {
	CONN_TYPE := getParam(1, "TYPE", "tcp")
	CONN_HOST := getParam(2, "HOST", "127.0.0.1")
	CONN_PORT := getParam(3, "PORT", "5001")
	// información de los parametros
	fmt.Printf("Listening in: %s:%s\n", CONN_HOST, CONN_PORT)
	
	listener, err := net.Listen(CONN_TYPE, CONN_HOST + ":" + CONN_PORT)
	checkError(err)
	
	for {
		conn, err := listener.Accept()
		checkError(err)
		handler(conn)
	}
}