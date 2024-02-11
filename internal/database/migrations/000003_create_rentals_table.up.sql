CREATE TABLE IF NOT EXISTS rentals (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    bike_id INTEGER NOT NULL,
    status TEXT,
    start_time TIMESTAMP ,
    end_time TIMESTAMP,
    start_latitude REAL ,
    start_longitude REAL,
    end_latitude REAL,
    end_longitude REAL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    duration_minutes INTEGER DEFAULT 0,
    cost REAL NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(bike_id) REFERENCES bikes(id)
);