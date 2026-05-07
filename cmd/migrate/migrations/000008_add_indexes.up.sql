CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- GIN: array containment queries (e.g. WHERE tags @> '{go}')
CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING GIN(tags);

-- GIN + trgm: ILIKE/LIKE search trên title (e.g. WHERE title ILIKE '%golang%')
CREATE INDEX IF NOT EXISTS idx_posts_title_trgm ON posts USING GIN(title gin_trgm_ops);

-- B-tree: lookup posts theo user, comments theo post, followers theo follower
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_followers_follower_id ON followers(follower_id);
