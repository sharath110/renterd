package autopilot

import (
	"net/http"
	"time"

	"go.sia.tech/jape"
)

func (ap *Autopilot) actionsHandler(jc jape.Context) {
	var since time.Time
	max := -1
	if jc.DecodeForm("since", (*paramTime)(&since)) != nil || jc.DecodeForm("max", &max) != nil {
		return
	}
	jc.Encode(ap.Actions(since, max))
}

func (ap *Autopilot) configHandlerGET(jc jape.Context) {
	jc.Encode(ap.Config())
}

func (ap *Autopilot) configHandlerPUT(jc jape.Context) {
	var c Config
	if jc.Decode(&c) != nil {
		return
	}
	if jc.Check("failed to set config", ap.SetConfig(c)) != nil {
		return
	}
}

func NewServer(ap *Autopilot) http.Handler {
	return jape.Mux(map[string]jape.Handler{
		"GET    /actions": ap.actionsHandler,
		"GET    /config":  ap.configHandlerGET,
		"PUT    /config":  ap.configHandlerPUT,
	})
}
