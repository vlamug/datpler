package input

import (
	"log"

	"gopkg.in/mcuadros/go-syslog.v2"
)

func NewSyslog() func(listenAddr string, logLines chan<- *LogLine) {
	return func(listenAddr string, logLines chan<- *LogLine) {
		syslogChannel := make(syslog.LogPartsChannel)
		channelHandler := syslog.NewChannelHandler(syslogChannel)

		server := syslog.NewServer()

		log.Printf("Listening syslog input on %s\n", listenAddr)
		err := server.ListenUDP(listenAddr)
		if err != nil {
			log.Fatalf("could not listned address: %s, error: %s\n", listenAddr, err)
		}

		server.SetHandler(channelHandler)
		server.SetFormat(syslog.Automatic)

		err = server.Boot()
		if err != nil {
			log.Fatalf("could nod boot syslog server: %s\n", err)
		}

		// send logs to channel
		go func(logsChannel syslog.LogPartsChannel) {
			for logLine := range logsChannel {
				log.Printf("new log line: %#v\n", logLine)
				logLines <- NewLogLine(logLine["hostname"].(string), logLine["content"].(string))
			}
		}(syslogChannel)

		server.Wait()
	}
}

type LogLine struct {
	Host string
	Data string
}

func NewLogLine(host, data string) *LogLine {
	return &LogLine{Host: host, Data: data}
}
