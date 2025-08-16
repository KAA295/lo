package pkg

import "log"

type LogEntry struct {
	Action  string
	Message string
}

type Logger struct {
	logChannel  chan LogEntry
	doneChannel chan struct{}
}

func NewLogger(bufferSize int) *Logger {
	return &Logger{
		logChannel:  make(chan LogEntry, bufferSize),
		doneChannel: make(chan struct{}),
	}
}

func (l *Logger) Start() {
	go func() {
		defer close(l.doneChannel)
		for entry := range l.logChannel {
			log.Printf("[%s] %s", entry.Action, entry.Message)
		}
	}()
}

func (l *Logger) Log(action string, message string) {
	select {
	case l.logChannel <- LogEntry{Action: action, Message: message}:
	default:
		log.Println("Log channel is full, dropping message")
	}
}

func (l *Logger) Stop() {
	close(l.logChannel)
	<-l.doneChannel
}
