package msggenerator

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
)

func Run(ctx context.Context, producer msgbroker.Producer, log logger.BaseLogger) {
	chGeneratedNames := make(chan []byte)
	go func(ctx context.Context, out <-chan []byte) {
		for {
			select {
			case <-ctx.Done():
				close(chGeneratedNames)
				return
			default:
				chGeneratedNames <- getRandomName()
				time.Sleep(3 * time.Second)
			}
		}
	}(ctx, chGeneratedNames)

	key := "person data"
	for randomName := range chGeneratedNames {
		err := producer.SendMessage(key, string(randomName))
		if err != nil {
			log.Info("send message failed: %v", err)
		}

		log.Info("Message sent: Key=%s, Value=%s\n", key, randomName)
	}
}

func getRandomName() []byte {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomSex := rand.Intn(2)
	isMale := randomSex == 0

	randomNationality := rand.Intn(2)
	isRussian := randomNationality == 0

	msg := msgbroker.MessageBody{}
	if isRussian {
		if isMale {
			msg.Name = getRandomValueFromSlice(russianFirstNamesMale)
			msg.Surname = getRandomValueFromSlice(russianLastNames)
			msg.Patronymic = getRandomValueFromSlice(russianPatronymicNamesMale)
		} else {
			msg.Name = getRandomValueFromSlice(russianFirstNamesFemale)
			msg.Surname = getRandomValueFromSlice(russianLastNames) + "a"
			msg.Patronymic = getRandomValueFromSlice(russianPatronymicNamesFemale)
		}
	} else {
		if isMale {
			msg.Name = getRandomValueFromSlice(englishFirstNamesMale)
		} else {
			msg.Name = getRandomValueFromSlice(englishFirstNamesFemale)
		}

		msg.Surname = getRandomValueFromSlice(englishLastNames)
	}

	msgJSON, _ := json.Marshal(msg)

	return msgJSON
}

func getRandomValueFromSlice(values []string) string {
	randomIdx := rand.Intn(len(values))
	return values[randomIdx]
}
