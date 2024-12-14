CREATE TABLE IF NOT EXISTS roles(
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL UNIQUE,
    level int NOT NULL DEFAULT 0,
    description text
);

INSERT INTO
    roles (name, description, level)
VALUES
    ('user', 'A user can create posts and comments', 1),
    ('moderator', 'A user can update other users posts', 2),
    ('admin', 'An admin can update and delete other users posts', 3);