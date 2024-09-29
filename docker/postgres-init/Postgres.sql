CREATE TABLE IF NOT EXISTS pack_size (
    id SERIAL PRIMARY KEY,
    "size" INTEGER NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "order" (
    id SERIAL PRIMARY KEY,
    "size" INTEGER NOT NULL,
    packs JSON NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT NOW(),
    modified TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted BOOLEAN NOT NULL DEFAULT FALSE
);

-- Default Pack Sizes
INSERT INTO pack_size ("size") VALUES
    (250),
    (500),
    (1000),
    (2000),
    (5000)
;