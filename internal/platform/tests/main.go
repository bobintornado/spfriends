package tests

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"github.com/bobintornado/spfriends/internal/platform/db"
	"github.com/bobintornado/spfriends/internal/platform/docker"
	"github.com/bobintornado/spfriends/internal/platform/web"
	"github.com/pborman/uuid"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

// Test owns state for running/shutting down tests.
type Test struct {
	MasterDB  *db.DB
	container *docker.Container
}

// New is the entry point for tests.
func New() *Test {
	var test Test

	// ============================================================
	// Startup Neo4j container

	var err error
	test.container, err = docker.StartNeo4j()
	if err != nil {
		log.Fatalln(err)
	}

	// ============================================================
	// Configuration
	bytes, err := exec.Command("docker-machine", "ip").Output()
	dockerHost := strings.TrimSpace(string(bytes))
	if err != nil {
		// if docker-machine is not avaiable, then probably localhost is docker host
		dockerHost = "localhost"
	}

	dbHost := fmt.Sprintf("bolt://%s:%s", dockerHost, test.container.Port)

	// ============================================================
	// Start Mongo

	log.Println("main : Started : Initialize Neo4j")
	test.MasterDB, err = db.New(dbHost, 25)
	if err != nil {
		log.Fatalf("startup : Register DB : %v", err)
	}

	return &test
}

// TearDown is used for shutting down tests. Calling this should be
// done in a defer immediately after calling New.
func (t *Test) TearDown() {
	t.MasterDB.Close()
	if err := docker.StopNeo4j(t.container); err != nil {
		log.Println(err)
	}
}

// Recover is used to prevent panics from allowing the test to cleanup.
func Recover(t *testing.T) {
	if r := recover(); r != nil {
		t.Fatal("Unhandled Exception:", string(debug.Stack()))
	}
}

// Context returns an app level context for testing.
func Context() context.Context {
	values := web.Values{
		TraceID: uuid.New(),
		Now:     time.Now(),
	}

	return context.WithValue(context.Background(), web.KeyValues, &values)
}
