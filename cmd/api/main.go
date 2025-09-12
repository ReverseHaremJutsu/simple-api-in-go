package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/segmentio/kafka-go"

	accountappservice "rest-api-in-gin/internal/account/application/service"
	accountdomainservice "rest-api-in-gin/internal/account/domain/service"
	accountinfrarepo "rest-api-in-gin/internal/account/infrastructure/repository"
	accountpresenthttp "rest-api-in-gin/internal/account/presentation/http"

	paymentappservice "rest-api-in-gin/internal/payment/application/service"
	paymentdomainservice "rest-api-in-gin/internal/payment/domain/service"
	paymentinfraconsumer "rest-api-in-gin/internal/payment/infrastructure/kafka/consumer"
	paymentinfrarepo "rest-api-in-gin/internal/payment/infrastructure/repository"
	paymentpresenthttp "rest-api-in-gin/internal/payment/presentation/http"

	env "rest-api-in-gin/internal/env"
	infrakafkaproducer "rest-api-in-gin/internal/infrastructure/kafka/producer"
	persistence "rest-api-in-gin/internal/infrastructure/persistence"
)

func main() {
	// initialize database connection
	dsn := env.GetEnvString("DATABASE_URL", "postgres://jacky:password@localhost:5432/mydb?sslmode=disable")
	db, err := persistence.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	// initialize infrastructure repositories
	userAccountChecker := accountinfrarepo.NewPostgresUserAccountChecker(db)
	userAccountWriter := accountinfrarepo.NewPostgresUserAccountWriter(db)

	walletWriter := paymentinfrarepo.NewPostgresWalletWriter(db)
	walletReader := paymentinfrarepo.NewPostgresWalletReader(db)
	walletChecker := paymentinfrarepo.NewPostgresWalletChecker(db)

	// initialize domain services
	userAccountDomainService := accountdomainservice.NewUserAccountDomainService(userAccountChecker)
	walletDomainService := paymentdomainservice.NewWalletDomainService(walletChecker, walletReader, walletWriter)

	// initialize application services
	registerService := accountappservice.NewRegisterAccountService(userAccountWriter, userAccountDomainService)
	depositService := paymentappservice.NewDepositFundService(walletWriter, walletDomainService)

	// initialize presentation handlers for all application services in each module
	accountHandler := accountpresenthttp.NewAccountHandler(registerService)
	paymentHandler := paymentpresenthttp.NewPaymentHandler(depositService)

	// register all routes
	router := SetupRouter(accountHandler, paymentHandler)

	// initialize context
	port := env.GetEnvString("PORT", "8080")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// initialize context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // ensures resources are released

	// prepare kafka producers/consumers
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
	})
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "user-created",
		GroupID: "payment-module",
	})
	outboxDelivery := infrakafkaproducer.NewOutboxDelivery(db, kafkaWriter)
	userCreatedConsumer := paymentinfraconsumer.NewUserCreatedConsumer(kafkaReader, walletDomainService)

	// initialize server
	go func() {
		log.Println("HTTP server starting on port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// initialize kafka producers/consumers
	go outboxDelivery.Run(ctx)
	go userCreatedConsumer.Run(ctx)

	// wait for termination signal for kafka producers/consumers
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	log.Println("Shutting down...")

	// shutdown kafka producers/consumers
	cancel()

	// shutdown HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP server shutdown failed: %v", err)
	}
}
