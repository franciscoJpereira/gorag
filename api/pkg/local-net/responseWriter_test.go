package localnet_test

import (
	"encoding/json"
	"net/http"
	localnet "ragAPI/pkg/local-net"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestWritingValues(t *testing.T) {
	r := localnet.NewResponseWriter()
	values := map[string]string{
		"value-1": "1",
		"value-2": "2",
		"value-3": "3",
	}
	m, _ := json.Marshal(values)
	if _, err := r.Write(m); err != nil {
		t.Fatalf("Writing buffer: %s\n", err)
	}
	var recovered map[string]string
	if err := json.Unmarshal(r.Buf.Bytes(), &recovered); err != nil {
		t.Fatalf("Recovering data from buffer: %s\n", err)
	}
	for key, value := range values {
		if value != recovered[key] {
			t.Fatalf("Failed at key %s with value %s", key, value)
		}
	}
}

func TestUsingItAsWriter(t *testing.T) {
	r := localnet.NewResponseWriter()
	rec, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	e := echo.New()
	c := e.NewContext(rec, r)
	basicGet := func(c echo.Context) error {
		c.Response().Header().Set("Value1", "1")
		return c.String(http.StatusFound, "OK Write")
	}
	basicGet(c)
	if r.Header().Get("value1") != "1" {
		t.Fatalf("Invalid header value: %s", r.Header().Get("value1"))
	}
	response := r.Buf.String()
	if response != "OK Write" {
		t.Fatalf("Invalid recovered response: %s\n", response)
	}
}
