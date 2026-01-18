package iplists

import (
	"errors"
	"fmt"
)

var (
	ErrNotInserted = fmt.Errorf("не удалось добавить ip в список%w", ErrBase)
	ErrNotUpdated  = fmt.Errorf("не удалось обновить ip в списке%w", ErrBase)
	ErrNotDeleted  = fmt.Errorf("не удалось удалить ip из списка%w", ErrBase)
	ErrNotIn       = fmt.Errorf("ip не содержится ни в одном из списков%w", ErrBase)
)

var ErrBase = errors.New("")
