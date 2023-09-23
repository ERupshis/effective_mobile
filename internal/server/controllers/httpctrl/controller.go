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
	"github.com/erupshis/effective_mobile/internal/server/storage/managers"
	"github.com/go-chi/chi/v5"
)

const packageName = "httpctrl"

type Controller struct {
	strg managers.BaseStorageManager
	log  logger.BaseLogger
}

func Create(strg managers.BaseStorageManager, log logger.BaseLogger) *Controller {
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
	r.Patch("/", c.updatePersonByIdPartiallyHandler)
	r.Get("/", c.getPersonsByFilterHandler)

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

	_, err = requestshelper.IsPersonDataValid(personData, true)
	if err != nil {
		c.log.Info("["+packageName+":Controller:createPersonHandler] data validation failed: %v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
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

	affectedCount, err := c.strg.DeletePersonById(r.Context(), int64(id))
	if err != nil {
		c.log.Info("["+packageName+":Controller:deletePersonByIdHandler] person id is not valid: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if affectedCount == 0 {
		c.log.Info("["+packageName+":Controller:deletePersonByIdHandler] request has no effect with id '%d'", id)
		w.WriteHeader(http.StatusNotModified)
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

	personData, err := requestshelper.ParsePersonDataFromJSON(buf.Bytes())
	if err != nil {
		c.log.Info("["+packageName+":Controller:updatePersonByIdHandler] failed to parse query params: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = requestshelper.IsPersonDataValid(personData, true)
	if err != nil {
		c.log.Info("["+packageName+":Controller:updatePersonByIdHandler] data validation failed: %v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	affectedCount, err := c.strg.UpdatePersonById(r.Context(), int64(id), personData)
	if err != nil {
		c.log.Info("["+packageName+":Controller:updatePersonByIdHandler] cannot process: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if affectedCount == 0 {
		c.log.Info("["+packageName+":Controller:updatePersonByIdHandler] request has no effect with id '%d'", id)
		w.WriteHeader(http.StatusNotModified)
		return
	}

	c.log.Info("["+packageName+":Controller:updatePersonByIdHandler] person with id '%d' successfully updated", id)
	responseBody := []byte("fully updated")
	w.Header().Add("Content-Length", strconv.FormatInt(int64(len(responseBody)), 10))
	w.Header().Add("Content-Type", "text/plain")
	_, _ = w.Write(responseBody)
}

func (c *Controller) updatePersonByIdPartiallyHandler(w http.ResponseWriter, r *http.Request) {
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
		c.log.Info("["+packageName+":Controller:updatePersonByIdPartiallyHandler] failed to read request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer helpers.ExecuteWithLogError(r.Body.Close, c.log)

	var valuesToUpdate map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &valuesToUpdate)
	if err != nil {
		c.log.Info("["+packageName+":Controller:updatePersonByIdPartiallyHandler] failed to parse query params: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	valuesToUpdate = requestshelper.FilterValues(valuesToUpdate)
	if len(valuesToUpdate) == 0 {
		c.log.Info("["+packageName+":Controller:updatePersonByIdPartiallyHandler] missing values to update in request '%v'", buf.Bytes())
		w.WriteHeader(http.StatusNotModified)
		return
	}

	affectedCount, err := c.strg.UpdatePersonByIdPartially(r.Context(), int64(id), valuesToUpdate)
	if err != nil {
		c.log.Info("["+packageName+":Controller:updatePersonByIdPartiallyHandler] cannot process: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if affectedCount == 0 {
		c.log.Info("["+packageName+":Controller:updatePersonByIdPartiallyHandler] request has no effect with id '%d'", id)
		w.WriteHeader(http.StatusNotModified)
		return
	}

	c.log.Info("["+packageName+":Controller:updatePersonByIdPartiallyHandler] person with id '%d' successfully updated", id)
	responseBody := []byte("partially updated")
	w.Header().Add("Content-Length", strconv.FormatInt(int64(len(responseBody)), 10))
	w.Header().Add("Content-Type", "text/plain")
	_, _ = w.Write(responseBody)
}

func (c *Controller) getPersonsByFilterHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	valuesToFilter, _ := requestshelper.ParseQueryValuesIntoMap(values)

	pageNum, pageSize := requestshelper.ParsePageAndPageSize(values)
	if pageSize < 0 {
		c.log.Info("[" + packageName + ":Controller:getPersonsByFilterHandler] negative page size")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if pageNum < 0 {
		c.log.Info("[" + packageName + ":Controller:getPersonsByFilterHandler] negative page num")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if len(values) != 0 {
		c.log.Info("["+packageName+":Controller:getPersonsByFilterHandler] unknown keys in request: %v", values)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	personsData, err := c.strg.SelectPersons(r.Context(), valuesToFilter, pageNum, pageSize)
	if err != nil {
		c.log.Info("["+packageName+":Controller:getPersonsByFilterHandler] cannot process: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var responseBody []byte
	if responseBody, err = json.MarshalIndent(personsData, "", "\t"); err != nil {
		c.log.Info("["+packageName+":Controller:getPersonsByFilterHandler] convert request result into JSON failed: %v", err)
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
