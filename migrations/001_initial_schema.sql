-- Create reviews table
CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    repository_name VARCHAR(255) NOT NULL,
    pull_request_id INTEGER NOT NULL,
    review_status VARCHAR(50) NOT NULL, -- 'pending', 'completed', 'failed'
    summary TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(repository_name, pull_request_id)
);

-- Create review_comments table
CREATE TABLE IF NOT EXISTS review_comments (
    id SERIAL PRIMARY KEY,
    review_id INTEGER NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    filename VARCHAR(500) NOT NULL,
    line_number INTEGER,
    comment_body TEXT NOT NULL,
    severity VARCHAR(20) NOT NULL, -- 'info', 'warning', 'error'
    posted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_reviews_repo_pr ON reviews(repository_name, pull_request_id);
CREATE INDEX IF NOT EXISTS idx_reviews_status ON reviews(review_status);
CREATE INDEX IF NOT EXISTS idx_review_comments_review_id ON review_comments(review_id);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for reviews table
DROP TRIGGER IF EXISTS update_reviews_updated_at ON reviews;
CREATE TRIGGER update_reviews_updated_at
    BEFORE UPDATE ON reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
