package events

import (
	"fmt"
	"time"

	"github.com/avast/retry-go"
)

var (
	DefaultAttempts  = uint(100)
	DefaultDelay     = 10 * time.Second
	DefaultDelayType = retry.FixedDelay
)

type RestartingSaramaConsumer struct {
	saramaConsumer *SaramaConsumer
	attempts       uint
	delay          time.Duration
	delayType      retry.DelayTypeFunc
}

func NewRestartingSaramaConsumer(saramaConsumer *SaramaConsumer) (*RestartingSaramaConsumer, error) {

	restartingSaramaConsumer := &RestartingSaramaConsumer{
		saramaConsumer: saramaConsumer,
		attempts:       DefaultAttempts,
		delay:          DefaultDelay,
		delayType:      DefaultDelayType,
	}

	return restartingSaramaConsumer, nil
}

func (s *RestartingSaramaConsumer) Start() error {
	if err := retry.Do(
		s.saramaConsumer.Start,
		retry.Attempts(s.attempts),
		retry.Delay(s.delay),
		retry.DelayType(s.delayType),
	); err != nil {
		fmt.Printf("Failed after %v retries to start the Sarama Consumer.", s.attempts)
		return err
	}

	return nil
}

func (s *RestartingSaramaConsumer) Stop() error {
	if err := s.saramaConsumer.Stop(); err != nil {
		return err
	}
	return nil
}
