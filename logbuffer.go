// this is the main program package
package main

// importing required modules / packages
import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

// structure with information about logs
type logtype struct {
	timestamp int64  // timestamp in epoch time
	message   string // log message
	level     int    // log level
}

// global variables - storing 25 log entries
var debug int = 6
var log_index = 0
var log_history int = 25
var logs = make([]logtype, log_history)

// get IP address
func get_local_ip() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// logger function saving into ring buffer
func mlog(level int, a ...interface{}) {
	if level <= debug {
		logs[log_index].level = level
		logs[log_index].message = fmt.Sprintln(a...)
		logs[log_index].timestamp = time.Now().Unix()
		log_index++
		if log_index >= log_history {
			log_index = 0
		}

		// logging to stdout as well
		s := fmt.Sprintln(a...)
		log.Printf("%s", strings.TrimSpace(s))
	}
}

// clearing all logs
func clear_logs() {
	log_index = 0
	for i := 0; i < log_history; i++ {
		logs[i].message = ""
		logs[i].level = 0
		logs[i].timestamp = 0
	}
}

// show logs
func show_logs() {
	// go through all logs - with the newest entry first, then back to the first slot and after this
	for x := (log_index - 1); x >= 0; x-- {
		if logs[x].timestamp > 0 {
			log.Printf(
				"{ \"id\":%d, \"ts\":%d, \"level\": %d, \"message\": \"%s\" }",
				x,
				logs[x].timestamp,
				logs[x].level,
				strings.TrimSpace(logs[x].message),
			)
		}
	}
	for x := (log_history - 1); x >= log_index; x-- {
		if logs[x].timestamp > 0 {
			log.Printf(
				"{ \"id\":%d, \"ts\":%d, \"level\": %d, \"message\": \"%s\" }",
				x,
				logs[x].timestamp,
				logs[x].level,
				strings.TrimSpace(logs[x].message),
			)
		}
	}
}

// MAIN
func main() {

	// setting up logging to memory buffer - initialise logs
	clear_logs()

	// adding some test log entries
	mlog(4, "INFO: running, IP address:", get_local_ip())
	mlog(1, "DEBUG: debug log entry")
	mlog(3, "WARN: warn log entry")

	// outputting all logs in the buffer
	show_logs()

	mlog(1, "INFO: main thread finished here")
}
