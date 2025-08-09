package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"procurement-system/internal/handlers"
	"procurement-system/internal/middleware"
	"procurement-system/internal/repository"
	"procurement-system/internal/services"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get database connection string and JWT secret from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}

	// Connect to the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}
	log.Println("Successfully connected to the database")

	// Initialize repositories
	userRepo := repository.NewPostgresUserRepository(db)
	vendorRepo := repository.NewPostgresVendorRepository(db)
	requisitionRepo := repository.NewPostgresRequisitionRepository(db)
	poRepo := repository.NewPostgresPurchaseOrderRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo)
	vendorService := services.NewVendorService(vendorRepo)
	pdfService := services.NewPDFService()
	poService := services.NewPurchaseOrderService(poRepo, vendorRepo, pdfService)
	requisitionService := services.NewRequisitionService(requisitionRepo, poService)
	navigationService := services.NewNavigationService()
	userService := services.NewUserService(userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	vendorHandler := handlers.NewVendorHandler(vendorService)
	requisitionHandler := handlers.NewRequisitionHandler(requisitionService)
	poHandler := handlers.NewPurchaseOrderHandler(poService)
	navigationHandler := handlers.NewNavigationHandler(navigationService)
	userHandler := handlers.NewUserHandler(userService)

	// Create router
	r := mux.NewRouter()

	// Setup routes
	api := r.PathPrefix("/api").Subrouter()

	// Auth routes
	api.HandleFunc("/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/login", authHandler.Login).Methods("POST")

	// User Management routes (Admin only)
	userRoutes := api.PathPrefix("/users").Subrouter()
	userRoutes.Use(middleware.AuthMiddleware, middleware.RoleMiddleware("Admin"))
	userRoutes.HandleFunc("", userHandler.GetAllUsers).Methods("GET")
	userRoutes.HandleFunc("/{id:[0-9]+}", userHandler.GetUserByID).Methods("GET")
	userRoutes.HandleFunc("/{id:[0-9]+}", userHandler.UpdateUser).Methods("PUT")
	userRoutes.HandleFunc("/{id:[0-9]+}", userHandler.DeleteUser).Methods("DELETE")

	// Navigation routes
	navRoutes := api.PathPrefix("/navigation").Subrouter()
	navRoutes.Use(middleware.AuthMiddleware)
	navRoutes.HandleFunc("/menu", navigationHandler.GetMenu).Methods("GET")
	navRoutes.HandleFunc("/breadcrumbs", navigationHandler.GetBreadcrumbs).Methods("GET")

	// Vendor routes (Admin only)
	vendorRoutes := api.PathPrefix("/vendors").Subrouter()
	vendorRoutes.Use(middleware.AuthMiddleware, middleware.RoleMiddleware("Admin"))
	vendorRoutes.HandleFunc("", vendorHandler.CreateVendor).Methods("POST")
	vendorRoutes.HandleFunc("", vendorHandler.GetAllVendors).Methods("GET")
	vendorRoutes.HandleFunc("/{id:[0-9]+}", vendorHandler.GetVendorByID).Methods("GET")
	vendorRoutes.HandleFunc("/{id:[0-9]+}", vendorHandler.UpdateVendor).Methods("PUT")
	vendorRoutes.HandleFunc("/{id:[0-9]+}", vendorHandler.DeleteVendor).Methods("DELETE")

	// Requisition routes
	reqRoutes := api.PathPrefix("/requisitions").Subrouter()
	reqRoutes.Use(middleware.AuthMiddleware) // All requisition routes require authentication
	reqRoutes.HandleFunc("", requisitionHandler.CreateRequisition).Methods("POST")
	reqRoutes.HandleFunc("/my", requisitionHandler.GetMyRequisitions).Methods("GET")
	reqRoutes.HandleFunc("/{id:[0-9]+}", requisitionHandler.UpdateRequisition).Methods("PUT")
	reqRoutes.HandleFunc("/{id:[0-9]+}", requisitionHandler.DeleteRequisition).Methods("DELETE")

	// Admin-only requisition routes
	adminReqRoutes := reqRoutes.PathPrefix("").Subrouter()
	adminReqRoutes.Use(middleware.RoleMiddleware("Admin"))
	adminReqRoutes.HandleFunc("/pending", requisitionHandler.GetPendingRequisitions).Methods("GET")
	adminReqRoutes.HandleFunc("/all", requisitionHandler.GetAllRequisitions).Methods("GET")
	adminReqRoutes.HandleFunc("/{id:[0-9]+}/approve", requisitionHandler.ApproveRequisition).Methods("POST")
	adminReqRoutes.HandleFunc("/{id:[0-9]+}/reject", requisitionHandler.RejectRequisition).Methods("POST")

	// Purchase Order routes
	poRoutes := api.PathPrefix("/purchase-orders").Subrouter()
	poRoutes.Use(middleware.AuthMiddleware)
	poRoutes.HandleFunc("/{id:[0-9]+}", poHandler.GetPurchaseOrderByID).Methods("GET")
	poRoutes.HandleFunc("/{id:[0-9]+}/pdf", poHandler.GetPurchaseOrderPDF).Methods("GET")

	// Admin-only PO routes
	adminPoRoutes := poRoutes.PathPrefix("").Subrouter()
	adminPoRoutes.Use(middleware.RoleMiddleware("Admin"))
	adminPoRoutes.HandleFunc("/all", poHandler.GetAllPurchaseOrders).Methods("GET")

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins for development
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
