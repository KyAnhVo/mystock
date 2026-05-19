CREATE SCHEMA IF NOT EXISTS market_data;

CREATE TABLE IF NOT EXISTS market_data.ticker (
  id            BIGSERIAL   PRIMARY KEY,
  ticker        TEXT        NOT NULL UNIQUE,
  name          TEXT        NOT NULL,
  sector        TEXT,
  cik           INT,
  description   TEXT
);

CREATE TABLE IF NOT EXISTS market_data.stock (
  ticker_id     BIGINT,
  date          DATE,
  open          NUMERIC(12, 4)  NOT NULL,
  high          NUMERIC(12, 4)  NOT NULL,
  low           NUMERIC(12, 4)  NOT NULL,
  close         NUMERIC(12, 4)  NOT NULL,
  volume        BIGINT          NOT NULL,
  adj_close     NUMERIC(12, 4)  NOT NULL,

  PRIMARY KEY (ticker_id, date),
  CONSTRAINT fk_ticker  FOREIGN KEY (ticker_id) REFERENCES market_data.ticker(id)
);

CREATE TABLE IF NOT EXISTS market_data.index (
  id          BIGSERIAL   PRIMARY KEY,
  symbol      TEXT        NOT NULL UNIQUE,
  name        TEXT        NOT NULL,
  description TEXT
);

CREATE TABLE IF NOT EXISTS market_data.index_components (
  ticker_id   BIGINT,
  index_id    BIGINT,
  weight      NUMERIC(8, 6) NOT NULL,
  added_on    DATE,
  removed_on  DATE,

  PRIMARY KEY (index_id, ticker_id),
  CONSTRAINT indcomp_ticker_ref FOREIGN KEY (ticker_id) REFERENCES market_data.ticker(id),
  CONSTRAINT indcomp_index_ref  FOREIGN KEY (index_id)  REFERENCES market_data.ticker(id)
);
