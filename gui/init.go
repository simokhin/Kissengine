package main

import (
	"MyChessEngine/engine"
	"bytes"
	"embed"
	"image/png"
	"io"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type SoundType int

const (
	MoveSound SoundType = iota
	CaptureSound
	CastleSound
	CheckSound
	PromoteSound
)

//go:embed assets/pieces/*.png
var pieceFS embed.FS

var pieceImages [15]*ebiten.Image

//go:embed assets/sounds/*.mp3
var soundFS embed.FS

var sounds [5]*audio.Player

var audioContext = audio.NewContext(44100)

func init() {
	load := func(piece engine.Piece, filename string) {
		data, err := pieceFS.ReadFile("assets/pieces/" + filename)
		if err != nil {
			log.Fatal(err)
		}

		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			log.Fatal(err)
		}

		pieceImages[piece] = ebiten.NewImageFromImage(img)
	}

	load(engine.White|engine.Pawn, "Chess_plt45.png")
	load(engine.White|engine.Knight, "Chess_nlt45.png")
	load(engine.White|engine.Bishop, "Chess_blt45.png")
	load(engine.White|engine.Rook, "Chess_rlt45.png")
	load(engine.White|engine.Queen, "Chess_qlt45.png")
	load(engine.White|engine.King, "Chess_klt45.png")
	load(engine.Black|engine.Pawn, "Chess_pdt45.png")
	load(engine.Black|engine.Knight, "Chess_ndt45.png")
	load(engine.Black|engine.Bishop, "Chess_bdt45.png")
	load(engine.Black|engine.Rook, "Chess_rdt45.png")
	load(engine.Black|engine.Queen, "Chess_qdt45.png")
	load(engine.Black|engine.King, "Chess_kdt45.png")

	loadSounds := func(sound SoundType, fileName string) {
		data, err := soundFS.ReadFile("assets/sounds/" + fileName)
		if err != nil {
			log.Fatal(err)
		}

		stream, err := mp3.DecodeWithoutResampling(bytes.NewReader(data))
		if err != nil {
			log.Fatal(err)
		}

		pcm, err := io.ReadAll(stream)
		if err != nil {
			log.Fatal(err)
		}

		sounds[sound] = audioContext.NewPlayerFromBytes(pcm)
	}

	loadSounds(MoveSound, "move.mp3")
	loadSounds(CaptureSound, "capture.mp3")
	loadSounds(CastleSound, "castle.mp3")
	loadSounds(CheckSound, "move-check.mp3")
	loadSounds(PromoteSound, "promote.mp3")
}
