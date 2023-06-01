package sql

CREATE DATABASE wallet

CREATE TABLE wallets (
  id SERIAL PRIMARY KEY,
  dni VARCHAR(20) NOT NULL,
  country_id VARCHAR(50) NOT NULL,
  creation_date TIMESTAMP DEFAULT NOW(),
  balance DECIMAL(10, 2)
);

CREATE TABLE logs (
  id SERIAL PRIMARY KEY,
  dni VARCHAR(20) NOT NULL,
  stage VARCHAR(50) NOT NULL,
  creation_date TIMESTAMP DEFAULT NOW()
);

CREATE TABLE transactions (
    transactionID SERIAL PRIMARY KEY,
    senderID INT,
    receiverID INT,
    amount DECIMAL(10, 2),
    transactionType VARCHAR(10) CHECK (transactionType IN ('Extraction', 'Deposit')),
    transactionDate TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (senderID) REFERENCES wallets(id),
    FOREIGN KEY (receiverID) REFERENCES wallets(id)
);
