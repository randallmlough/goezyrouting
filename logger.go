package goezyrouting

import "log"

type Logger interface {
	Log(msg string, v ...interface{})
	Debug(msg string, v ...interface{})
	Error(msg string, v ...interface{})
	Fatal(msg string, v ...interface{})
	Warn(msg string, v ...interface{})
	Panic(msg string, v ...interface{})
}

func NewLogger() *Log {
	return &Log{}
}

type Log struct{}

var offset = 2

func (x *Log) ErrorLog(v ...interface{}) *log.Logger {
	return &log.Logger{}
}
func (x *Log) Log(msg string, v ...interface{}) {
	ii := make([]interface{}, len(v)+offset, len(v)+offset)
	ii[0] = "Info:"
	ii[1] = msg
	for i, vv := range v {
		ii[i+offset] = vv
	}
	log.Println(ii...)
}
func (x *Log) Debug(msg string, v ...interface{}) {
	ii := make([]interface{}, len(v)+offset, len(v)+offset)
	ii[0] = "Debug:"
	ii[1] = msg
	for i, vv := range v {
		ii[i+offset] = vv
	}
	log.Println(ii...)
}
func (x *Log) Error(msg string, v ...interface{}) {
	ii := make([]interface{}, len(v)+offset, len(v)+offset)
	ii[0] = "Error:"
	ii[1] = msg
	for i, vv := range v {
		ii[i+offset] = vv
	}
	log.Println(ii...)
}
func (x *Log) Fatal(msg string, v ...interface{}) {
	ii := make([]interface{}, len(v)+offset, len(v)+offset)
	ii[0] = "Fatal:"
	ii[1] = msg
	for i, vv := range v {
		ii[i+offset] = vv
	}
	log.Fatal(ii...)
}
func (x *Log) Warn(msg string, v ...interface{}) {
	ii := make([]interface{}, len(v)+offset, len(v)+offset)
	ii[0] = "Warning:"
	ii[1] = msg
	for i, vv := range v {
		ii[i+offset] = vv
	}
	log.Println(ii...)
}
func (x *Log) Panic(msg string, v ...interface{}) {
	ii := make([]interface{}, len(v)+offset, len(v)+offset)
	ii[0] = "Panic:"
	ii[1] = msg
	for i, vv := range v {
		ii[i+offset] = vv
	}
	log.Panic(ii...)
}
