package server

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	internalErrors "github.com/srodrmendz/api-product-catalog/errors"
	"github.com/srodrmendz/api-product-catalog/model"
	"github.com/srodrmendz/api-product-catalog/utils"
)

// Healthcheck godoc
// @Tags healthcheck
// @Accept  json
// @Produce  json
// @Success 200
// @Router /health-check [get]
func (a *App) healthCheck(w http.ResponseWriter, _ *http.Request) {
	response := map[string]string{
		"version":      a.Config.Version,
		"build_date":   a.Config.BuildDate,
		"service_name": "api-product-catalog",
	}

	utils.DataJSON(w, http.StatusOK, response)
}

// Create godoc
// @Tags create
// @Description Create product
// @Accept  json
// @Produce  json
// @Param request body model.Product true "Request body"
// @Success 201 {object} model.Product
// @Failure 500
// @Router /v1 [post]
func (a *App) create(w http.ResponseWriter, r *http.Request) {
	builder := model.NewCreateBuilder(r)

	request, err := builder.Build()
	if err != nil {
		utils.ErrJSON(w, http.StatusBadRequest, err)

		return
	}

	product, err := a.Services.ProductsService.Create(r.Context(), *request)
	if err != nil {
		if errors.Is(err, internalErrors.ErrProductSKUAlreadyExist) {
			utils.ErrJSON(w, http.StatusBadRequest, err)

			return
		}

		utils.ErrJSON(w, http.StatusInternalServerError, err)

		return
	}

	utils.DataJSON(w, http.StatusCreated, product)
}

// Get By ID godoc
// @Tags get by id
// @Description Get product by id
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} model.Product
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /v1/{id}/ [get]
func (a *App) getByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		utils.ErrJSON(w, http.StatusBadRequest, errors.New("id must be provided"))

		return
	}

	product, err := a.Services.ProductsService.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, internalErrors.ErrProductNotFound) {
			utils.ErrJSON(w, http.StatusNotFound, err)

			return
		}

		utils.ErrJSON(w, http.StatusInternalServerError, err)

		return
	}

	utils.DataJSON(w, http.StatusOK, product)
}

// Get By sku godoc
// @Tags get by sku
// @Description Get product by sku
// @Accept  json
// @Produce  json
// @Param sku path string true "sku"
// @Success 200 {object} model.Product
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /v1/sku/{id}/ [get]
func (a *App) getBySKU(w http.ResponseWriter, r *http.Request) {
	sku := mux.Vars(r)["sku"]

	if sku == "" {
		utils.ErrJSON(w, http.StatusBadRequest, errors.New("id must be provided"))

		return
	}

	product, err := a.Services.ProductsService.GetBySKU(r.Context(), sku)
	if err != nil {
		if errors.Is(err, internalErrors.ErrProductNotFound) {
			utils.ErrJSON(w, http.StatusNotFound, err)

			return
		}

		utils.ErrJSON(w, http.StatusInternalServerError, err)

		return
	}

	utils.DataJSON(w, http.StatusOK, product)
}

// Delete godoc
// @Tags delete
// @Description Delete product
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 204
// @Failure 400
// @Failure 500
// @Router /v1/{id}/ [delete]
func (a *App) delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		utils.ErrJSON(w, http.StatusBadRequest, errors.New("id must be provided"))

		return
	}

	if err := a.Services.ProductsService.Delete(r.Context(), id); err != nil {
		utils.ErrJSON(w, http.StatusInternalServerError, err)

		return
	}

	utils.DataJSON(w, http.StatusNoContent, nil)
}

// Update godoc
// @Tags update
// @Description Update product
// @Accept  json
// @Produce  json
// @Param request body model.UpdateRequest true "Request body"
// @Param id path string true "id"
// @Success 201 {object} model.Product
// @Failure 404
// @Failure 500
// @Router /v1/{id}/ [put]
func (a *App) update(w http.ResponseWriter, r *http.Request) {
	builder := model.NewUpdateBuilder(r)

	request, err := builder.Build()
	if err != nil {
		utils.ErrJSON(w, http.StatusBadRequest, err)

		return
	}

	product, err := a.Services.ProductsService.Update(r.Context(), request)
	if err != nil {
		if errors.Is(err, internalErrors.ErrProductNotFound) {
			utils.ErrJSON(w, http.StatusNotFound, err)

			return
		}

		utils.ErrJSON(w, http.StatusInternalServerError, err)

		return
	}

	utils.DataJSON(w, http.StatusCreated, product)
}

// Search godoc
// @Tags search
// @Description Search products
// @Accept  json
// @Produce  json
// @Param name query string false "name"
// @Param sort query string false "sort"
// @Param in_stock query string false "in stock"
// @Param limit query int true "limit"
// @Param offset query int true "offset"
// @Success 200 {object} model.SearchResponse
// @Failure 500
// @Router /v1 [get]
func (a *App) search(w http.ResponseWriter, r *http.Request) {
	builder := model.NewSearchBuilder(r)

	request, err := builder.Build()
	if err != nil {
		utils.ErrJSON(w, http.StatusBadRequest, err)

		return
	}

	products, err := a.Services.ProductsService.Search(r.Context(), *request)
	if err != nil {
		utils.ErrJSON(w, http.StatusInternalServerError, err)

		return
	}

	utils.DataJSON(w, http.StatusOK, products)
}
