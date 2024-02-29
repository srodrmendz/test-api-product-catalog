package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	internalError "github.com/srodrmendz/api-product-catalog/errors"
	"github.com/srodrmendz/api-product-catalog/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
)

// Create new product repository
func New(client *mongo.Client, database string, collection string) *ProductsCatalogRepository {
	return &ProductsCatalogRepository{
		collection: client.Database(database).Collection(collection),
	}
}

// Create a new product
func (r *ProductsCatalogRepository) Create(ctx context.Context, product model.Product) (*model.Product, error) {
	product.ID = uuid.NewString()

	now := time.Now()

	product.CreatedAt = now

	product.UpdatedAt = now

	product.InStock = product.Qty > 0

	// Create user on repository
	if _, err := r.collection.InsertOne(ctx, product); err != nil {
		writeError, ok := err.(mongo.WriteException)
		if !ok {
			return nil, fmt.Errorf("creating product %s on repository %w", product.Name, err)
		}

		for _, wErr := range writeError.WriteErrors {
			if wErr.Code == mongoDBDuplicatedKeyErrorCode {
				return nil, internalError.ErrProductSKUAlreadyExist
			}
		}

		return nil, fmt.Errorf("creating product %s on repository %w", product.Name, err)
	}

	return &product, nil
}

// Get a product by id
func (r *ProductsCatalogRepository) GetByID(ctx context.Context, id string) (*model.Product, error) {
	resp := r.collection.FindOne(ctx, bson.M{"_id": id})

	if resp.Err() != nil {
		return nil, internalError.ErrProductNotFound
	}

	var product model.Product

	if err := resp.Decode(&product); err != nil {
		return nil, fmt.Errorf("decoding product from repository %w", err)
	}

	return &product, nil
}

// Get a product by sku
func (r *ProductsCatalogRepository) GetBySKU(ctx context.Context, sku string) (*model.Product, error) {
	resp := r.collection.FindOne(ctx, bson.M{"sku": sku})

	if resp.Err() != nil {
		return nil, internalError.ErrProductNotFound
	}

	var product model.Product

	if err := resp.Decode(&product); err != nil {
		return nil, fmt.Errorf("decoding product from repository %w", err)
	}

	return &product, nil
}

// Delete product
func (r *ProductsCatalogRepository) Delete(ctx context.Context, id string) error {
	if _, err := r.collection.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return fmt.Errorf("deleting product from repository %w", err)
	}

	return nil
}

// Update a product
func (r *ProductsCatalogRepository) Update(ctx context.Context, id string, qty uint64) (*model.Product, error) {
	filter := bson.M{"_id": id}

	update := bson.M{
		"$set": bson.M{
			"qty":        qty,
			"updated_at": time.Now(),
			"in_stock":   qty > 0,
		},
	}

	var product model.Product

	err := r.
		collection.
		FindOneAndUpdate(
			ctx,
			filter,
			update,
			options.FindOneAndUpdate().SetReturnDocument(options.After)).
		Decode(&product)
	if err != nil {
		return nil, internalError.ErrProductNotFound
	}

	return &product, nil
}

// Search products
func (r *ProductsCatalogRepository) Search(ctx context.Context, limit int, offset int) ([]model.Product, *int64, error) {
	eg, _ := errgroup.WithContext(ctx)

	var products []model.Product

	var total int64

	// Call concurrenlty search method
	eg.Go(func() error {
		prds, err := r.search(ctx, limit, offset)
		if err != nil {
			return err
		}

		products = prds

		return nil
	})

	// Call concurrenlty get total items on db
	eg.Go(func() error {
		tl, err := r.getTotal(ctx)
		if err != nil {
			return err
		}

		total = tl

		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, nil, err
	}

	return products, &total, nil
}

func (r *ProductsCatalogRepository) search(ctx context.Context, limit int, offset int) ([]model.Product, error) {
	opt := options.Find()

	opt.SetLimit(int64(limit))

	opt.SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, bson.D{}, opt)
	if err != nil {
		return nil, err
	}

	var products []model.Product

	for cursor.Next(ctx) {
		var product model.Product

		if err := cursor.Decode(&product); err != nil {
			return nil, fmt.Errorf("decoding product from repository %w", err)
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *ProductsCatalogRepository) getTotal(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}
