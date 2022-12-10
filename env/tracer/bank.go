package tracer

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Bank struct {
	dur    time.Duration
	tracer *Tracer
}

func (b *Bank) Clear() {
	b.tracer = nil
	b.dur = 0
}

func (b *Bank) SaveSample(file string) {
	if b.tracer == nil {
		return
	}

	f, err := os.Create(file)
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()

	m := json.NewEncoder(f)

	f.WriteString("[")
	for i, t := range b.tracer.Events {
		if i != 0 {
			f.WriteString(",")
		}
		err := m.Encode(t)
		if err != nil {
			log.Panicln(err)
		}

	}
	f.WriteString("]")

}

func (b *Bank) Open() *Tracer {
	return new(b)
}

func (b *Bank) collect(t *Tracer) {

	if t.Disabled() {
		return
	}

	elapsed := time.Since(t.Started)
	if b.dur > elapsed {
		return

	}

	b.tracer = t
	b.dur = elapsed
}
