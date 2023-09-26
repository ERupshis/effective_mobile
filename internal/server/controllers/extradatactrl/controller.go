// Package extradatactr remote api calls controller.
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

// Controller struct of controller. Receives data which should be filled with additional data, taken from remote services.
// In case of success puts person data with extra data to further operating. Otherwise - generates error message in error channel.
type Controller struct {
	//INPUT channels
	chIn <-chan datastructs.ExtraDataFilling

	//OUTPUT channels
	chError chan<- msgbroker.Message
	chOut   chan<- datastructs.ExtraDataFilling

	client client.BaseClient
	log    logger.BaseLogger
}

// Create creates controller.
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

// Run main goroutine listens input channel and processes it.
func (c *Controller) Run(ctx context.Context) {
	c.log.Info("[%s:Controller:Run] start work", packageName)

	for {
		select {
		case <-ctx.Done():
			c.log.Info("[%s:Controller:Run] stop work due to context was canceled", packageName)
			close(c.chError)
			close(c.chOut)
			return

		case personDataIn, ok := <-c.chIn:
			if !ok {
				c.log.Info("[%s:Controller:Run] stop work due to input channel was closed", packageName)
				close(c.chError)
				close(c.chOut)
				return
			}

			if c.fillExtraData(ctx, &personDataIn) {
				c.log.Info("[%s:Controller:Run] person data was filled with extra data: %v", packageName, personDataIn.Data)
				c.chOut <- personDataIn
			}
		}
	}
}

// fillExtraData cycle for remote apis calls for additional person data.
func (c *Controller) fillExtraData(ctx context.Context, dataIn *datastructs.ExtraDataFilling) bool {
	for _, service := range remoteServices {
		if err := c.doRemoteApiRequest(ctx, service, dataIn); err != nil {
			c.log.Info("[%s:Controller:fillExtraData] request to remote API unsuccessful: %v", packageName, err)
			return false
		}
	}

	return true
}

// fillExtraData makes request to remote service and handling response.
func (c *Controller) doRemoteApiRequest(ctx context.Context, serviceUrl string, dataIn *datastructs.ExtraDataFilling) error {
	statusCode, body, err := c.client.DoGetURIWithQuery(ctx, serviceUrl, map[string]string{"name": dataIn.Data.Name})
	if err != nil {
		c.log.Info("[%s:Controller:doRemoteApiRequest] failed to get person's age. source: %v statusCode: %d, body: %s, error: %v",
			packageName, dataIn.Data, statusCode, body, err)
	}

	if statusCode != http.StatusOK {
		errAPI := datastructs.Error{}
		err = json.Unmarshal(body, &errAPI)
		if err != nil {
			c.log.Info("JSON unmarshalling fail: %v", err)
			return fmt.Errorf("failed to parse response json: %w", err)
		}

		wrappedErr := fmt.Errorf("call remote API '%s' failed: %v", serviceUrl, errAPI.Data)
		if err = msghelper.PutErrorMessageInChan(c.chError, &dataIn.Raw, "error-remote-api", wrappedErr); err != nil {
			c.log.Info("[%s:Controller:doRemoteApiRequest] failed to send error message. error: %v", packageName, err)
		}
		return fmt.Errorf("bad remote API response: %w", wrappedErr)
	}

	err = json.Unmarshal(body, &dataIn.Data)
	if err != nil {
		countries := datastructs.Countries{}
		err = json.Unmarshal(body, &countries)
		if err != nil {
			c.log.Info("[%s:Controller:doRemoteApiRequest] JSON unmarshalling fail: %v", packageName, err)
			return err
		}

		dataIn.Data.Country = countries.Data[0].Id
	}

	return nil
}
