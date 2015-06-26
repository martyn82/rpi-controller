CREATE TABLE devices (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT    NOT NULL    UNIQUE,
    model       TEXT    NOT NULL,
    protocol    TEXT    NOT NULL    DEFAULT '',
    address     TEXT    NOT NULL    DEFAULT '',
    extra       TEXT    NOT NULL    DEFAULT ''
);
