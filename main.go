package main

import (
	"encoding/binary"
	"github.com/gordonklaus/portaudio"
	"net/http"
	"time"
)

const sampleRate = 44100
const seconds = 1

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	buffer := make([]float32, sampleRate * seconds)

	stream, err := portaudio.OpenDefaultStream(1, 0, sampleRate,   len(buffer), func(in []float32) {
		for i := range buffer {
			buffer[i] = in[i]
		}
	})

	chk(err)
	chk(stream.Start())
	time.Sleep(time.Second * 40)
	chk(stream.Stop())
	defer stream.Close()

	http.HandleFunc("/audio", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			panic("expected http.ResponseWriter to be an http.Flusher")
		}
		w.Header().Set("Connection", "Keep-Alive")
		w.Header().Set("Transfer-Encoding", "chunked")

		for true {
			binary.Write(w, binary.BigEndian, &buffer)
			flusher.Flush() // Trigger "chunked" encoding
			return
		}

	})
}


func chk(err error) {
	if err != nil {
		panic(err)
	}
}