CREATE TABLE IF NOT EXISTS usdt_rates
(
	id SERIAL PRIMARY KEY,
	timestamp TIMESTAMP NOT NULL,
	ask TEXT NOT NULL,
	bid TEXT NOT NULL
);