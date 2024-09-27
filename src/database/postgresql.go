package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Luthor91/DiscordBot/database/migration"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB est la variable globale qui stocke la connexion à la base de données
var DB *gorm.DB

func InitDatabase() {
	validateEnvVars()
	CreatePostgresDatabase()
	DB = ConnectToPostgres()
	migration.MigrateAllPostgresql(DB)
}

func validateEnvVars() {
	requiredVars := []string{"DB_ADMIN_PASSWORD", "DB_ADMIN_USER", "DB_HOST", "DB_PORT", "DB_SSL_MODE", "DB_USER", "DB_PASSWORD", "DB_NAME"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			log.Fatalf("La variable d'environnement %s n'est pas définie", v)
		}
	}
}

func CreatePostgresDatabase() bool {
	dbName := os.Getenv("DB_NAME")
	password := os.Getenv("DB_ADMIN_PASSWORD")
	user := os.Getenv("DB_ADMIN_USER")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	sslMode := os.Getenv("DB_SSL_MODE")

	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=%s", user, password, host, port, sslMode)

	// Connexion en tant qu'administrateur pour vérifier et créer la base de données
	dbAdmin, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Erreur lors de la connexion en tant qu'administrateur: %v", err)
	}
	defer dbAdmin.Close()

	// Vérifier si la base de données existe déjà
	dbExists, err := CheckDatabaseExists(dbName)
	if err != nil {
		log.Fatalf("Erreur lors de la vérification de l'existence de la base de données: %v", err)
	}
	if dbExists {
		log.Printf("La base de données %s existe déjà.", dbName)
		return false
	}

	// Créer la base de données si elle n'existe pas
	stmt := fmt.Sprintf("CREATE DATABASE %s", dbName)
	_, err = dbAdmin.Exec(stmt)
	if err != nil {
		log.Printf("Erreur lors de la création de la base de données: %v", err)
		return false
	}

	log.Printf("La base de données %s a été créée avec succès.", dbName)
	return true
}

func CheckDatabaseExists(dbName string) (bool, error) {
	password := os.Getenv("DB_ADMIN_PASSWORD")
	user := os.Getenv("DB_ADMIN_USER")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	sslMode := os.Getenv("DB_SSL_MODE")

	// Connexion sans spécifier le nom de la base de données cible
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=%s", user, password, host, port, sslMode)

	dbAdmin, err := sql.Open("postgres", connStr)
	if err != nil {
		return false, fmt.Errorf("Erreur lors de la connexion pour vérifier l'existence de la base de données: %v", err)
	}
	defer dbAdmin.Close()

	var dbFound string
	err = dbAdmin.QueryRow("SELECT datname FROM pg_database WHERE datname = $1", dbName).Scan(&dbFound)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("Erreur lors de la vérification de l'existence de la base de données: %v", err)
	}

	return true, nil
}

// ConnectToPostgres initialise la connexion à la base de données PostgreSQL
func ConnectToPostgres() *gorm.DB {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSL_MODE")

	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s", user, password, host, port, dbName, sslMode)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erreur lors de la connexion à la base de données: %v", err)
	}

	return db
}
