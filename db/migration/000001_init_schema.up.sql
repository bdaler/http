CREATE TABLE banners (
    id BIGSERIAL primary KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    button TEXT NOT NULL,
    link TEXT NOT NULL,
    image TEXT NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
