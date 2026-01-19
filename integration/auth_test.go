package integration

import (
	"context"
	"fmt"
	"github.com/MaximBayurov/rate-limiter/internal/client"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestAuthAPI struct {
	TestSuite
}

func TestAuth(t *testing.T) {
	if testing.Short() {
		t.Skip("Пропускаем интеграционные тесты в режиме short")
	}

	suite.Run(t, new(TestAuthAPI))
}

func (s *TestAuthAPI) TestTryAuth() {
	ctx := context.Background()

	s.T().Run("Превышение rate limit по login", func(t *testing.T) {
		attempts := 15
		allowed := s.tryAuthN(ctx, attempts, "login101", "password101", "120.0.0.1")

		assert.Equal(t, 10, allowed)
	})
	s.T().Run("Превышение rate limit по password", func(t *testing.T) {
		ip := "120.0.0.2"
		pass := "password102"
		allowed := 0
		attempts := 105

		var r client.Response
		var login string
		for i := 0; i < attempts; i++ {
			login = fmt.Sprintf("pass-check-%d", i)
			r, _ = s.apiClient.TryAuth(ctx, login, pass, ip)
			if r.Success {
				allowed += 1
			}
		}
		assert.Equal(t, 100, allowed)
	})
	s.T().Run("Превышение rate limit по IP", func(t *testing.T) {
		ip := "120.0.0.3"
		var allowed int32 = 0
		wg := sync.WaitGroup{}
		for i := 0; i < 32; i++ {
			wg.Add(1)
			i := i
			go func() {
				defer wg.Done()
				var r client.Response
				var cred string
				for j := 0; j < 32; j++ {
					cred = fmt.Sprintf("%d-ip-check-%d", i, j)
					r, _ = s.apiClient.TryAuth(ctx, cred, cred, ip)
					if r.Success {
						atomic.AddInt32(&allowed, 1)
					}
				}
			}()
		}

		wg.Wait()
		assert.Equal(t, int32(1000), allowed)
	})
}

func (s *TestAuthAPI) TestClearBucket() {
	ctx := context.Background()

	s.T().Run("Успешный сброс bucket", func(t *testing.T) {
		ip := "120.0.0.4"
		password := "pass"
		login := "log"

		allowed := s.tryAuthN(ctx, 11, login, password, ip)
		assert.Equal(t, 10, allowed)

		r, _ := s.apiClient.ClearBucket(ctx, ip, login)
		assert.True(t, r.Success)

		allowed = s.tryAuthN(ctx, 11, login, password, ip)
		assert.Equal(t, 10, allowed)
	})
}
