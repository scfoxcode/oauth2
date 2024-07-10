package counters 

import (
    "sync"
)

type Counters struct {
    counters map[string]int
    mu       sync.Mutex
}

func (r *Counters) Init() {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.counters = make(map[string]int)
}

func (r *Counters) Increment(ip string) {
    r.mu.Lock()
    defer r.mu.Unlock()
    count, ok := r.counters[ip]
    if !ok {
        r.counters[ip] = 0
    } else {
        r.counters[ip] = count + 1
    }
}

func (r *Counters) GetCount(ip string) int {
    r.mu.Lock()
    defer r.mu.Unlock()
    count, ok := r.counters[ip]
    if !ok {
        return 0
    } else {
        return count
    }
}

func (r *Counters) Clear() {
    r.Init()
}


