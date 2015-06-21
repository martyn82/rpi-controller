CREATE TABLE trigger_action (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    trigger_id      INTEGER NOT NULL    DEFAULT 0,
    agentName       TEXT    NOT NULL    DEFAULT '',
    propertyName    TEXT    NOT NULL    DEFAULT '',
    propertyValue   TEXT    NOT NULL    DEFAULT ''
);
