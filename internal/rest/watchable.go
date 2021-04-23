package rest

import (
	"sync"

	"github.com/toms1441/chess-server/internal/model"
)

var watchablemtx sync.Mutex
var watchables []model.Watchable
