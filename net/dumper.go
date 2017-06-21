package net

import (
	"net/http"
	"net/http/httputil"
	"time"

	log "github.com/sirupsen/logrus"
)

type RequestDumperInterface interface {
	DumpRequest(*http.Request)
	DumpResponse(*http.Response)
}

type RequestDumper struct {
}

func NewRequestDumper() RequestDumper {
	return RequestDumper{}
}

func (p RequestDumper) DumpRequest(req *http.Request) {
	dumpedRequest, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Debugf("Error dumping request\n%s\n", err)
	} else {
		log.Debugf("\nREQUEST: [%s]\n%s\n", time.Now().Format(time.RFC3339), string(dumpedRequest))
	}
}

func (p RequestDumper) DumpResponse(res *http.Response) {
	dumpedResponse, err := httputil.DumpResponse(res, true)
	if err != nil {
		log.Debugf("Error dumping response\n%s\n", err)
	} else {
		log.Debugf("\nRESPONSE: [%s]\n%s\n", time.Now().Format(time.RFC3339), string(dumpedResponse))
	}
}
