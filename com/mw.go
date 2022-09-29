/*
* AUTOR: Javier Pardos  & Hector Toral
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2022
* FICHERO: mw.go
* DESCRIPCIÓN:
 */
package com

import "net"

type AUX1 struct {
    C net.Conn
    Request Request
}

type AUX2 struct {
    C net.Conn
    Reply Reply
}

