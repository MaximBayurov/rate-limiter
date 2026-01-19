package integration

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestIPListAPI struct {
	TestSuite
}

func TestIPList(t *testing.T) {
	if testing.Short() {
		t.Skip("Пропускаем интеграционные тесты в режиме short")
	}

	suite.Run(t, new(TestIPListAPI))
}

func (s *TestIPListAPI) TestAdd() {
	ctx := context.Background()
	s.T().Run("Успешное добавление в белый список IP", func(t *testing.T) {
		ip := "127.0.0.1"
		listType := "white"
		resp, err := s.apiClient.AddIP(ctx, ip, listType, false)
		assert.NoError(t, err)
		assert.True(t, resp.Success)

		attempts := 15
		allowed := s.tryAuthN(ctx, attempts, "login", "password", ip)
		assert.Equal(t, attempts, allowed)
	})
	s.T().Run("Успешное добавление в белый список IP с перезаписью", func(t *testing.T) {
		ip := "127.0.0.2"
		resp, err := s.apiClient.AddIP(ctx, ip, "black", false)
		assert.NoError(t, err)
		assert.True(t, resp.Success)

		resp, err = s.apiClient.AddIP(ctx, ip, "white", true)
		assert.NoError(t, err)
		assert.True(t, resp.Success)

		attempts := 15
		allowed := s.tryAuthN(ctx, attempts, "login1", "password1", ip)
		assert.Equal(t, attempts, allowed)
	})
	s.T().Run("Неуспешное добавление в белый список IP", func(t *testing.T) {
		ip := "127.0.0.3"
		resp, err := s.apiClient.AddIP(ctx, ip, "black", false)
		assert.NoError(t, err)
		assert.True(t, resp.Success)

		resp, err = s.apiClient.AddIP(ctx, ip, "white", false)
		assert.NoError(t, err)
		assert.False(t, resp.Success)

		allowed := s.tryAuthN(ctx, 1, "login2", "password2", ip)
		assert.Equal(t, 0, allowed)
	})
	s.T().Run("Добавление подсетей в белый список IP", func(t *testing.T) {
		ip := "127.0.1.1/24"
		resp, err := s.apiClient.AddIP(ctx, ip, "white", false)
		assert.NoError(t, err)
		assert.True(t, resp.Success)

		wg := sync.WaitGroup{}
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ip = fmt.Sprintf("127.0.1.%d", i)
				attempts := 11
				allowed := s.tryAuthN(ctx, attempts, "login3", "password3", ip)
				assert.Equal(t, attempts, allowed)
			}()
		}
		wg.Wait()
	})
	s.T().Run("Успешное добавление в черный список IP", func(t *testing.T) {
		ip := "127.0.0.4"
		resp, err := s.apiClient.AddIP(ctx, ip, "black", false)
		assert.NoError(t, err)
		assert.True(t, resp.Success)

		allowed := s.tryAuthN(ctx, 5, "login4", "password4", ip)
		assert.Equal(t, 0, allowed)
	})
	s.T().Run("Успешное добавление в черный список IP с перезаписью", func(t *testing.T) {
		ip := "127.0.0.5"
		resp, err := s.apiClient.AddIP(ctx, ip, "white", false)
		assert.NoError(t, err)
		assert.True(t, resp.Success)

		resp, err = s.apiClient.AddIP(ctx, ip, "black", true)
		assert.NoError(t, err)
		assert.True(t, resp.Success)

		allowed := s.tryAuthN(ctx, 5, "login5", "password5", ip)
		assert.Equal(t, 0, allowed)
	})
	s.T().Run("Неуспешное добавление в черный список IP", func(t *testing.T) {
		ip := "127.0.0.6"
		resp, err := s.apiClient.AddIP(ctx, ip, "white", false)
		assert.NoError(t, err)
		assert.True(t, resp.Success)

		resp, err = s.apiClient.AddIP(ctx, ip, "black", false)
		assert.NoError(t, err)
		assert.False(t, resp.Success)

		attempts := 15
		allowed := s.tryAuthN(ctx, attempts, "login6", "password6", ip)
		assert.Equal(t, attempts, allowed)
	})
	s.T().Run("Добавление подсетей в черный список IP", func(t *testing.T) {
		ip := "127.0.2.1/24"
		resp, err := s.apiClient.AddIP(ctx, ip, "black", false)
		assert.NoError(t, err)
		assert.True(t, resp.Success)

		for i := 0; i < 10; i++ {
			ip = fmt.Sprintf("127.0.2.%d", i)
			allowed := s.tryAuthN(ctx, 1, "login7", "password7", ip)
			assert.Equal(t, 0, allowed)
		}
	})
}

func (s *TestIPListAPI) TestDelete() {
	ctx := context.Background()
	s.T().Run("Успешное удаление из белого списка", func(t *testing.T) {
		ip := "127.0.3.1"
		listType := "white"
		login := "login8"
		password := "password8"
		resp, err := s.apiClient.AddIP(ctx, ip, listType, false)
		assert.NoError(t, err)
		assert.True(t, resp.Success)

		attempts := 15
		allowed := s.tryAuthN(ctx, attempts, login, password, ip)
		assert.Equal(t, attempts, allowed)

		resp, err = s.apiClient.DeleteIP(ctx, ip, listType)
		assert.NoError(t, err)
		assert.True(t, resp.Success)

		allowed = s.tryAuthN(ctx, 1, login, password, ip)
		assert.Equal(t, 0, allowed)
	})
	s.T().Run("Успешное удаление из черного списка", func(t *testing.T) {
		ip := "127.0.3.2"
		listType := "black"
		login := "login9"
		password := "password9"
		resp, err := s.apiClient.AddIP(ctx, ip, listType, false)
		assert.NoError(t, err)
		assert.True(t, resp.Success)

		allowed := s.tryAuthN(ctx, 1, login, password, ip)
		assert.Equal(t, 0, allowed)

		resp, err = s.apiClient.DeleteIP(ctx, ip, listType)
		assert.NoError(t, err)
		assert.True(t, resp.Success)

		attempts := 12
		allowed = s.tryAuthN(ctx, attempts, login, password, ip)
		assert.Equal(t, 9, allowed)
	})
}
