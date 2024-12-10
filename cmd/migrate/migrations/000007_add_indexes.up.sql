CREATE INDEX idx_comments_content ON comments (content);

CREATE INDEX IF NOT EXISTS idx_posts_title ON posts(title);

CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts (tags);

CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);