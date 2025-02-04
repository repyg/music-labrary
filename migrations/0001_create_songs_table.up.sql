CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    group VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    release_date VARCHAR(255),
    text TEXT,
    link VARCHAR(255)
);
