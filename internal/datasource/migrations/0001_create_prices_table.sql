CREATE TABLE IF NOT EXISTS prices 
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    create_date DATE NOT NULL
);
