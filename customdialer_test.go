package mqtt

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var myDialer CustomDialer = func(uri *url.URL, tlsc *tls.Config, timeout time.Duration, headers http.Header) (net.Conn, error) {
	return nil, nil
}

const testScheme = "my_network_scheme"

func Test_CustomMgr_AddDialers(t *testing.T) {
	CustomDialerInit()
	dialerMgr.AddDialer(testScheme, myDialer)
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
