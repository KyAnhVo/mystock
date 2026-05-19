CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE IF NOT EXISTS users.users (
  id                UUID    PRIMARY KEY,
  email             TEXT    UNIQUE NOT NULL,
  username          TEXT    UNIQUE NOT NULL,
  password_hashed   TEXT    NOT NULL
);

CREATE TABLE IF NOT EXISTS users.session (
  session_token     TEXT          PRIMARY KEY,
  user_id           UUID          NOT NULL,
  starts_at         TIMESTAMPTZ   NOT NULL, 
  expires_at        TIMESTAMPTZ   NOT NULL,

  CONSTRAINT fk_userid_session FOREIGN KEY (user_id) REFERENCES users.users(id)
);

