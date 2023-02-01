CREATE TABLE users
(
    id       SERIAL PRIMARY KEY,
    login    VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(50)        NOT NULL
);

CREATE TABLE balances
(
    id      SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users (id) NOT NULL,
    balance DECIMAL                       NOT NULL
        CONSTRAINT positive_balance CHECK (balance >= 0)
);

CREATE TABLE transactions
(
    id       SERIAL PRIMARY KEY,
    user_id  INTEGER REFERENCES users (id) NOT NULL,
    attempts INTEGER                       NOT NULL
        CONSTRAINT check_attempts CHECK (0 <= attempts AND attempts <= 5),
    status   INTEGER                       NOT NULL,
    type     INTEGER                       NOT NULL,
    amount   DECIMAL                       NOT NULL
        CONSTRAINT positive_balance CHECK (amount > 0),
    date     TIMESTAMP                     NOT NULL DEFAULT now()
);