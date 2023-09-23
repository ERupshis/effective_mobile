package httpctrl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/erupshis/effective_mobile/internal/helpers"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/server/helpers/requestshelper"
	"github.com/erupshis/effective_mobile/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

const packageName = "httpctrl"

type Controller struct {
	strg storage.BaseStorage
	log  logger.BaseLogger
}

func Create(strg storage.BaseStorage, log logger.BaseLogger) *Controller {
	return &Controller{
		strg: strg,
		log:  log,
	}
}

func (c *Controller) Route() *chi.Mux {
	r := chi.NewRouter()

	r.Use(c.log.LogHandler)

	r.Post("/", c.createPersonHandler)
	r.Delete("/", c.deletePersonByIdHandler)
	r.Put("/", c.updatePersonByIdHandler)
	r.Patch("/", c.updatePersonByIdHandler)
	r.Get("/", c.selectPersonsByFilterHandler)

	r.NotFound(c.badRequestHandler)

	return r
}

func (c *Controller) createPersonHandler(w http.ResponseWriter, r *http.Request) {
	buf := bytes.Buffer{}
	if _, err := buf.ReadFrom(r.Body); err != nil {
		c.log.Info("["+packageName+":Controller:createPersonHandler] failed to read request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer helpers.ExecuteWithLogError(r.Body.Close, c.log)

	if len(buf.Bytes()) == 0 {
		c.log.Info("[" + packageName + ":Controller:createPersonHandler] empty request's body")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	personData, err := requestshelper.ParsePersonDataFromJSON(buf.Bytes())
	if err != nil {
		c.log.Info("["+packageName+":Controller:createPersonHandler] failed to parse query params: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newPersonId, err := c.strg.AddPerson(r.Context(), personData)
	if err != nil {
		c.log.Info("["+packageName+":Controller:createPersonHandler] cannot process: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.log.Info("["+packageName+":Controller:createPersonHandler] person with id '%d' successfully added", newPersonId)
	responseBody := []byte(fmt.Sprintf("Added with id: %d", newPersonId))
	w.Header().Add("Content-Length", strconv.FormatInt(int64(len(responseBody)), 10))
	w.Header().Add("Content-Type", "text/plain")
	_, _ = w.Write(responseBody)
}

func (c *Controller) deletePersonByIdHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if values.Get("id") == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(values.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	_, err = c.strg.DeletePersonById(r.Context(), int64(id))
	if err != nil {
		c.log.Info("["+packageName+":Controller:deletePersonByIdHandler] delete person: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.log.Info("["+packageName+":Controller:deletePersonByIdHandler] person with id '%d' successfully deleted", id)
	responseBody := []byte("successfully deleted")
	w.Header().Add("Content-Length", strconv.FormatInt(int64(len(responseBody)), 10))
	w.Header().Add("Content-Type", "text/plain")
	_, _ = w.Write(responseBody)
}

func (c *Controller) updatePersonByIdHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if values.Get("id") == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(values.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	buf := bytes.Buffer{}
	if _, err := buf.ReadFrom(r.Body); err != nil {
		c.log.Info("["+packageName+":Controller:updatePersonByIdHandler] failed to read request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer helpers.ExecuteWithLogError(r.Body.Close, c.log)

	var valuesToUpdate map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &valuesToUpdate)
	if err != nil {
		c.log.Info("["+packageName+":Controller:updatePersonByIdHandler] failed to parse query params: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = c.strg.UpdatePersonById(r.Context(), int64(id), valuesToUpdate)
	if err != nil {
		c.log.Info("["+packageName+":Controller:updatePersonByIdHandler] cannot process: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.log.Info("["+packageName+":Controller:updatePersonByIdHandler] person with id '%d' successfully updated", id)
	responseBody := []byte(fmt.Sprintf("person with id '%d' updated", id))
	w.Header().Add("Content-Length", strconv.FormatInt(int64(len(responseBody)), 10))
	w.Header().Add("Content-Type", "text/plain")
	_, _ = w.Write(responseBody)
}

func (c *Controller) selectPersonsByFilterHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	valuesToFilter, err := requestshelper.ParseQueryValuesIntoMap(values)
	if err != nil {
		c.log.Info("["+packageName+":Controller:selectPersonsByFilterHandler] cannot process: %v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	personsData, err := c.strg.SelectPersons(r.Context(), valuesToFilter)
	if err != nil {
		c.log.Info("["+packageName+":Controller:selectPersonsByFilterHandler] cannot process: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var responseBody []byte
	if responseBody, err = json.MarshalIndent(personsData, "", "\t"); err != nil {
		c.log.Info("["+packageName+":Controller:selectPersonsByFilterHandler] convert request result into JSON failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Length", strconv.FormatInt(int64(len(responseBody)), 10))
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(responseBody)
}

func (c *Controller) badRequestHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
