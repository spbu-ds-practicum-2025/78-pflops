package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"78-pflops/services/ad_service/internal/model"
)

type AdRepository struct {
	pool *pgxpool.Pool
}

func NewAdRepository(pool *pgxpool.Pool) *AdRepository {
	return &AdRepository{pool: pool}
}

func (r *AdRepository) Create(ctx context.Context, ad *model.Ad) error {
	if ad.ID == "" {
		ad.ID = uuid.New().String()
	}
	ad.CreatedAt = time.Now()
	ad.UpdatedAt = ad.CreatedAt
	_, err := r.pool.Exec(ctx, `INSERT INTO ads (id, author_id, title, description, price, category_id, condition, status, seller_rating_cached, created_at, updated_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		ad.ID, ad.AuthorID, ad.Title, ad.Description, ad.Price, ad.CategoryID, ad.Condition, ad.Status, ad.SellerRatingCached, ad.CreatedAt, ad.UpdatedAt,
	)
	return err
}

func (r *AdRepository) Get(ctx context.Context, id string) (*model.Ad, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, author_id, title, description, price, category_id, condition, status, seller_rating_cached, created_at, updated_at FROM ads WHERE id=$1`, id)
	var ad model.Ad
	var rating *float64
	if err := row.Scan(&ad.ID, &ad.AuthorID, &ad.Title, &ad.Description, &ad.Price, &ad.CategoryID, &ad.Condition, &ad.Status, &rating, &ad.CreatedAt, &ad.UpdatedAt); err != nil {
		return nil, err
	}
	ad.SellerRatingCached = rating
	return &ad, nil
}

func (r *AdRepository) ListImages(ctx context.Context, adID string) ([]model.AdImage, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, ad_id, url, is_primary, position FROM ad_images WHERE ad_id=$1 ORDER BY position ASC, id ASC`, adID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var images []model.AdImage
	for rows.Next() {
		var img model.AdImage
		if err := rows.Scan(&img.ID, &img.AdID, &img.URL, &img.IsPrimary, &img.Position); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}

func (r *AdRepository) Search(ctx context.Context, text string, categoryID *string, priceMin, priceMax *int64, condition *string, limit, offset int) ([]model.Ad, int, error) {
	// Simplified search (will extend later with proper builder)
	query := `SELECT id, author_id, title, description, price, category_id, condition, status, seller_rating_cached, created_at, updated_at FROM ads WHERE 1=1`
	args := []any{}
	idx := 1
	appendCond := func(cond string, val any) {
		query += fmt.Sprintf(" AND %s $%d", cond, idx)
		args = append(args, val)
		idx++
	}
	if text != "" {
		query += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", idx, idx)
		args = append(args, "%"+text+"%")
		idx++
	}
	if categoryID != nil {
		appendCond("category_id =", *categoryID)
	}
	if priceMin != nil {
		appendCond("price >=", *priceMin)
	}
	if priceMax != nil {
		appendCond("price <=", *priceMax)
	}
	if condition != nil {
		appendCond("condition =", *condition)
	}
	// Pagination
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", idx, idx+1)
	args = append(args, limit, offset)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []model.Ad
	for rows.Next() {
		var ad model.Ad
		var rating *float64
		if err := rows.Scan(&ad.ID, &ad.AuthorID, &ad.Title, &ad.Description, &ad.Price, &ad.CategoryID, &ad.Condition, &ad.Status, &rating, &ad.CreatedAt, &ad.UpdatedAt); err != nil {
			return nil, 0, err
		}
		ad.SellerRatingCached = rating
		list = append(list, ad)
	}
	return list, len(list), nil
}

func (r *AdRepository) Update(ctx context.Context, id string, authorID string, title, description *string, price *int64) error {
	set := "updated_at = NOW()"
	args := []any{}
	idx := 1
	add := func(fragment string, val any) {
		set += ", " + fragment + fmt.Sprintf(" $%d", idx)
		args = append(args, val)
		idx++
	}
	if title != nil {
		add("title =", *title)
	}
	if description != nil {
		add("description =", *description)
	}
	if price != nil {
		add("price =", *price)
	}
	// WHERE id and author
	query := fmt.Sprintf("UPDATE ads SET %s WHERE id = $%d AND author_id = $%d", set, idx, idx+1)
	args = append(args, id, authorID)
	res, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("not found or no permission")
	}
	return nil
}

func (r *AdRepository) AttachMedia(ctx context.Context, adID, mediaID string) error {
	// store mediaID as URL for simplicity
	id := uuid.New().String()
	_, err := r.pool.Exec(ctx, `INSERT INTO ad_images (id, ad_id, url, is_primary, position) VALUES ($1,$2,$3,false,0)`, id, adID, mediaID)
	return err
}

func (r *AdRepository) Delete(ctx context.Context, id string, authorID string) error {
	res, err := r.pool.Exec(ctx, `DELETE FROM ads WHERE id=$1 AND author_id=$2`, id, authorID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("not found or no permission")
	}
	return nil
}
