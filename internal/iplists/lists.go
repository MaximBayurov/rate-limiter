package iplists

import (
	"fmt"

	"github.com/MaximBayurov/rate-limiter/internal/logger"
	"github.com/jmoiron/sqlx"
)

type IPList interface {
	// Add добавляет IP в список
	Add(ip string, listType ListType) error

	// Update обновляет IP в списке
	Update(ip string, listType ListType) error

	// Delete удаляет IP из списка
	Delete(ip string, listType ListType) error

	// In проверяет наличие IP и возвращает тип списка
	In(ip string) (ListType, error)
}

// NewIPList возвращает реализацию списка IP.
func NewIPList(db *sqlx.DB, logger logger.Logger) IPList {
	var list IPList = &SQLIPList{
		db:     db,
		logger: logger,
	}
	return list
}

type SQLIPList struct {
	db     *sqlx.DB
	logger logger.Logger
}

// Add добавляет IP в список.
func (l *SQLIPList) Add(ip string, listType ListType) error { //nolint: dupl
	// Начинаем транзакцию
	tx, err := l.db.Beginx()
	if err != nil {
		return fmt.Errorf("transaction begin: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Вставляем IP
	query := `
		INSERT INTO ip_list (ip, type)
		VALUES (:ip, :list_type)
	`

	result, err := tx.NamedExec(
		query,
		map[string]interface{}{
			"ip":        ip,
			"list_type": listType,
		},
	)
	if err != nil {
		l.logger.Error(
			fmt.Errorf("add ip to list: %w", err).Error(),
			"ip", ip,
			"listType", listType,
		)
		return ErrNotInserted
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		l.logger.Error(
			fmt.Errorf("obtain inserted rows count: %w", err).Error(),
			"ip", ip,
			"listType", listType,
		)
		return ErrNotInserted
	}

	if rowsAffected == 0 {
		return ErrNotInserted
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit: %w", err)
	}

	return nil
}

// Update обновляет IP в списке.
func (l *SQLIPList) Update(ip string, listType ListType) error { //nolint: dupl
	// Начинаем транзакцию
	tx, err := l.db.Beginx()
	if err != nil {
		return fmt.Errorf("transaction begin: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Вставляем IP
	query := `
		UPDATE ip_list
		SET type = :list_type
		WHERE ip = :ip
	`

	result, err := tx.NamedExec(
		query,
		map[string]interface{}{
			"ip":        ip,
			"list_type": listType,
		},
	)
	if err != nil {
		l.logger.Error(
			fmt.Errorf("update ip in list: %w", err).Error(),
			"ip", ip,
			"listType", listType,
		)
		return ErrNotUpdated
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		l.logger.Error(
			fmt.Errorf("obtain updated rows count: %w", err).Error(),
			"ip", ip,
			"listType", listType,
		)
		return ErrNotUpdated
	}

	if rowsAffected == 0 {
		return ErrNotUpdated
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit: %w", err)
	}

	return nil
}

// Delete удаляет IP из списка.
func (l *SQLIPList) Delete(ip string, listType ListType) error { //nolint: dupl
	// Начинаем транзакцию
	tx, err := l.db.Beginx()
	if err != nil {
		return fmt.Errorf("transaction begin: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Вставляем IP
	query := `
		DELETE FROM ip_list
		WHERE ip=:ip AND type=:list_type
		RETURNING ip;
	`

	result, err := tx.NamedExec(
		query,
		map[string]interface{}{
			"ip":        ip,
			"list_type": listType,
		},
	)
	if err != nil {
		l.logger.Error(
			fmt.Errorf("delete ip from list: %w", err).Error(),
			"ip", ip,
			"listType", listType,
		)
		return ErrNotDeleted
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		l.logger.Error(
			fmt.Errorf("obtain deleted rows count: %w", err).Error(),
			"ip", ip,
			"listType", listType,
		)
		return ErrNotDeleted
	}

	if rowsAffected == 0 {
		return ErrNotDeleted
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit: %w", err)
	}

	return nil
}

func (l *SQLIPList) In(ip string) (ListType, error) {
	// Начинаем транзакцию
	tx, err := l.db.Beginx()
	if err != nil {
		return "", fmt.Errorf("transaction begin: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `
		SELECT type
		FROM ip_list
		WHERE ip >>= :ip
		ORDER BY ip DESC
		LIMIT 1
	`

	var listType string
	rows, err := tx.NamedQuery(
		query,
		map[string]interface{}{
			"ip": ip,
		},
	)
	defer closeRows(rows)

	if err != nil {
		return "", fmt.Errorf("check ip in list: %w", err)
	}

	rows.Next()
	if err := rows.Scan(&listType); err != nil {
		return "", fmt.Errorf("check ip scan result: %w", err)
	}

	result, ok := ParseType(listType)
	if !ok {
		return "", ErrNotIn
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("transaction commit: %w", err)
	}

	return result, nil
}

func closeRows(rows *sqlx.Rows) {
	if rows == nil {
		return
	}
	if err := rows.Close(); err != nil {
		return
	}
}
