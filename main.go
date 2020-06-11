package main

import (
	"fmt"
	"github.com/miekg/dns"
	"log"
	"strconv"
)

var records = map[string]string {
	"a.service.": "192.168.0.1",
	"b.service.": "192.168.0.2",
	"test.com.": "192.168.0.3", // just for test pattern when create HandleFunc
}

func parseQuery(msg *dns.Msg){
	for _,q := range msg.Question {
		switch q.Qtype {
			case dns.TypeA: { // Address Record
				ip := records[q.Name]
				if ip != "" {
					rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip)) // create new resource record
					if err != nil {
						log.Fatalf("Failed to create RR: %s\n ", err.Error())
					}
					msg.Answer = append(msg.Answer, rr)
				}
			}
		}
	}
}
func handler(writer dns.ResponseWriter, reqMsg *dns.Msg){
	replyMsg := &dns.Msg{}
	replyMsg.SetReply(reqMsg) //create a reply message from a request message.
	replyMsg.Compress = false

	switch reqMsg.Opcode {
		case dns.OpcodeQuery:
			parseQuery(replyMsg)
	}
	writer.WriteMsg(replyMsg)
}

func main()  {
	dns.HandleFunc("service.", handler) // request message should have domain that follow the pattern

	port := 5123
	server := &dns.Server{
		Addr:              ":" + strconv.Itoa(port),
		Net:               "udp", //UDP protocol
	}
	log.Printf("Started at : %d", port)

	err := server.ListenAndServe()
	defer server.Shutdown()

	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
