CREATE SCHEMA user;

CREATE TABLE user.user (
  id                UUID    PRIMARY KEY,
  email             TEXT    UNIQUE NOT NULL,
  username          TEXT    UNIQUE NOT NULL,
  password_hashed   TEXT    NOT NULL
);

CREATE TABLE user.session (
  session_token     TEXT          PRIMARY KEY,
  user_id            TEXT          NOT NULL,
  starts_at         TIMESTAMPTZ   NOT NULL, 
  expires_at        TIMESTAMPTZ   NOT NULL,

  CONSTRAINT fk_userid_session FOREIGN KEY (userid) REFERENCES user.user(id)
);

