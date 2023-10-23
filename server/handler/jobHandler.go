package handler

import (
	"net/http"

	"github.com/prajwalad101/vincent/server/util"
)

func (broker *Broker) JobHandler(w http.ResponseWriter, r *http.Request) {
	jobId := util.GenerateJobId()

	w.Write([]byte(jobId))
	return
}
