package utils

import (
	"log"
	"sync"
)

type PubSub struct {
	Streamers []*Streamer
	SMutex    *sync.Mutex
}

type Streamer struct {
	Id         int64
	Ch         chan string
	SubCounter int
	SCMutex    *sync.Mutex
	Alive      bool
}

func (ps *PubSub) GetStreamer(id int64) *Streamer {
	ps.SMutex.Lock()
	defer ps.SMutex.Unlock()
	h := 0
	for i := 0; i < len(ps.Streamers); i += h {
		ps.Streamers[i].SCMutex.Lock()
		if !(ps.Streamers[i].Alive){
			ps.Streamers = append(ps.Streamers[:i], ps.Streamers[i+1:]...)
			h = 0
		} else {
			h = 1
			ps.Streamers[i].SCMutex.Unlock()
		}
	}
	for _, v := range ps.Streamers {
		if v.Id == id {
			return v
		}
	}
	newStreamer := &Streamer{
		Id:         id,
		Ch:         make(chan string, 0),
		SubCounter: 0,
		SCMutex:    &sync.Mutex{},
		Alive:      true,
	}
	ps.Streamers = append(ps.Streamers, newStreamer)
	return newStreamer
}

func (s *Streamer) Subscribe() {
	s.SCMutex.Lock()
	defer s.SCMutex.Unlock()

	s.SubCounter++
}

func (s *Streamer) Unsubscribe() {
	s.SCMutex.Lock()
	defer s.SCMutex.Unlock()

	s.SubCounter--

	if s.SubCounter < 1 {
		s.Alive = false
	}
}

func (ps *PubSub) Notify(id int64, msg string) {
	var s *Streamer
	for _, v := range ps.Streamers {
		if v.Id == id {
			s = v
		}
	}

	if s == nil {
		log.Println("notifier: no subscribers")
		return
	}

	s.SCMutex.Lock()
	defer s.SCMutex.Unlock()

	for i := 0; i < s.SubCounter; i++ {
		s.Ch <- msg
	}
}
