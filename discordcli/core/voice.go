package core

import (
	"fmt"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/gordonklaus/portaudio"
	"layeh.com/gopus"
)

const (
	channels  int = 1                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) * 2 // max size of opus data
)

type VoiceConnection struct {
	outputStream    *portaudio.Stream
	inputStream     *portaudio.Stream
	voiceConnection *discordgo.VoiceConnection
	outputBuffer    []int16
	closed          bool
	sendChannel     chan []int16
	recvChannel     chan *discordgo.Packet
}

func CreateVoiceConnection(session *discordgo.Session, channel *discordgo.Channel) *VoiceConnection {
	h, o := portaudio.DefaultHostApi()
	chk(o)
	device := h.Devices[17]
	voiceConnection, _ := session.ChannelVoiceJoin(channel.GuildID, channel.ID, false, false)
	voiceConn := &VoiceConnection{}
	voiceConn.voiceConnection = voiceConnection
	voiceConn.sendChannel = make(chan []int16, 2)
	voiceConn.recvChannel = make(chan *discordgo.Packet, 2)
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

func (vc *VoiceConnection) Start() {
	vc.closed = false
	vc.voiceConnection.Speaking(true)
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

func (vc *VoiceConnection) processInput(in []int16) {
	data := make([]int16, 960)
	for x, d := range in {
		data[x] = d / 4
	}
	vc.sendChannel <- data
}

func (vc *VoiceConnection) processOutput() {
	for {
		if vc.closed {
			return
		}
		data, ok := <-vc.recvChannel
		if ok {
			pcm := make([]int16, 960)
			for x, d := range data.PCM {
				pcm[x] = int16(float32(d) * 0.4)
			}
			vc.outputBuffer = pcm
			vc.outputStream.Write()
		}
	}
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
