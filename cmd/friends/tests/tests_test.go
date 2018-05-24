package tests

import (
	"os"
	"testing"

	"github.com/bobintornado/spfriends/cmd/friends/handlers"
	"github.com/bobintornado/spfriends/internal/platform/tests"
	"github.com/bobintornado/spfriends/internal/platform/web"
)

var a *web.App
var test *tests.Test

// TestMain is the entry point for testing.
func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	test = tests.New()
	defer test.TearDown()
	a = handlers.API(test.MasterDB).(*web.App)
	return m.Run()
}
