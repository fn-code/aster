package aster

import (
	"net"

	"github.com/google/logger"
)

// Listen running audio server
func Listen(addr string) (*AudioServer, error) {
	lis, err := net.ListenPacket("udp", addr)
	if err != nil {
		logger.Fatalf("error resolve addr : %v\n", err)
		return nil, err
	}
	return &AudioServer{
		buffer: make([]byte, 1024),
		pc:     lis,
		done:   make(chan struct{}),
	}, nil
}

// ReadDataPlay is read data from client,
// using ffplay to play audio
func (pc *AudioServer) ReadDataPlay() {
	for {
		_, addr, err := pc.pc.ReadFrom(pc.buffer)
		if err != nil {
			logger.Infof("error read data: %v", err)
		}
		ok, err := pc.readHeader(addr)
		if err != nil {
			logger.Errorf("error header : %v", err)
		}
		if ok {
			err = pc.procesPlay(addr)
			if err != nil {
				logger.Errorf("error process audio : %v", err)
			}
		}

	}
}

// ReadDataStream is read data from client
// to stream audio to reeiver
func (pc *AudioServer) ReadDataStream() {
	for {
		_, addr, err := pc.pc.ReadFrom(pc.buffer)
		if err != nil {
			logger.Infof("error read data: %v", err)
		}
		ok, err := pc.readHeader(addr)
		if err != nil {
			logger.Errorf("error header : %v", err)
		}
		if ok {
			pc.procesStream(addr)
		}

	}
}

// Close is to close audio server
func (pc *AudioServer) Close() error {
	if err := pc.pc.Close(); err != nil {
		return err
	}
	return nil
}
