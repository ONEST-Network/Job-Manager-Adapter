CREATE TABLE IF NOT EXISTS applications (
    id VARCHAR(50) PRIMARY KEY,
    scheme_id VARCHAR(50) NOT NULL REFERENCES schemes(id),
    applicant_id VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    credentials JSONB NOT NULL,
    submitted_at TIMESTAMP WITH TIME ZONE NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_applications_scheme_id ON applications(scheme_id);
CREATE INDEX idx_applications_applicant_id ON applications(applicant_id);
CREATE INDEX idx_applications_status ON applications(status); 