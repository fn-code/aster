package aster

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

const (
	// StatusOK is status ok
	StatusOK = byte(0x80)
	// StatusNotOK is status is not ok
	StatusNotOK = byte(0x81)
	// PlayAudio is command to play audio
	PlayAudio = byte(0x33)
	// StopAudio is command to stop audio
	StopAudio = byte(0x34)
)

// AudioServer contain output from audio server
type AudioServer struct {
	pc     net.PacketConn
	buffer []byte
	done   chan struct{}
}

// AudioConn hold udp connection to the server
type AudioConn struct {
	conn     *net.UDPConn
	bufferIn []byte
}

const (
	// linux
	// ffmpeg = "ffmpeg"
	// audio  = "-f alsa -ac 2 -i default -acodec libmp3lame -ab 32k -f rtp rtp://192.168.0.104:1234"
	// windows
	ffmpeg  = "./ffmpeg.exe"
	micName = "Microphone (High Definition Audio Device)"
	base    = "-f dshow -i"
	ffplay  = "./ffmpeg-4.1/ffplay"
)

var addr = "192.168.31.15:1234"

// runAudio is starting audio
func runAudioClient() *exec.Cmd {
	remote := fmt.Sprintf("rtp://%s", addr)
	cmd := exec.Command(ffplay, "-nodisp", "-i", remote)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd
}

func runAudioStream(addrs []string) (*exec.Cmd, error) {
	// if len(addrs) == 0 {
	// 	return errEmptyAddress
	// }

	client := []string{}
	for _, v := range addrs {
		// -ab 32k
		client = append(client, fmt.Sprintf("-acodec libmp3lame -ab 32k -f rtp rtp://%s", v))
	}
	aud := fmt.Sprintf("audio=%s", micName)

	sub := strings.Split(strings.Join(client, " "), " ")
	mst := strings.Split(base, " ")
	mst = append(mst, aud)
	mst = append(mst, sub...)

	cmd := exec.Command(ffmpeg, mst...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd, nil
}

func (pc *AudioServer) sendStatus(status byte, addr net.Addr) error {
	msg := []byte{status}
	_, err := pc.pc.WriteTo(msg, addr)
	if err != nil {
		return err
	}
	return nil
}

func (pc *AudioServer) readHeader(addr net.Addr) (bool, error) {
	err := checkHeader(pc.buffer)
	if err != nil {
		if err == errMaskInvalid || err == errModeInvalid {
			err := pc.sendStatus(StatusNotOK, addr)
			if err != nil {
				return false, err
			}
		}
		return false, err
	}
	return true, nil
}

func checkHeader(r []byte) error {
	if len(r) == 0 {
		return errEmptyByte
	}
	if r[0] != 0x81 {
		return errModeInvalid
	}
	if (r[1] & 0x80) != 0x80 {
		return errMaskInvalid
	}
	return nil
}
