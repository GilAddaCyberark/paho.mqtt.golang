package mqtt

import (
	"crypto/tls"
	"net"
	"net/url"
	"testing"
	"time"
)

const testScheme = "my_network_scheme"

func Test_CustomMgr_AddDialers(t *testing.T) {
	CustomDialerInit()
	var myDialerFunc CustomDialer = func(uri *url.URL, tlsc *tls.Config, timeout time.Duration, args ...interface{}) (net.Conn, error) {
		return nil, nil
	}
	dialerMgr.AddDialer(testScheme, myDialerFunc)
	customDialer := dialerMgr.GetDialer(testScheme)
	if customDialer == nil {
		t.Errorf("no customn dialer able to be pulled")
	}
}

func Test_CustomMgr_NoInit(t *testing.T) {
	customDialer2 := dialerMgr.GetDialer("")
	if customDialer2 != nil {
		t.Errorf("custom dialer pulled with empty scheme")
	}
}
func Test_CustomMgr_EmptyScheme(t *testing.T) {
	CustomDialerInit()
	customDialer2 := dialerMgr.GetDialer("")
	if customDialer2 != nil {
		t.Errorf("custom dialer pulled with empty scheme")
	}
}

func Test_CustomMgr_GetWrongDialers(t *testing.T) {
	CustomDialerInit()
	customDialer2 := dialerMgr.GetDialer("wrong scheme")
	if customDialer2 != nil {
		t.Errorf("custom dialer pulled with wrong scheme")
	}
}

func Test_Custom_Dialer(t *testing.T) {
	CustomDialerInit()
	var myDialerFunc CustomDialer = func(uri *url.URL, tlsc *tls.Config, timeout time.Duration, args ...interface{}) (net.Conn, error) {
		server, client := net.Pipe()
		go func() {
			time.Sleep(5 * time.Second)
			server.Close()
		}()
		return client, nil
	}
	dialerMgr.AddDialer(testScheme, myDialerFunc)
	customDialerfunc := dialerMgr.GetDialer(testScheme)
	if customDialerfunc == nil {
		t.Error("no func to retreive")
	}
	timeout := time.Duration(1 * time.Second)
	endpoint, err := url.Parse("openssl://localhost:9090")
	if err != nil {
		t.Errorf("%v\n", err)
	}
	conn, err := customDialerfunc(endpoint, &tls.Config{}, timeout, 12, 15)
	if err != nil {
		t.Errorf("%v\n", err)
	}
	if conn == nil {
		t.Error("connection not created", err)
	}
	conn.Close()
}
