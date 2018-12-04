package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 320
	screenHeight = 240
	sampleRate   = 44100
	bufferSize   = 4096
)

type stream struct {
	position int64
}

func (s *stream) Read(data []byte) (int, error) {
	fillBuffer(len(data) / 4)
	for i := 0; i < len(data)/4; i++ {
		// PCM 16-bit 2 channel stereo
		data[4*i] = byte(buffer[i])
		data[4*i+1] = byte(buffer[i] >> 8)
		data[4*i+2] = byte(buffer[i])
		data[4*i+3] = byte(buffer[i] >> 8)
	}
	return len(data), nil
}

func (s *stream) Close() error {
	return nil
}

var player *audio.Player

func update(screen *ebiten.Image) error {
	if player == nil {
		var err error
		player, err = audio.NewPlayer(audioContext, &stream{})
		if err != nil {
			return err
		}
		player.Play()
	}
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	msg := fmt.Sprintf("TPS: %0.2f\n", ebiten.CurrentTPS())
	ebitenutil.DebugPrint(screen, msg)
	return nil
}

var audioContext *audio.Context
var scanner *bufio.Scanner
var buffer [bufferSize]int16

func fillBuffer(size int) {
	for i := 0; i < size; i++ {
		scanner.Scan()
		text := scanner.Text()
		value, _ := strconv.ParseInt(text, 10, 16)
		buffer[i] = int16(value)
	}
}

func main() {
	var err error
	audioContext, err = audio.NewContext(sampleRate)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Open("wave.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)

	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Sound Example"); err != nil {
		log.Fatal(err)
	}
}
