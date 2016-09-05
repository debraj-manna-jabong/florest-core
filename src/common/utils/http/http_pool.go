package http

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

type Config struct {
	// MaxConn maximum number of connections
	MaxConn int
	// MaxIdleConnsPerHost, if non-zero, controls the maximum idle
	// (keep-alive) to keep per-host.
	MaxIdleConns int
	// ResponseHeaderTimeout, if non-zero, specifies the amount of
	// time to wait for a server's response headers after fully
	// writing the request (including its body, if any). This
	// time does not include the time to read the response body.
	ResponseHeaderTimeout int // in seconds
	// DisableKeepAlives, if true, prevents re-use of TCP connections
	// between different HTTP requests.
	DisableKeepAlives bool
}

type pool struct {
	transport *http.Transport
	mutex     *sync.Mutex
	maxConn   int
	numConn   int
}

var poolObj *pool

const (
	ERR_MAXCON_REACHED = "Max number of connections reached, cant take new connections"
	ERR_MINCON_REACHED = "Min number of connections reached, cant decrement connections !!"
	MIN_NUM_CON        = 0
)

// InitConnPool initialized connection pool
func InitConnPool(conf *Config) {
	poolObj = new(pool)
	// init pool fields
	poolObj.transport = &http.Transport{MaxIdleConnsPerHost: conf.MaxIdleConns,
		ResponseHeaderTimeout: time.Duration(conf.ResponseHeaderTimeout) * time.Second,
		DisableKeepAlives:     conf.DisableKeepAlives}
	poolObj.mutex = &sync.Mutex{}
	poolObj.maxConn = conf.MaxConn
}

// isPoolSet is connection pool set
func isPoolSet() bool {
	return poolObj != nil
}

// incNumCon increment number of connections
func incNumCon() error {
	poolObj.mutex.Lock()
	defer poolObj.mutex.Unlock()
	if poolObj.numConn < poolObj.maxConn {
		poolObj.numConn++
		return nil
	} else {
		return errors.New(ERR_MAXCON_REACHED)
	}

}

// decNumCon decrement number of connections
func decNumCon() error {
	poolObj.mutex.Lock()
	defer poolObj.mutex.Unlock()
	if poolObj.numConn == MIN_NUM_CON {
		return errors.New(ERR_MINCON_REACHED)
	} else {
		poolObj.numConn--
		return nil
	}
}
