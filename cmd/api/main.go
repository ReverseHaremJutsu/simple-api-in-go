package main

import (
	"log"
	"net/http"

	accountappservice "rest-api-in-gin/internal/account/application/service"
	accountdomainservice "rest-api-in-gin/internal/account/domain/service"
	"rest-api-in-gin/internal/env"

	accountinfrapersistence "rest-api-in-gin/internal/account/infrastructure/persistence"
	accountinfrarepo "rest-api-in-gin/internal/account/infrastructure/repository"
	accountpresenthttp "rest-api-in-gin/internal/account/presentation/http"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// -----------------------------
	// 1️⃣ Initialize DB
	// -----------------------------
	dsn := env.GetEnvString("DATABASE_URL", "postgres://jacky:password@localhost:5432/mydb?sslmode=disable")
	db, err := accountinfrapersistence.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	// -----------------------------
	// 2️⃣ Initialize Repository
	// -----------------------------
	userAccountChecker := accountinfrarepo.NewPostgresUserAccountChecker(db)
	userAccountWriter := accountinfrarepo.NewPostgresUserAccountWriter(db)

	// -----------------------------
	// 3️⃣ Initialize Domain Service
	// -----------------------------
	userAccountDomainService := accountdomainservice.NewUserAccountDomainService(userAccountChecker)

	// -----------------------------
	// 4️⃣ Initialize Application Service
	// -----------------------------
	registerService := accountappservice.NewRegisterAccountService(userAccountWriter, userAccountDomainService)
	// other services like updateService, deleteService can be initialized here

	// -----------------------------
	// 5️⃣ Initialize Handlers
	// -----------------------------

	// accountHandler := accountpresenthttp.NewAccountHandler(registerService)
	accountHandler := accountpresenthttp.NewAccountHandler(registerService)

	// -----------------------------
	// 6️⃣ Setup router
	// -----------------------------
	router := SetupRouter(accountHandler)

	// // -----------------------------
	// // 6️⃣ Setup HTTP routes
	// // -----------------------------
	// accountpresenthttp.RegisterAccountRoutes(accountHandler)

	// -----------------------------
	// 7️⃣ Start server
	// -----------------------------
	port := env.GetEnvString("PORT", "8080")

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

// package main

// import (
// 	"database/sql"
// 	"log"
// 	_ "rest-api-in-gin/docs"
// 	"rest-api-in-gin/internal/database"
// 	"rest-api-in-gin/internal/env"

// 	_ "github.com/joho/godotenv/autoload"
// 	_ "github.com/mattn/go-sqlite3"
// )

// // @title Go Gin Rest API
// // @version 1.0
// // @description A rest API in Go using Gin framework.
// // @securityDefinitions.apikey BearerAuth
// // @in header
// // @name Authorization
// // @description Enter your bearer token in the format **Bearer &lt;token&gt;**

// // Apply the security definition to your endpoints
// // @security BearerAuth

// type application struct {
// 	port      int
// 	jwtSecret string
// 	models    database.Models
// }

// func main() {

// 	db, err := sql.Open("sqlite3", "./data.db")
// 	if err != nil {
// 		log.Fatal((err))
// 	}

// 	defer db.Close()

// 	models := database.NewModels(db)

// 	app := &application{
// 		port:      env.GetEnvInt("PORT", 8080),
// 		jwtSecret: env.GetEnvString("JWT_SECRET", "some-secret-123456"),
// 		models:    models,
// 	}

// 	if err := app.serve(); err != nil {
// 		log.Fatal(err)
// 	}
// }
