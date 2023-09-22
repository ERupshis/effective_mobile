package httpctrl

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/erupshis/effective_mobile/internal/helpers"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/server/helpers/httphelper"
	"github.com/erupshis/effective_mobile/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

const packageName = "httpctrl"

type Controller struct {
	strg storage.BaseStorageManager
	log  logger.BaseLogger
}

func Create(strg storage.BaseStorageManager, log logger.BaseLogger) *Controller {
	return &Controller{
		strg: strg,
		log:  log,
	}
}

func (c *Controller) Route() *chi.Mux {
	r := chi.NewRouter()

	r.Use(c.log.LogHandler)

	r.Post("/", c.createPersonHandler)
	r.Delete("/{id}", c.deletePersonByIdHandler)
	r.Put("/{id}", c.updatePersonByIdHandler)
	r.Patch("/{id}", c.updatePersonByIdPartiallyHandler)
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

	personData, err := httphelper.ParsePersonDataFromJSON(buf.Bytes())
	if err != nil {
		c.log.Info("["+packageName+":Controller:createPersonHandler] failed to parse query params: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = httphelper.IsPersonDataValid(personData, true)
	if err != nil {
		c.log.Info("["+packageName+":Controller:createPersonHandler] data validation: %v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = c.strg.AddPerson(r.Context(), personData)
	if err != nil {
		c.log.Info("["+packageName+":Controller:createPersonHandler] cannot process: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) deletePersonByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = c.strg.DeletePersonDataById(r.Context(), int64(id))
	if err != nil {
		c.log.Info("["+packageName+":Controller:deletePersonByIdHandler] person id is not valid: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) updatePersonByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
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

	personData, err := httphelper.ParsePersonDataFromJSON(buf.Bytes())
	if err != nil {
		c.log.Info("["+packageName+":Controller:updatePersonByIdHandler] failed to parse query params: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = httphelper.IsPersonDataValid(personData, true)
	if err != nil {
		c.log.Info("["+packageName+":Controller:updatePersonByIdHandler] data validation: %v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = c.strg.UpdatePersonById(r.Context(), int64(id), personData)
	if err != nil {
		c.log.Info("["+packageName+":Controller:updatePersonByIdHandler] cannot process: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) updatePersonByIdPartiallyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	valuesToUpdate, err := httphelper.ParseQueryValuesIntoMap(r.URL.Query())
	if err != nil {
		c.log.Info("["+packageName+":Controller:updatePersonByIdHandler] failed to parse query params: %v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = c.strg.UpdatePersonByIdPartially(r.Context(), int64(id), valuesToUpdate)
	if err != nil {
		c.log.Info("["+packageName+":Controller:updatePersonByIdPartiallyHandler] cannot process: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) getPersonsByFilterHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	valuesToFilter, _ := httphelper.ParseQueryValuesIntoMap(values)

	page, pageSize := httphelper.ParsePageAndPageSize(values)

	personsData, err := c.strg.GetPersons(r.Context(), valuesToFilter, page, pageSize)
	if err != nil {
		c.log.Info("["+packageName+":Controller:getPersonsByFilterHandler] cannot process: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var responseBody []byte
	if responseBody, err = json.Marshal(personsData); err != nil {
		c.log.Info("["+packageName+":Controller:getPersonsByFilterHandler] convert request result into JSON failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(responseBody)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) badRequestHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
