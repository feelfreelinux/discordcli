package core

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gordonklaus/portaudio"
	"layeh.com/gopus"
)

const (
	channels         int   = 1           // 1 for mono, 2 for stereo
	frameRate        int   = 48000       // audio sampling rate
	frameSize        int   = 960         // uint16 size of each audio frame
	maxBytes         int   = (frameSize) // max size of opus data
	talkTriggerMilis int64 = 500
)

type VoiceConnection struct {
	outputStream    *portaudio.Stream
	inputStream     *portaudio.Stream
	voiceConnection *discordgo.VoiceConnection
	outputBuffer    []int16
	closed          bool
	sendChannel     chan []int16
	recvChannel     chan *discordgo.Packet
	talkTrigger     int64
	speaking        bool
}

func CreateVoiceConnection(session *discordgo.Session, channel *discordgo.Channel) *VoiceConnection {
	h, o := portaudio.DefaultHostApi()
	chk(o)
	device := h.Devices[7]
	voiceConnection, _ := session.ChannelVoiceJoin(channel.GuildID, channel.ID, false, false)
	voiceConn := &VoiceConnection{}
	voiceConn.voiceConnection = voiceConnection
	voiceConn.sendChannel = make(chan []int16, 2)
	voiceConn.recvChannel = make(chan *discordgo.Packet, 2)
	voiceConn.talkTrigger = 0
	voiceConn.speaking = false
	go SendPCM(voiceConnection, voiceConn.sendChannel)
	go ReceivePCM(voiceConnection, voiceConn.recvChannel)
	inputParams := portaudio.LowLatencyParameters(device, nil)
	inputParams.Input.Channels = 1
	inputParams.Output.Channels = 0
	inputParams.SampleRate = float64(frameRate)
	inputParams.FramesPerBuffer = frameSize

	outputParams := portaudio.LowLatencyParameters(nil, device)
	outputParams.Input.Channels = 0
	outputParams.Output.Channels = 1
	outputParams.SampleRate = float64(frameRate)
	outputParams.FramesPerBuffer = frameSize

	voiceConn.outputBuffer = make([]int16, 960)

	voiceConn.inputStream, _ = portaudio.OpenStream(inputParams, voiceConn.processInput)
	voiceConn.outputStream, _ = portaudio.OpenStream(outputParams, &voiceConn.outputBuffer)

	return voiceConn
}

func chk(e error) {
	if e != nil {
		panic(e)
	}
}

func (vc *VoiceConnection) isSpeaking(speaking bool) {
	if speaking {
		vc.talkTrigger = makeTimestamp()
		if !vc.speaking {
			vc.speaking = true
			vc.talkTrigger = makeTimestamp()
			go vc.voiceConnection.Speaking(true)
		}
	} else {
		if vc.speaking && makeTimestamp() > (vc.talkTrigger+talkTriggerMilis) {
			vc.speaking = false
			go vc.voiceConnection.Speaking(false)
		}
	}
}
func (vc *VoiceConnection) Start() {
	vc.closed = false
	vc.speaking = false
	go vc.processOutput()
	vc.inputStream.Start()
	vc.outputStream.Start()

}

func (vc *VoiceConnection) Stop() {
	vc.inputStream.Stop()
	vc.outputStream.Stop()
	vc.voiceConnection.Close()
	vc.closed = true
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func (vc *VoiceConnection) processInput(in []int16) {
	data := make([]int16, 960)
	for x, d := range in {
		data[x] = d / 4
		if data[x] > 2000 {
			vc.isSpeaking(true)
		} else {
			vc.isSpeaking(false)
		}
	}

	if vc.speaking {
		vc.sendChannel <- data
	}
}

func (vc *VoiceConnection) processOutput() {
	for {
		if vc.closed {
			return
		}
		pcm := make([]int16, 960)
		data, ok := <-vc.recvChannel
		if ok {
			for x, d := range data.PCM {
				pcm[x] = d * 2
			}

		}
		vc.outputBuffer = pcm
		vc.outputStream.Write()
	}
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

var (
	speakers    map[uint32]*gopus.Decoder
	opusEncoder *gopus.Encoder
	mu          sync.Mutex
)

// OnError gets called by dgvoice when an error is encountered.
// By default logs to STDERR
var OnError = func(str string, err error) {
	prefix := "dgVoice: " + str

	if err != nil {
		os.Stderr.WriteString(prefix + ": " + err.Error())
	} else {
		os.Stderr.WriteString(prefix)
	}
}

// SendPCM will receive on the provied channel encode
// received PCM data into Opus then send that to Discordgo
func SendPCM(v *discordgo.VoiceConnection, pcm <-chan []int16) {
	if pcm == nil {
		return
	}

	var err error

	opusEncoder, err = gopus.NewEncoder(frameRate, channels, gopus.Audio)

	if err != nil {
		OnError("NewEncoder Error", err)
		return
	}

	for {

		// read pcm from chan, exit if channel is closed.
		recv, ok := <-pcm
		if !ok {
			OnError("PCM Channel closed", nil)
			return
		}

		// try encoding pcm frame with Opus
		opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
		if err != nil {
			OnError("Encoding Error", err)
			return
		}

		if v.Ready == false || v.OpusSend == nil {
			// OnError(fmt.Sprintf("Discordgo not ready for opus packets. %+v : %+v", v.Ready, v.OpusSend), nil)
			// Sending errors here might not be suited
			return
		}
		// send encoded opus data to the sendOpus channel
		v.OpusSend <- opus
	}
}

// ReceivePCM will receive on the the Discordgo OpusRecv channel and decode
// the opus audio into PCM then send it on the provided channel.
func ReceivePCM(v *discordgo.VoiceConnection, c chan *discordgo.Packet) {
	if c == nil {
		return
	}

	var err error

	for {
		if v.Ready == false || v.OpusRecv == nil {
			OnError(fmt.Sprintf("Discordgo not to receive opus packets. %+v : %+v", v.Ready, v.OpusSend), nil)
			return
		}

		p, ok := <-v.OpusRecv
		if !ok {
			return
		}

		if speakers == nil {
			speakers = make(map[uint32]*gopus.Decoder)
		}

		_, ok = speakers[p.SSRC]
		if !ok {
			speakers[p.SSRC], err = gopus.NewDecoder(48000, 1)
			if err != nil {
				OnError("error creating opus decoder", err)
				continue
			}
		}

		p.PCM, err = speakers[p.SSRC].Decode(p.Opus, 960, false)
		if err != nil {
			OnError("Error decoding opus data", err)
			continue
		}

		c <- p
	}
}
