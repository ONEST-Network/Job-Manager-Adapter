package db

import (
	"database/sql"
	"fmt"
	"time"
	"github.com/VinVorteX/beneficiary-manager/internal/db/migrate"
	"github.com/VinVorteX/beneficiary-manager/internal/logger"
	"github.com/VinVorteX/beneficiary-manager/internal/models"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	db     *sql.DB
	logger logger.Logger
}

func NewPostgresDB(connStr, migrationsPath string, logger logger.Logger) (*PostgresDB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Error("Failed to connect to database", err, nil)
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err = db.Ping(); err != nil {
		logger.Error("Failed to ping database", err, nil)
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	logger.Info("Successfully connected to database", nil)

	// Run migrations
	if err := migrate.MigrateDB(db, migrationsPath); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	return &PostgresDB{db: db, logger: logger}, nil
}

func (p *PostgresDB) GetSchemes(filter models.SchemeFilter) ([]models.Scheme, error) {
	p.logger.Debug("Executing GetSchemes query", map[string]interface{}{
		"filter": filter,
	})

	query := `SELECT id, name, description, provider, criteria, amount, start_date, end_date, status 
              FROM schemes 
              WHERE ($1 = '' OR provider = $1)
              AND ($2 = 0 OR amount >= $2)
              AND ($3 = 0 OR amount <= $3)
              AND ($4 = '' OR status = $4)`

	rows, err := p.db.Query(query, filter.Provider, filter.MinAmount, filter.MaxAmount, filter.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to query schemes: %v", err)
	}
	defer rows.Close()

	var schemes []models.Scheme
	for rows.Next() {
		var s models.Scheme
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.Provider, &s.Criteria,
			&s.Amount, &s.StartDate, &s.EndDate, &s.Status); err != nil {
			return nil, fmt.Errorf("failed to scan scheme row: %v", err)
		}
		schemes = append(schemes, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating scheme rows: %v", err)
	}

	return schemes, nil
}

// Add implementation for GetSchemeByID
func (p *PostgresDB) GetSchemeByID(id string) (*models.Scheme, error) {
	query := `SELECT id, name, description, provider, criteria, amount, start_date, end_date, status 
              FROM schemes 
              WHERE id = $1`

	var scheme models.Scheme
	err := p.db.QueryRow(query, id).Scan(
		&scheme.ID, &scheme.Name, &scheme.Description, &scheme.Provider,
		&scheme.Criteria, &scheme.Amount, &scheme.StartDate, &scheme.EndDate, &scheme.Status,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get scheme: %v", err)
	}

	return &scheme, nil
}

// Add implementation for SaveApplication
func (p *PostgresDB) SaveApplication(app models.Application) error {
	query := `INSERT INTO applications (id, scheme_id, applicant_id, status, credentials, submitted_at, last_updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := p.db.Exec(query,
		app.ID, app.SchemeID, app.ApplicantID, app.Status,
		app.Credentials, app.SubmittedAt, app.LastUpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save application: %v", err)
	}

	return nil
}

// Add implementation for GetApplicationStatus
func (p *PostgresDB) GetApplicationStatus(applicationID string) (*models.ApplicationStatus, error) {
	query := `SELECT id, status, last_updated_at 
              FROM applications 
              WHERE id = $1`

	var status models.ApplicationStatus
	err := p.db.QueryRow(query, applicationID).Scan(
		&status.ApplicationID, &status.Status, &status.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get application status: %v", err)
	}

	return &status, nil
}

// ConfigurePool sets up the database connection pool
func (p *PostgresDB) ConfigurePool(maxOpen, maxIdle int) error {
	p.db.SetMaxOpenConns(maxOpen)
	p.db.SetMaxIdleConns(maxIdle)
	p.db.SetConnMaxLifetime(time.Hour)

	p.logger.Info("Configured database connection pool", map[string]interface{}{
		"max_open_conns": maxOpen,
		"max_idle_conns": maxIdle,
	})
	return nil
}

// Add other necessary database methods
