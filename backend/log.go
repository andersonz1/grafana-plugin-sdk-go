package backend

import (
	"github.com/andersonz1/grafana-plugin-sdk-go/backend/log"
)

// Logger is the default logger instance. This can be used directly to log from
// your plugin to grafana-server with calls like backend.Logger.Debug(...).
var Logger = log.DefaultLogger
