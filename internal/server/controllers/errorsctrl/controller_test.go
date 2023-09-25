package errorsctrl

import (
	"context"
	"testing"
	"time"

	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
	"github.com/stretchr/testify/assert"
)

func TestController_Run(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "valid base case",
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch1 := make(chan msgbroker.Message, 3)
			ch2 := make(chan msgbroker.Message, 3)
			ch3 := make(chan msgbroker.Message, 3)
			defer close(ch1)
			defer close(ch2)
			defer close(ch3)
			chOut := make(chan msgbroker.Message, 3)

			c := &Controller{
				chansIn: []<-chan msgbroker.Message{ch1, ch2, ch3},
				chOut:   chOut,
				log:     log,
			}

			ctxWithCancel, cancel := context.WithCancel(tt.args.ctx)
			go c.Run(ctxWithCancel)

			ch1 <- msgbroker.Message{}
			ch2 <- msgbroker.Message{}
			ch3 <- msgbroker.Message{}

			var inCount int
			go func() {
				for range chOut {
					inCount++
				}
			}()

			waitChecks := make(chan struct{})
			go func() {
				_, closed := <-chOut
				assert.True(t, closed, "check out put channel is closed")
				waitChecks <- struct{}{}
			}()

			cancel()
			<-waitChecks
			time.Sleep(time.Second)
		})
	}
}

func TestController_RunStoppedFanInByChannelsClose(t *testing.T) {
	log, _ := logger.CreateZapLogger("info")
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "valid base case",
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch1 := make(chan msgbroker.Message, 3)
			ch2 := make(chan msgbroker.Message, 3)
			ch3 := make(chan msgbroker.Message, 3)
			chOut := make(chan msgbroker.Message, 3)

			c := &Controller{
				chansIn: []<-chan msgbroker.Message{ch1, ch2, ch3},
				chOut:   chOut,
				log:     log,
			}

			ctxWithCancel, _ := context.WithCancel(tt.args.ctx)
			go c.Run(ctxWithCancel)

			ch1 <- msgbroker.Message{}
			ch2 <- msgbroker.Message{}
			ch3 <- msgbroker.Message{}

			close(ch1)
			close(ch2)
			close(ch3)

			waitChecks := make(chan struct{})
			go func() {
				_, closed := <-chOut
				assert.True(t, closed, "check out put channel is closed")
				waitChecks <- struct{}{}
			}()

			<-waitChecks
		})
	}
}
