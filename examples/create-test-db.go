package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("Creating test SQLite database...")

	dbName := os.Getenv("TEST_DB_NAME")
	if dbName == "" {
		dbName = "test.db"
	}

	fmt.Printf("Database file: %s\n", dbName)
	os.Remove(dbName)

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	queries := []string{
		`CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL,
			first_name TEXT,
			last_name TEXT,
			phone_number TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			price DECIMAL(10,2),
			category_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			total_amount DECIMAL(10,2),
			status TEXT DEFAULT 'pending',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,

		`CREATE TABLE order_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_id INTEGER,
			product_id INTEGER,
			quantity INTEGER,
			price DECIMAL(10,2),
			FOREIGN KEY (order_id) REFERENCES orders(id),
			FOREIGN KEY (product_id) REFERENCES products(id)
		)`,

		`CREATE TABLE categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT
		)`,

		`CREATE TABLE user_profiles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			bio TEXT,
			website TEXT,
			social_media TEXT,
			preferences TEXT,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("Error creating table: %v", err)
		}
	}

	fmt.Println("Inserting sample data...")

	userInserts := []string{
		`INSERT INTO users (email, first_name, last_name, phone_number) VALUES 
		('john.doe@vizcore.com', 'John', 'Doe', '555-0101'),
		('jane.smith@vizcore.com', 'Jane', 'Smith', '555-0102'),
		('bob.johnson@vizcore.com', 'Bob', 'Johnson', '555-0103'),
		('alice.brown@vizcore.com', 'Alice', 'Brown', '555-0104'),
		('charlie.davis@vizcore.com', 'Charlie', 'Davis', '555-0105')`,
	}

	categoryInserts := []string{
		`INSERT INTO categories (name, description) VALUES 
		('Electronics', 'Electronic devices and gadgets'),
		('Books', 'Books and publications'),
		('Clothing', 'Apparel and accessories'),
		('Home & Garden', 'Home improvement and garden supplies')`,
	}

	productInserts := []string{
		`INSERT INTO products (name, description, price, category_id) VALUES 
		('Laptop Computer', 'High-performance laptop', 999.99, 1),
		('Programming Book', 'Learn Go programming', 49.99, 2),
		('T-Shirt', 'Comfortable cotton t-shirt', 19.99, 3),
		('Garden Hose', '50ft garden hose', 29.99, 4),
		('Smartphone', 'Latest model smartphone', 699.99, 1),
		('Novel', 'Bestselling fiction novel', 14.99, 2),
		('Jeans', 'Premium denim jeans', 79.99, 3),
		('Plant Pot', 'Ceramic plant pot', 12.99, 4)`,
	}

	orderInserts := []string{
		`INSERT INTO orders (user_id, total_amount, status) VALUES 
		(1, 1049.98, 'completed'),
		(2, 64.98, 'pending'),
		(3, 99.98, 'completed'),
		(1, 29.99, 'shipped')`,
	}

	orderItemInserts := []string{
		`INSERT INTO order_items (order_id, product_id, quantity, price) VALUES 
		(1, 1, 1, 999.99),
		(1, 5, 1, 49.99),
		(2, 2, 1, 49.99),
		(2, 4, 1, 14.99),
		(3, 3, 1, 19.99),
		(3, 7, 1, 79.99),
		(4, 4, 1, 29.99)`,
	}

	profileInserts := []string{
		`INSERT INTO user_profiles (user_id, bio, website, social_media, preferences) VALUES 
		(1, 'Software developer', 'https://paulmothapo.co.za', '@paulmothapo', 'tech,programming'),
		(2, NULL, NULL, '@janesmith', 'books,travel'),
		(3, 'Marketing specialist', NULL, NULL, NULL),
		(4, 'Designer and artist', 'https://mahlako.com', '@mahlako', 'art,design,photography'),
		(5, NULL, NULL, NULL, NULL)`,
	}

	allInserts := append(userInserts, categoryInserts...)
	allInserts = append(allInserts, productInserts...)
	allInserts = append(allInserts, orderInserts...)
	allInserts = append(allInserts, orderItemInserts...)
	allInserts = append(allInserts, profileInserts...)

	for _, insert := range allInserts {
		_, err := db.Exec(insert)
		if err != nil {
			log.Fatalf("Error inserting data: %v", err)
		}
	}

	indexes := []string{
		`CREATE INDEX idx_users_email ON users(email)`,
		`CREATE INDEX idx_products_category ON products(category_id)`,
		`CREATE INDEX idx_orders_user ON orders(user_id)`,
		`CREATE INDEX idx_orders_status ON orders(status)`,
	}

	for _, index := range indexes {
		_, err := db.Exec(index)
		if err != nil {
			log.Printf("Warning: Could not create index: %v", err)
		}
	}

	fmt.Printf("✓ Test database '%s' created successfully!\n", dbName)
	fmt.Println("✓ Sample data inserted")
	fmt.Println("✓ Indexes created")
	fmt.Println("")
	fmt.Println("You can now test the system with:")
	fmt.Printf("  $env:DATABASE_URL = \"./%s\"\n", dbName)
	fmt.Println("  go run examples/env-database-example.go")
	fmt.Println("Or:")
	fmt.Printf("  $env:DATABASE_URL = \"./%s\"\n", dbName)
	fmt.Println("  go run examples/real-database-example.go")
}
