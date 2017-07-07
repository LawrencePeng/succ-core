package succinct

import (
	"os"
	"sync"
)

type SuccinctStore struct {
	dir string
	files []os.File
	rMux  sync.RWMutex
}

