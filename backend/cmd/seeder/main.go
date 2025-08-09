package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"procurement-system/internal/models"
	"procurement-system/internal/repository"
	"procurement-system/internal/services"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	db                 *sql.DB
	userRepo           repository.UserRepository
	vendorRepo         repository.VendorRepository
	requisitionRepo    repository.RequisitionRepository
	poRepo             repository.PurchaseOrderRepository
	requisitionService services.RequisitionService
)

func main() {
	// Load .env file from the current or parent directories
	err := godotenv.Load()
	if err != nil {
		// If .env is not in the current dir, try the parent (backend/)
		if err = godotenv.Load("../.env"); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	// Connect to database
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	// Initialize repositories
	userRepo = repository.NewPostgresUserRepository(db)
	vendorRepo = repository.NewPostgresVendorRepository(db)
	requisitionRepo = repository.NewPostgresRequisitionRepository(db)
	poRepo = repository.NewPostgresPurchaseOrderRepository(db)

	// Initialize services (we need this for the approval logic)
	activityLogRepo := repository.NewPostgresActivityLogRepository(db)
	logService := services.NewActivityLogService(activityLogRepo)
	pdfService := services.NewPDFService()
	poService := services.NewPurchaseOrderService(poRepo, vendorRepo, pdfService)
	requisitionService = services.NewRequisitionService(requisitionRepo, poService, logService)

	fmt.Println("Starting database seeding...")

	// Clean existing data to avoid duplicates on re-run
	cleanData()

	// Seed data
	users := seedUsers()
	vendors := seedVendors()
	seedRequisitions(users, vendors)

	fmt.Println("Database seeding completed successfully!")
}

func cleanData() {
	fmt.Println("Cleaning existing data...")
	// Order is important due to foreign key constraints
	if _, err := db.Exec("DELETE FROM purchase_orders;"); err != nil {
		log.Printf("Warn: could not delete from purchase_orders: %v", err)
	}
	if _, err := db.Exec("DELETE FROM requisitions;"); err != nil {
		log.Printf("Warn: could not delete from requisitions: %v", err)
	}
	if _, err := db.Exec("DELETE FROM vendors;"); err != nil {
		log.Printf("Warn: could not delete from vendors: %v", err)
	}
	if _, err := db.Exec("DELETE FROM users;"); err != nil {
		log.Printf("Warn: could not delete from users: %v", err)
	}

	// Reset sequences
	fmt.Println("Resetting table sequences...")
	db.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1;")
	db.Exec("ALTER SEQUENCE vendors_id_seq RESTART WITH 1;")
	db.Exec("ALTER SEQUENCE requisitions_id_seq RESTART WITH 1;")
	db.Exec("ALTER SEQUENCE purchase_orders_id_seq RESTART WITH 1;")
	fmt.Println("Data cleaned.")
}

func seedUsers() []models.User {
	fmt.Println("Seeding users...")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	usersToCreate := []models.User{
		{Name: "Admin User", Email: "admin@example.com", HashedPassword: string(hashedPassword), Role: "Admin"},
		{Name: "Employee One", Email: "employee1@example.com", HashedPassword: string(hashedPassword), Role: "Employee"},
		{Name: "Employee Two", Email: "employee2@example.com", HashedPassword: string(hashedPassword), Role: "Employee"},
		{Name: "Procurement Officer", Email: "procurement@example.com", HashedPassword: string(hashedPassword), Role: "Procurement Officer"},
		{Name: "Approver One", Email: "approver@example.com", HashedPassword: string(hashedPassword), Role: "Approver"},
	}

	var createdUsers []models.User
	for _, user := range usersToCreate {
		createdUser, err := userRepo.CreateUser(&user)
		if err != nil {
			log.Fatalf("Error creating user %s: %v", user.Name, err)
		}
		fmt.Printf("Created user: %s (ID: %d)\n", createdUser.Name, createdUser.ID)
		createdUsers = append(createdUsers, *createdUser)
	}
	return createdUsers
}

func seedVendors() []models.Vendor {
	fmt.Println("Seeding vendors...")
	vendorsToCreate := []models.Vendor{
		{Name: "Tech Supplies Inc.", ContactPerson: "John Smith", Email: "contact@techsupplies.com", Phone: "123-456-7890", Address: "123 Tech Park, Silicon Valley, CA"},
		{Name: "Office Furniture Co.", ContactPerson: "Jane Doe", Email: "sales@officefurn.com", Phone: "098-765-4321", Address: "456 Business Rd, New York, NY"},
		{Name: "Global IT Solutions", ContactPerson: "Peter Jones", Email: "info@globalit.com", Phone: "555-555-5555", Address: "789 Enterprise Way, London, UK"},
	}

	var createdVendors []models.Vendor
	for _, vendor := range vendorsToCreate {
		err := vendorRepo.CreateVendor(&vendor)
		if err != nil {
			log.Fatalf("Error creating vendor %s: %v", vendor.Name, err)
		}
		fmt.Printf("Created vendor: %s (ID: %d)\n", vendor.Name, vendor.ID)
		createdVendors = append(createdVendors, vendor)
	}
	return createdVendors
}

func seedRequisitions(users []models.User, vendors []models.Vendor) {
	fmt.Println("Seeding requisitions...")
	employee1 := users[1]
	employee2 := users[2]
	vendor1 := vendors[0]
	vendor2 := vendors[1]

	requisitionsToCreate := []models.Requisition{
		{RequesterID: employee1.ID, VendorID: &vendor1.ID, ItemDescription: "10x New Dell Laptops", Quantity: 10, EstimatedPrice: 1200.00, TotalPrice: 12000.00, Justification: "New hire setup", Status: "Pending"},
		{RequesterID: employee2.ID, VendorID: &vendor2.ID, ItemDescription: "5x Ergonomic Office Chairs", Quantity: 5, EstimatedPrice: 350.00, TotalPrice: 1750.00, Justification: "Replace old chairs", Status: "Pending"},
		{RequesterID: employee1.ID, VendorID: &vendor2.ID, ItemDescription: "20x Standing Desks", Quantity: 20, EstimatedPrice: 500.00, TotalPrice: 10000.00, Justification: "Office wellness initiative", Status: "Pending"},
		{RequesterID: employee2.ID, VendorID: &vendor1.ID, ItemDescription: "1x VR Headset", Quantity: 1, EstimatedPrice: 800.00, TotalPrice: 800.00, Justification: "Research and development", Status: "Rejected"},
	}

	for i, req := range requisitionsToCreate {
		createdReq, err := requisitionRepo.CreateRequisition(&req)
		if err != nil {
			log.Fatalf("Error creating requisition: %v", err)
		}
		fmt.Printf("Created requisition for: %s (ID: %d)\n", createdReq.ItemDescription, createdReq.ID)

		// Approve the third requisition to trigger PO creation
		if i == 2 {
			fmt.Printf("Approving requisition ID %d to generate a Purchase Order...\n", createdReq.ID)
			// We'll use the Admin User (ID: 1) to approve this in the seeder
			adminID := 1
			err := requisitionService.ApproveRequisition(createdReq.ID, adminID)
			if err != nil {
				log.Fatalf("Error approving requisition: %v", err)
			}
			fmt.Println("Requisition approved and PO created.")
		}
	}
}
