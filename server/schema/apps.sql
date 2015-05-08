CREATE TABLE apps (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT    NOT NULL    UNIQUE,
    protocol    TEXT    NOT NULL    DEFAULT '',
    address     TEXT    NOT NULL    DEFAULT ''
);
