package relationship_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/bobintornado/spfriends/internal/platform/tests"
	"github.com/bobintornado/spfriends/internal/relationship"
)

var test *tests.Test

// TestMain is the entry point for testing.
func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	test = tests.New()
	defer test.TearDown()
	return m.Run()
}

// TestCreate validates we can create friendship and retrieve it from DB
func TestFriendship(t *testing.T) {
	defer tests.Recover(t)

	t.Log("Given the need to validate creation of friendship relationship.")
	{

		ctx := tests.Context()

		dbConn, err := test.MasterDB.OpenPool()
		// waiting for neo4j to be ready
		for i := 1; i <= 1000; i++ {
			if err == nil {
				break
			}

			time.Sleep(10 * time.Millisecond)
			dbConn, err = test.MasterDB.OpenPool()
		}

		defer (*dbConn).Close()

		cfr := relationship.CreateFriendshipRequest{
			Friends: [2]string{"winston@sp.com", "miccheng@sp.com"},
		}

		err = relationship.CreateFriendship(ctx, dbConn, &cfr)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to create friendship relationship: %s.", tests.Failed, err)
		}

		winFriends, err := relationship.GetFriendsOfUserByEmail(dbConn, "winston@sp.com")
		if err != nil {
			t.Fatalf("\t%s\tShould be able to get friends list : %s.", tests.Failed, err)
		}

		miccFriends, err := relationship.GetFriendsOfUserByEmail(dbConn, "miccheng@sp.com")
		if err != nil {
			t.Fatalf("\t%s\tShould be able to get friends list : %s.", tests.Failed, err)
		}

		list1 := []string{"miccheng@sp.com"}
		list2 := []string{"winston@sp.com"}

		if !reflect.DeepEqual(winFriends, list1) {
			t.Logf("\t\tGot : %+v", winFriends)
			t.Logf("\t\tWant: %+v", list1)
			t.Fatalf("\t%s\tShould get back as friends after friendship.", tests.Failed)
		}

		if !reflect.DeepEqual(miccFriends, list2) {
			t.Logf("\t\tGot : %+v", winFriends)
			t.Logf("\t\tWant: %+v", list2)
			t.Fatalf("\t%s\tShould get back as friends after friendship.", tests.Failed)
		}
	}
}
