package suite

import (
	"context"
	"database/sql"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"catalog-microservice/internal/app"
	"catalog-microservice/internal/config"
	catalogpb "catalog-microservice/internal/pb/catalog"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const grpcHost = "localhost"

type Suite struct {
	*testing.T
	Cfg           *config.Config
	CatalogClient catalogpb.CatalogServiceClient
	DB            *sql.DB
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config/local_test.yaml")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	t.Cleanup(cancel)

	connStr := mustStartPostgres(ctx, t)
	mustRunMigrations(t, connStr)

	// находим свободный порт чтобы параллельные тесты не конфликтовали
	port := mustGetFreePort(t)
	cfg.GRPC.Port = port

	mustStartServer(t, cfg, connStr)

	// открываем отдельное соединение для хелперов вставки данных
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err, "failed to open db for helpers")
	t.Cleanup(func() { db.Close() })

	cc, err := grpc.NewClient(
		net.JoinHostPort(grpcHost, strconv.Itoa(port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err, "grpc connection failed")
	t.Cleanup(func() { cc.Close() })

	return ctx, &Suite{
		T:             t,
		Cfg:           cfg,
		CatalogClient: catalogpb.NewCatalogServiceClient(cc),
		DB:            db, // теперь передаём DB
	}
}

// mustGetFreePort просит ОС выдать свободный порт
func mustGetFreePort(t *testing.T) int {
	t.Helper()

	// слушаем на порту 0 — ОС сама выберет свободный
	ln, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err, "failed to get free port")
	defer ln.Close()

	return ln.Addr().(*net.TCPAddr).Port
}

func mustStartPostgres(ctx context.Context, t *testing.T) string {
	t.Helper()

	container, err := tcpostgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		tcpostgres.WithDatabase("catalog_test"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err, "failed to start postgres container")

	t.Cleanup(func() {
		if err := container.Terminate(context.Background()); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	})

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	return connStr
}

func mustRunMigrations(t *testing.T, connStr string) {
	t.Helper()

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(currentFile), "..", "..")
	migrationsPath := filepath.Join(projectRoot, "migrations")

	err = goose.Up(db, migrationsPath)
	require.NoError(t, err, "migrations failed")
}

func mustStartServer(t *testing.T, cfg *config.Config, connStr string) {
	t.Helper()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	application := app.NewWithDSN(log, cfg.GRPC.Port, connStr)

	go func() {
		application.GRPCSrv.MustRun()
	}()

	waitForPort(t, cfg.GRPC.Port)

	t.Cleanup(func() {
		application.GRPCSrv.Stop()
	})
}

func waitForPort(t *testing.T, port int) {
	t.Helper()

	addr := net.JoinHostPort(grpcHost, strconv.Itoa(port))

	require.Eventually(t,
		func() bool {
			conn, err := net.DialTimeout("tcp", addr, time.Second)
			if err != nil {
				return false
			}
			conn.Close()
			return true
		},
		10*time.Second,
		100*time.Millisecond,
		"grpc server did not start on port %d", port,
	)
}
