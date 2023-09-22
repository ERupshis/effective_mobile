package extradatactrl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/erupshis/effective_mobile/internal/client"
	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
	"github.com/erupshis/effective_mobile/internal/server/helpers/msghelper"
)

const packageName = "extradatactrl"

var remoteServices = []string{
	"https://api.agify.io",
	"https://api.genderize.io",
	"https://api.nationalize.io",
}

type Controller struct {
	//INPUT channels
	chIn <-chan datastructs.ExtraDataFilling

	//OUTPUT channels
	chError chan<- msgbroker.Message
	chOut   chan<- datastructs.ExtraDataFilling

	client client.BaseClient
	log    logger.BaseLogger
}

func Create(chIn <-chan datastructs.ExtraDataFilling, chOut chan<- datastructs.ExtraDataFilling, chError chan<- msgbroker.Message,
	client client.BaseClient, log logger.BaseLogger) *Controller {
	return &Controller{
		chIn:    chIn,
		chOut:   chOut,
		chError: chError,
		client:  client,
		log:     log,
	}
}

func (c *Controller) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(c.chError)
			close(c.chOut)
			return

		case personDataIn, ok := <-c.chIn:
			if !ok {
				close(c.chError)
				close(c.chOut)
				return
			}

			if c.fillExtraData(ctx, &personDataIn) {
				c.log.Info("["+packageName+":Controller:Run] person data was filled with extra data: %v", personDataIn.Data)
				c.chOut <- personDataIn
			}
		}
	}
}

func (c *Controller) fillExtraData(ctx context.Context, dataIn *datastructs.ExtraDataFilling) bool {
	for _, service := range remoteServices {
		if err := c.doRemoteApiRequest(ctx, service, dataIn); err != nil {
			c.log.Info("["+packageName+":Controller:fillExtraData] request to remote API unsuccessful: %v", err)
			return false
		}
	}

	return true
}

func (c *Controller) doRemoteApiRequest(ctx context.Context, serviceUrl string, dataIn *datastructs.ExtraDataFilling) error {
	statusCode, body, err := c.client.DoGetURIWithQuery(ctx, serviceUrl, map[string]string{"name": dataIn.Data.Name})
	if err != nil {
		c.log.Info("["+packageName+":Controller:doRemoteApiRequest] failed to get person's age. source: %v statusCode: %d, body: %s, error: %v",
			dataIn.Data, statusCode, body, err)
	}

	if statusCode != http.StatusOK {
		errAPI := datastructs.Error{}
		err = json.Unmarshal(body, &errAPI)
		if err != nil {
			c.log.Info("JSON unmarshaling fail: %v", err)
			return fmt.Errorf("failed to parse response json: %w", err)
		}

		wrappedErr := fmt.Errorf("call remote API '%s' failed: %v", serviceUrl, errAPI.Data)
		if err = msghelper.PutErrorMessageInChan(c.chError, &dataIn.Raw, "error-remote-api", wrappedErr); err != nil {
			c.log.Info("["+packageName+":Controller:doRemoteApiRequest] failed to send error message. error: %v", err)
		}
		return fmt.Errorf("bad remote API response: %w", wrappedErr)
	}

	err = json.Unmarshal(body, &dataIn.Data)
	if err != nil {
		countries := datastructs.Countries{}
		err = json.Unmarshal(body, &countries)
		if err != nil {
			c.log.Info("["+packageName+":Controller:doRemoteApiRequest] JSON unmarshaling fail: %v", err)
			return err
		}

		dataIn.Data.Country = countries.Data[0].Id
	}

	return nil
}
