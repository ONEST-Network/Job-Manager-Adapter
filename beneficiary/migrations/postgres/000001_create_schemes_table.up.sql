CREATE TABLE IF NOT EXISTS schemes (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    provider VARCHAR(100) NOT NULL,
    criteria TEXT[],
    amount DECIMAL(12,2) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_schemes_provider ON schemes(provider);
CREATE INDEX idx_schemes_status ON schemes(status);
CREATE INDEX idx_schemes_amount ON schemes(amount); 