package engine

import (
	"giot/internal/core/model"
	"giot/internal/processor"
)

var EnChan = make(chan processor.RemoteData, 1024)

func loop() {

	data := <-EnChan
}

type Interface interface {
	Trigger(s string) error
}

type engine struct {
	de model.Device
}

func (engine *engine) Trigger(s string) error {
	for _, v := range engine.de.Rule {
		for _, tr := range v.Triggers {
			tr.
		}
	}
}
