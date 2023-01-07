package tracer

import (
	"encoding/json"
	"fmt"
	"time"
)

type Tracer struct {
	Started time.Time
	Events  []Event
	bank    *Bank
	stack   []string
}

var (
	// Disabled - a tracer singleton that never logs
	Disabled = &Tracer{}
	// MaxEventCapacity - tracer will stop logging if exceeded
	MaxEventCapacity = 5000
)

func begin(b *Bank) *Tracer {
	return &Tracer{
		Started: time.Now(),
		bank:    b,
	}
}

var printOutput = false

func PrintAllOutput() {
	printOutput = true
}

func (t *Tracer) Disabled() bool {
	return t.Started.IsZero() || len(t.Events) >= MaxEventCapacity
}

func (t *Tracer) append(e Event) {
	t.Events = append(t.Events, e)
	if printOutput {
		result, _ := json.Marshal(e)
		fmt.Println(string(result))
	}
}

func (t *Tracer) AliasProcess(id int, name string) {
	if t.Disabled() {
		return
	}
	t.append(Event{
		Name:  "process_name",
		Phase: "M",
		PID:   id,
		Args: map[string]interface{}{
			"name": name,
		},
	})
}

func (t *Tracer) Begin(name string) {
	if t.Disabled() {
		return
	}
	t.stack = append(t.stack, name)
	t.append(Event{
		Timestamp: time.Since(t.Started).Microseconds(),
		Name:      name,
		PID:       1,
		TID:       1,
		Phase:     "B",
	})
}

func (t *Tracer) Arg(kvs map[string]interface{}) {

	if t.Disabled() {
		return
	}
	t.Events[len(t.Events)-1].Args = kvs
}

func (t *Tracer) End() {
	if t.Disabled() {
		return
	}

	elapsed := time.Since(t.Started)

	name := t.stack[len(t.stack)-1]

	counter := t.bank.gross[name]
	counter.cnt += 1
	counter.dur += elapsed
	t.bank.gross[name] = counter

	t.stack = t.stack[0 : len(t.stack)-1]
	t.append(Event{
		Timestamp: elapsed.Microseconds(),
		PID:       1,
		TID:       1,
		Phase:     "E",
	})
}

func (t *Tracer) Close() {
	t.bank.collect(t)
}
