CREATE SCHEMA IF NOT EXISTS market_data;

CREATE TABLE IF NOT EXISTS market_data.ticker (
  id                BIGSERIAL   PRIMARY KEY,
  ticker            TEXT        NOT NULL UNIQUE,
  cik               TEXT,

  name              TEXT        NOT NULL,
  exchange          TEXT,
  description       TEXT
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
