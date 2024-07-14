CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL
    );

CREATE TABLE IF NOT EXISTS groups (
                                      id SERIAL PRIMARY KEY,
                                      name VARCHAR(50) NOT NULL,
    description TEXT
    );
