package relationship

import (
	"context"
	"errors"
	"io"

	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

func CreateFriendship(ctx context.Context, dbConn *bolt.Conn, cfr *CreateFriendshipRequest) error {
	// if there is block, return error
	stmt, err := (*dbConn).PrepareNeo(`
		RETURN EXISTS( (:User {email: {email1} })-[:BLOCK]-(:User {email: {email2} }) )
	`)
	if err != nil {
		return err
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"email1": cfr.Friends[0],
		"email2": cfr.Friends[1],
	})
	if err != nil {
		return err
	}

	data, _, err := rows.NextNeo()
	if err != nil {
		return err
	}

	// block existed
	if data[0].(bool) {
		return errors.New("User block existed, can't create friendship")
	}

	stmt.Close()

	stmt, err = (*dbConn).PrepareNeo(`
		MERGE (user1:User { email: {email1} })
		MERGE (user2:User { email: {email2} })
		MERGE (user1)-[:FRIEND_OF]->(user2)
	`)
	if err != nil {
		return err
	}

	_, err = stmt.ExecNeo(map[string]interface{}{
		"email1": cfr.Friends[0],
		"email2": cfr.Friends[1],
	})
	if err != nil {
		return err
	}

	stmt.Close()

	return nil
}

func ListFriendsOfUser(ctx context.Context, dbConn *bolt.Conn, gflr *GetFriendsListRequest) (list []string, err error) {
	list, err = GetFriendsOfUserByEmail(dbConn, gflr.Email)
	if err != nil {
		return list, err
	}

	return list, nil
}

func ListCommonFriends(ctx context.Context, dbConn *bolt.Conn, gcflr *GetCommonFriendsListRequest) (list []string, err error) {
	list = []string{}

	stmt, err := (*dbConn).PrepareNeo(`
		MATCH (:User {email: { email1 }})-[:FRIEND_OF]-(commonFriends)
				-[:FRIEND_OF]-(:User {email: { email2 }})
		RETURN commonFriends.email
	`)
	if err != nil {
		return list, err
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"email1": gcflr.Friends[0],
		"email2": gcflr.Friends[1],
	})
	if err != nil {
		return list, err
	}

	for true {
		row, _, err := rows.NextNeo()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []string{}, err
		}
		list = append(list, row[0].(string))
	}

	stmt.Close()

	return list, nil
}

func CreateSubscription(ctx context.Context, dbConn *bolt.Conn, csr *CreateSubscriptionRequest) error {
	// create user and relationship if not existed yet
	stmt, err := (*dbConn).PrepareNeo(`
		MERGE (requestor:User { email: {requestor} })
		MERGE (target:User { email: {target} })
		MERGE (requestor)-[:SUBSCRIBE_TO]->(target)
	`)
	if err != nil {
		return err
	}

	_, err = stmt.ExecNeo(map[string]interface{}{
		"requestor": csr.Requestor,
		"target":    csr.Target,
	})
	if err != nil {
		return err
	}

	stmt.Close()

	return nil
}

func CreateBlock(ctx context.Context, dbConn *bolt.Conn, csr *CreateBlockRequest) error {
	// create user and relationship if not existed yet
	stmt, err := (*dbConn).PrepareNeo(`
		MERGE (requestor:User { email: {requestor} })
		MERGE (target:User { email: {target} })
		MERGE (requestor)-[:BLOCK]->(target)
	`)
	if err != nil {
		return err
	}

	_, err = stmt.ExecNeo(map[string]interface{}{
		"requestor": csr.Requestor,
		"target":    csr.Target,
	})
	if err != nil {
		return err
	}

	stmt.Close()

	return nil
}

// get list of user who are friends of target user
func GetFriendsOfUserByEmail(dbConn *bolt.Conn, email string) (list []string, err error) {
	// create user and relationship if not existed yet
	stmt, err := (*dbConn).PrepareNeo(`
		MATCH (:User { email: {email} })-[:FRIEND_OF]-(friend)
		RETURN friend.email
	`)
	if err != nil {
		return list, err
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"email": email,
	})
	if err != nil {
		return list, err
	}

	for true {
		row, _, err := rows.NextNeo()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []string{}, err
		}
		list = append(list, row[0].(string))
	}

	stmt.Close()

	return list, nil
}

// get list of user who subscribes to the target user
func GetSubscribersOfUserByEmail(dbConn *bolt.Conn, email string) (list []string, err error) {
	stmt, err := (*dbConn).PrepareNeo(`
		MATCH (:User { email: {email} })<-[:SUBSCRIBE_TO]-(subscriber)
		RETURN subscriber.email
	`)
	if err != nil {
		return list, err
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"email": email,
	})
	if err != nil {
		return list, err
	}

	for true {
		row, _, err := rows.NextNeo()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []string{}, err
		}
		list = append(list, row[0].(string))
	}

	stmt.Close()

	return list, nil
}

// get list of user who blocks the target user
func GetBlockersToUserByEmail(dbConn *bolt.Conn, email string) (list []string, err error) {
	stmt, err := (*dbConn).PrepareNeo(`
		MATCH (:User { email: {email} })<-[:BLOCK]-(blocker)
		RETURN blocker.email
	`)
	if err != nil {
		return list, err
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"email": email,
	})
	if err != nil {
		return list, err
	}

	for true {
		row, _, err := rows.NextNeo()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []string{}, err
		}
		list = append(list, row[0].(string))
	}

	stmt.Close()

	return list, nil
}

func GetUpdateList(ctx context.Context, dbConn *bolt.Conn, csr *GetUpdateListRequest) (list []string, err error) {

	friendsList, err := GetFriendsOfUserByEmail(dbConn, csr.Sender)
	if err != nil {
		return list, err
	}

	subscriberList, err := GetSubscribersOfUserByEmail(dbConn, csr.Sender)
	if err != nil {
		return list, err
	}

	mentionedSet := ExtractEmails(csr.Text)

	blockerList, err := GetBlockersToUserByEmail(dbConn, csr.Sender)
	if err != nil {
		return list, err
	}

	totalList := append(friendsList, subscriberList...)
	totalList = append(totalList, mentionedSet...)

	updateMap := map[string]bool{}
	for _, v := range totalList {
		updateMap[v] = true
	}
	for _, v := range blockerList {
		delete(updateMap, v)
	}

	for k := range updateMap {
		list = append(list, k)
	}

	return list, nil
}
