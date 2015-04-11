CREATE TABLE devices (
    id          INT     NOT NULL    PRIMARY KEY,
    name        TEXT    NOT NULL    UNIQUE,
    model       TEXT    NOT NULL,
    protocol    TEXT    NOT NULL    DEFAULT '',
    address     TEXT    NOT NULL    DEFAULT ''
);
