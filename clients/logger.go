package clients

import "log"

type Logger struct{}

func (l Logger) Debug(msg string, keyvals ...interface{}) {
	printed := []interface{}{"[DEBUG]", msg}
	printed = append(printed, keyvals)
	log.Println(printed...)
}

func (l Logger) Info(msg string, keyvals ...interface{}) {
	printed := []interface{}{"[INFO]", msg}
	printed = append(printed, keyvals)
	log.Println(printed...)
}

func (l Logger) Warn(msg string, keyvals ...interface{}) {
	printed := []interface{}{"[WARN]", msg}
	printed = append(printed, keyvals)
	log.Println(printed...)
}

func (l Logger) Error(msg string, keyvals ...interface{}) {
	printed := []interface{}{"[ERROR]", msg}
	printed = append(printed, keyvals)
	log.Println(printed...)
}
