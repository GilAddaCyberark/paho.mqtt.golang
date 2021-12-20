package mqtt

import (
	"crypto/tls"
	"errors"
	"net"
	"net/url"
	"sync"
	"time"
)

type CustomDialer = func(uri *url.URL, tlsc *tls.Config, timeout time.Duration, args ...interface{}) (net.Conn, error)

type CustomDialerMgr struct {
	lock    sync.Mutex
	dialers map[string]CustomDialer
}

var (
	dialerMgr *CustomDialerMgr
)

func CustomDialerInit() {
	dialerMgr = &CustomDialerMgr{
		dialers: make(map[string]CustomDialer),
	}
}

func AddCustomDialer(schema string, fn CustomDialer) error {
	return dialerMgr.AddDialer(schema, fn)
}

func GetCustomDialer(schema string) CustomDialer {
	return dialerMgr.GetDialer(schema)
}

func (t *CustomDialerMgr) AddDialer(schema string, fn CustomDialer) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.dialers == nil {
		return nil
	}
	if _, ok := t.dialers[schema]; ok {
		return errors.New("dialer of schema " + schema + " already exists")
	}

	t.dialers[schema] = fn
	return nil
}

func (t *CustomDialerMgr) GetDialer(schema string) CustomDialer {
	if t == nil {
		return nil
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	if t.dialers == nil {
		return nil
	}
	fn, ok := t.dialers[schema]
	if ok {
		return fn
	} else {
		return nil
	}
}
