package httpctrl

import (
	"net/http"
	"strconv"

	"github.com/erupshis/effective_mobile/internal/datastructs"
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
	personData, err := httphelper.ParseQueryValues(r.URL.Query())
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

	err = c.strg.AddPersonData(r.Context(), personData)
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

	personData, err := httphelper.ParseQueryValues(r.URL.Query())
	if err != nil {
		c.log.Info("["+packageName+":Controller:updatePersonByIdHandler] failed to parse query params: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = httphelper.IsPersonDataValid(personData, true)
	if err != nil {
		c.log.Info("["+packageName+":Controller:updatePersonByIdHandler] person data is not valid: %v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = c.strg.UpdatePersonDataById(r.Context(), int64(id), personData)
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

	err = c.strg.UpdatePersonDataByIdPartially(r.Context(), int64(id), valuesToUpdate)
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

	//TODO: handling limit and offset.

	personsData, err := c.strg.GetPersonsData(r.Context(), valuesToFilter, 0, 0)
	if err != nil {
		c.log.Info("["+packageName+":Controller:getPersonsByFilterHandler] cannot process: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//TODO: write personsData in response body as Json
	personsData = append(personsData, datastructs.PersonData{})

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) badRequestHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
