package wrtc

import "github.com/FurmanovVitaliy/pixel-cloud/pkg/logger"

type StreamerService struct {
	videoCodec string
	audioCodec string
	logger     logger.Logger
}

func NewStreamerService(videoCodec, audioCodec string) *StreamerService {
	return &StreamerService{
		videoCodec: videoCodec,
		audioCodec: audioCodec,
		logger:     logger.Init("debug"),
	}
}

func (s *StreamerService) CreateRTC(onMessage func(data []byte)) (WebRTC, error) {
	rtc := newWRTC(onMessage, s.videoCodec, s.audioCodec, s.logger)
	return rtc, nil
}
