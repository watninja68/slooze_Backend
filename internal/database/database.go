package database

import (
	"backend/internal/models"
	"backend/migrations"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pressly/goose/v3"
	"io/fs"
	"log"
	"os"
	"strconv"
	"time"
)

type Service interface {
	Health() map[string]string
	CreateOrderDB(ctx context.Context, userID, restaurantID int64, totalPrice float64) (int64, error)
	GetUserByID(ctx context.Context, id int) (models.User, error)
	ListRestaurantsDB(ctx context.Context, countryFilter *int64) ([]models.Restaurant, error)
	GetMenuItemsDB(ctx context.Context, restaurantID int64) ([]models.MenuItem, error)
	ListOrdersDB(ctx context.Context) ([]Order, error)
	GetOrderDB(ctx context.Context, id int64) (Order, error)
	ListPaymentMethodsDB(ctx context.Context, userID int64) ([]PaymentMethod, error)
	AddPaymentMethodDB(ctx context.Context, pm PaymentMethod) (int64, error)
	UpdatePaymentMethodDB(ctx context.Context, pm PaymentMethod) error
	Close() error
}
type service struct {
	db *sql.DB
}

var (
	database   = os.Getenv("BLUEPRINT_DB_DATABASE")
	password   = os.Getenv("BLUEPRINT_DB_PASSWORD")
	username   = os.Getenv("BLUEPRINT_DB_USERNAME")
	port       = os.Getenv("BLUEPRINT_DB_PORT")
	host       = os.Getenv("BLUEPRINT_DB_HOST")
	schema     = os.Getenv("BLUEPRINT_DB_SCHEMA")
	dbInstance *service
)

func MigrateFs(db *sql.DB, migrationFS fs.FS, dir string) error {
	// Ensure the custom filesystem is set for Goose
	goose.SetBaseFS(migrationFS)
	// Defer resetting to nil to avoid interfering with other potential Goose users
	defer func() {
		goose.SetBaseFS(nil) // Restore default filesystem interaction
	}()

	// Set the dialect before running migrations
	if err := goose.SetDialect("postgres"); err != nil {
		// Wrap the error for better context
		return fmt.Errorf("failed to set goose dialect to 'postgres': %w", err)
	}

	// Run the migrations
	log.Printf("Running database migrations from directory '%s' within embedded FS...", dir)
	// Use the Up function to apply all pending migrations
	if err := goose.Up(db, dir); err != nil {
		// Log the failure and attempt to get the status for debugging
		log.Printf("Goose 'up' migration failed: %v. Checking migration status...", err)
		// Use a separate function or inline the status check logic
		if statusErr := goose.Status(db, dir); statusErr != nil {
			log.Printf("Additionally failed to get goose migration status after 'up' failure: %v", statusErr)
		}
		// Return the original error from the 'Up' command
		return fmt.Errorf("goose 'up' migration failed: %w", err)
	}

	log.Println("Database migrations 'up' completed successfully.")
	return nil // Migrations applied successfully
}
func MigrateStatus(db *sql.DB, dir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}
	log.Println("Checking migration status...")
	if err := goose.Status(db, dir); err != nil {
		return fmt.Errorf("failed to get goose status: %w", err)
	}
	return nil
}
func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	db, err := sql.Open("pgx", connStr)

	if err := MigrateFs(db, migrations.FS, "."); err != nil {
		if statusErr := MigrateStatus(db, "."); statusErr != nil {
			log.Printf("Additionally failed to get migration status: %v", statusErr)
		}
		log.Panicf("Migration error during New(): %v", err)
	}
	log.Println("Database migrations applied successfully.")

	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func CreateOrderDB() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	db, err := sql.Open("pgx", connStr)

	if err := MigrateFs(db, migrations.FS, "."); err != nil {
		if statusErr := MigrateStatus(db, "."); statusErr != nil {
			log.Printf("Additionally failed to get migration status: %v", statusErr)
		}
		log.Panicf("Migration error during New(): %v", err)
	}
	log.Println("Database migrations applied successfully.")

	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}
func (s *service) GetUserByID(ctx context.Context, id int) (models.User, error) {
	const q = `
        SELECT u.id, u.name, u.email,
               r.id, r.name,
               c.id, c.name
        FROM users u
        JOIN roles      r ON r.id = u.role_id
        JOIN countries  c ON c.id = u.country_id
        WHERE u.id = $1
        LIMIT 1;
    `
	var out models.User
	if err := s.db.QueryRowContext(ctx, q, id).
		Scan(&out.ID, &out.Name, &out.Email,
			&out.Role.ID, &out.Role.Name,
			&out.Country.ID, &out.Country.Name); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return out, fmt.Errorf("user %d not found", id)
		}
		return out, fmt.Errorf("get user by id: %w", err)
	}
	return out, nil
}
