package app

import (
	"errors"
	"fmt"
)

var ErrUnsupportedIPListType = errors.New("неподдерживаемый тип списка IP")

var (
	ErrLoginDisallow    = fmt.Errorf("превышено количество попыток по логину%w", ErrBase)
	ErrPasswordDisallow = fmt.Errorf("превышено количество попыток по паролю%w", ErrBase)
	ErrIPDisallow       = fmt.Errorf("превышено количество попыток по IP%w", ErrBase)
)

var ErrBase = errors.New("")
