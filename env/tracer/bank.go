package tracer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type counter struct {
	dur time.Duration
	cnt int
}

type Bank struct {
	dur    time.Duration
	tracer *Tracer

	gross map[string]counter
}

func NewBank() *Bank {
	return &Bank{gross: map[string]counter{}}
}

func (b *Bank) Clear() {
	b.tracer = nil
	b.dur = 0
	b.gross = map[string]counter{}
}

type KV struct {
	name string
	cnt  int
	dur  time.Duration
}

func (b *Bank) SaveReport(file string) {

	if b.tracer == nil {
		return
	}

	var kvs []KV

	for k, v := range b.gross {
		kvs = append(kvs, KV{name: k, dur: v.dur, cnt: v.cnt})
	}
	sort.SliceStable(kvs, func(i, j int) bool {
		return kvs[i].dur > kvs[j].dur
	})

	f, err := os.Create(file)
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()

	for _, i := range kvs {

		line := strings.Replace(i.name, "\n", " ", -1)
		line = strings.TrimSpace(line)
		f.WriteString(fmt.Sprintf("%12v\t%6d\t%s\n", i.dur, i.cnt, line))

	}

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
	return begin(b)
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
