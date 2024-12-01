CREATE TABLE IF NOT EXISTS currency (
    id SERIAL PRIMARY KEY,
    time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    currency VARCHAR(3) NOT NULL,
    type VARCHAR(4) NOT NULL,
    value DECIMAL(10,2) NOT NULL,
    
    CONSTRAINT currency_type_check CHECK (type IN ('buy', 'sell')),
    CONSTRAINT currency_code_check CHECK (currency IN ('USD', 'EUR'))
);

CREATE INDEX idx_currency_time ON currency(time);

CREATE INDEX idx_currency_composite ON currency(currency, type, time);

