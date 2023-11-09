package v1

import (
	"backend/internal/api/interfaces/helper"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// ProductResources provide products interfaces
type ProductResources struct {
}

// NewProductResources returns productResources struct
func NewProductResources() *ProductResources {
	return &ProductResources{}
}

// Routes creates a REST router for the product resource
func (rs ProductResources) Routes() chi.Router {
	r := chi.NewRouter()
	r.With().Get("/", rs.list)
	r.With().Get("/{productId}", rs.getDetail)
	return r
}

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductDetail struct {
	Product
	Description string `json:"description"`
}

func (rs ProductResources) list(w http.ResponseWriter, r *http.Request) {
	products := []Product{
		{ID: "1111", Name: "Product 1", Price: 100.0},
		{ID: "2222", Name: "Product 2", Price: 200.0},
	}
	logrus.Debug(products)
	helper.Succeed(w, products)
}

func (rs ProductResources) getDetail(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productId")
	product := ProductDetail{
		Product: Product{
			ID:    productID,
			Name:  "Product " + productID,
			Price: 100.0,
		},
		Description: "This is a great collaboration product!",
	}
	logrus.Debug(product)
	helper.Succeed(w, product)
}
