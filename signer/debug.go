package signer

import "log"

var debugEnabled = false

func debugf(format string, v ...interface{}) {
	if debugEnabled {
		log.Printf(format, v...)
	}
}
func debugln(v ...interface{}) {
	if debugEnabled {
		log.Println(v...)
	}
}
