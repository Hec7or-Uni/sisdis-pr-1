/*
* AUTOR: Rafael Tolosana Calasanz
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2021
* FICHERO: definitions.go
* DESCRIPCIÓN: contiene las definiciones de estructuras de datos necesarias para
*			la práctica 1
 */
package com

import (
	"fmt"
	"os"
	"time"
)

type TPInterval struct {
    A int
    B int
}

type Request struct {
    Id int
    Interval TPInterval
}

type TimeRequest struct {
    Id int
    T time.Time
}

type Reply struct {
    Id int
    Primes []int
}

type TimeReply struct {
    Id int
    T time.Time
}

type CustomReply struct {
    T time.Duration
    Primes []int
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}