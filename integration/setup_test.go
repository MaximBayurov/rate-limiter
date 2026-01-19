package integration

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MaximBayurov/rate-limiter/internal/client"
	"github.com/MaximBayurov/rate-limiter/internal/configuration"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// TestSuite - базовый класс для всех интеграционных тестов.
type TestSuite struct {
	suite.Suite
	apiClient client.Client
	db        *sql.DB
}

// SetupSuite выполняется перед всеми тестами.
func (s *TestSuite) SetupSuite() {
	// Читаем конфигурацию из переменных окружения
	host := os.Getenv("API_HOST")
	if host == "" {
		host = "http://localhost"
	}
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "80"
	}

	DSN := os.Getenv("DATABASE_DSN")
	if DSN == "" {
		DSN = "postgres://postgres:postgres@localhost:5432/testing?sslmode=disable"
	}

	// Инициализируем клиент
	s.apiClient = client.New(configuration.ClientConf{
		Host:    host,
		Port:    port,
		Timeout: 100,
	})

	// Подключаемся к БД
	var err error
	s.db, err = sql.Open("postgres", DSN)
	require.NoError(s.T(), err)

	// Проверяем подключение
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = s.db.PingContext(ctx)
	require.NoError(s.T(), err)

	// Очищаем базу данных перед запуском тестов
	s.cleanDatabase()
}

// TearDownSuite выполняется после всех тестов.
func (s *TestSuite) TearDownSuite() {
	if s.db != nil {
		_ = s.db.Close()
	}
}

// SetupTest выполняется перед каждым тестом.
func (s *TestSuite) SetupTest() {
	// Очищаем базу данных перед каждым тестом
	s.cleanDatabase()
}

// cleanDatabase очищает базу данных.
func (s *TestSuite) cleanDatabase() {
	tables := []string{
		"ip_list",
	}

	for _, table := range tables {
		_, err := s.db.Exec(fmt.Sprintf("DELETE FROM %s", table)) //nolint:noctx
		if err != nil {
			// Таблица может не существовать, игнорируем ошибку
			log.Printf("Warning: не удалось очистить таблицу %s: %v", table, err)
		}
	}
}

func (s *TestSuite) tryAuthN(
	ctx context.Context,
	attempts int,
	login,
	password,
	ip string,
) (allowed int) {
	var resp client.Response
	for i := 0; i < attempts; i++ {
		resp, _ = s.apiClient.TryAuth(ctx, login, password, ip)
		if resp.Success {
			allowed++
		}
	}
	return allowed
}
