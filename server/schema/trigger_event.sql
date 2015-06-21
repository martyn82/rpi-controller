CREATE TABLE trigger_event (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    trigger_id      INTEGER NOT NULL    DEFAULT 0,
    agent_name      TEXT    NOT NULL    DEFAULT '',
    property_name   TEXT    NOT NULL    DEFAULT '',
    property_value  TEXT    NOT NULL    DEFAULT ''
);
