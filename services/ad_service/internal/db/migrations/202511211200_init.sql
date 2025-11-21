-- Initial schema for ad_service
CREATE TABLE IF NOT EXISTS ads (
    id UUID PRIMARY KEY,
    author_id UUID NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    price BIGINT NOT NULL,
    category_id UUID NOT NULL,
    condition TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'ACTIVE',
    seller_rating_cached DOUBLE PRECISION,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ad_images (
    id UUID PRIMARY KEY,
    ad_id UUID NOT NULL REFERENCES ads(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    position INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY,
    slug TEXT UNIQUE,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS favorites (
    user_id UUID NOT NULL,
    ad_id UUID NOT NULL REFERENCES ads(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, ad_id)
);

CREATE TABLE IF NOT EXISTS ad_reviews (
    id UUID PRIMARY KEY,
    ad_id UUID NOT NULL REFERENCES ads(id) ON DELETE CASCADE,
    reviewer_id UUID NOT NULL,
    rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_ads_category_price ON ads(category_id, price);
CREATE INDEX IF NOT EXISTS idx_ads_author ON ads(author_id);
CREATE INDEX IF NOT EXISTS idx_ad_images_ad_id ON ad_images(ad_id);
CREATE INDEX IF NOT EXISTS idx_ad_reviews_ad_id ON ad_reviews(ad_id);
-- Fulltext (will require extension)
-- CREATE EXTENSION IF NOT EXISTS pg_trgm;
-- CREATE INDEX IF NOT EXISTS idx_ads_fulltext ON ads USING GIN (to_tsvector('russian', title || ' ' || description));
