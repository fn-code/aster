package aster

import (
	"fmt"
	"log"
	"net"

	"github.com/google/logger"
)

// procesStream is to start or stop stream audio
func (pc *AudioServer) procesStream(addr net.Addr) {
	res := pc.buffer[2]

	switch res {
	case PlayAudio:
		log.Println("Audio server is playing")
		pc.sendStatus(StatusOK, addr)

		// List ip to stream audio
		addrs := []string{
			"192.168.31.10:1234",
			"192.168.31.11:1234",
			"192.168.31.12:1234",
			"192.168.31.13:1234",
			"192.168.31.14:1234",
			"192.168.31.15:1234",
		}
		cmd, _ := runAudioStream(addrs)

		go func() {
			select {
			case <-pc.done:
				cmd.Process.Kill()
				return
			}
		}()

		err := cmd.Start()
		if err != nil {
			log.Println(err)
		}

	case StopAudio:
		pc.done <- struct{}{}
		log.Println("Audio server is stoped")
		pc.sendStatus(StatusOK, addr)

	default:
		log.Println("command error")
		pc.sendStatus(StatusNotOK, addr)
	}
}

// procesPlay is to start play or stop audio
func (pc *AudioServer) procesPlay(addr net.Addr) error {
	res := pc.buffer[2]

	switch res {
	case PlayAudio:
		fmt.Println("play audio")
		pc.sendStatus(StatusOK, addr)
		cmd := runAudioClient()
		go func() {

			for {
				select {
				case <-pc.done:
					fmt.Println("stoped")
					err := cmd.Process.Kill()
					if err != nil {
						logger.Errorf("error kill process: %v\n", err)
						return
					}
					return

				}
			}
		}()

		err := cmd.Start()
		if err != nil {
			logger.Errorf("error start command:%v\n", err)
		}

	case StopAudio:
		pc.done <- struct{}{}
		fmt.Println("stop audio")
		pc.sendStatus(StatusOK, addr)

	default:
		err := pc.sendStatus(StatusNotOK, addr)
		if err != nil {
			return err
		}
		return fmt.Errorf("invalid command")
	}
	return nil
}
