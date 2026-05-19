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
  starts_at         TIMESTAMPTZ   NOT NULL DEFAULT NOW(), 
  expires_at        TIMESTAMPTZ   NOT NULL DEFAULT NOW() + INTERVAL '7 days',

  CONSTRAINT fk_userid_session FOREIGN KEY (user_id) REFERENCES users.users(id)
);

CREATE TABLE IF NOT EXISTS users.watchlist (
  user_id     UUID,
  ticker_id   BIGINT,

  PRIMARY KEY (user_id, ticker_id),
  CONSTRAINT fk_watchlist_user    FOREIGN KEY (user_id)   REFERENCES users.users(id),
  CONSTRAINT fk_watchlist_ticker  FOREIGN KEY (ticker_id) REFERENCES market_data.ticker(id)
);
