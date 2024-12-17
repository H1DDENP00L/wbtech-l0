package consts

import (
	"l0/producer"
	"sync"
)

// Просто глобальные константы для работы с кэшем)

var OrderCache = make(map[string]producer.Order)
var CacheMutex sync.RWMutex
