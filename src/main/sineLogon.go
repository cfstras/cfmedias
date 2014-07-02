package main

import (
	"math"
	"time"

	"code.google.com/p/portaudio-go/portaudio"
)

const sampleRate = 44100

func sineLogon() {
	portaudio.Initialize()
	defer portaudio.Terminate()
	s := newStereoSine(132, 198/2, sampleRate)
	defer s.Close()
	chk(s.Start())

	doADSR(s, 0.4, 0.1, 0.1)
	doADSR(s, 0.1, 0.8, 0.9)

	// wait for buffer to drain
	time.Sleep(500 * time.Millisecond)
	chk(s.Stop())
}

func doADSR(s *stereoSine, attack, decay, release float64) {
	// attack
	s.fade = 1.0 / (sampleRate * attack)
	<-s.control

	// decay/sustain
	time.Sleep(time.Duration(decay*1000.0) * time.Millisecond)
	//release
	s.fade = -1.0 / (sampleRate * release)

	<-s.control
}

type stereoSine struct {
	*portaudio.Stream
	step1, phase1, pan1 float64
	step2, phase2, pan2 float64
	volume              float64
	fade                float64

	control chan bool
}

func newStereoSine(freqL, freqR, sampleRate float64) *stereoSine {
	s := &stereoSine{nil,
		// timestep, phase, pan
		freqL / sampleRate, 0, -1,
		freqR / sampleRate, 0, 1,
		0, 0, make(chan bool)}
	var err error
	s.Stream, err = portaudio.OpenDefaultStream(0, 2, sampleRate, 0, s.processAudio)
	chk(err)
	return s
}

func (g *stereoSine) processAudio(out [][]float32) {
	for i := range out[0] {
		val1 := g.volume * math.Sin(2*math.Pi*g.phase1)
		val2 := g.volume * math.Sin(2*math.Pi*g.phase2)

		val1l := math.Max(val1*g.pan1, 0) * 0.5
		val2l := math.Max(val2*g.pan2, 0) * 0.5

		val1r := math.Max(val1*-g.pan1, 0) * 0.5
		val2r := math.Max(val2*-g.pan2, 0) * 0.5

		out[0][i] = float32(val1l + val2l)
		out[1][i] = float32(val1r + val2r)

		_, g.phase1 = math.Modf(g.phase1 + g.step1)
		_, g.phase2 = math.Modf(g.phase2 + g.step2)

		g.volume += g.fade
		if g.volume < 0 {
			g.volume = 0
			g.fade = 0
			g.control <- true
		}
		if g.volume > 1 {
			g.volume = 1
			g.fade = 0
			g.control <- true
		}
	}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
