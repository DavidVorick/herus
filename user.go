package main

import (
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

const (
	loginPage = "/login.go"
	userPrefix= "/u/"
)

// flags can be given to users if they are believed to be behaving in an
// abusive way.
type flags struct {
	Time        time.Time
	FlaggerName string
}

// User contains all of the information
type User struct {
	// Generic user information.
	Email          string
	HashedPassword []byte // blake2b(salt+username+password)
	Salt           string
	Username       string

	Flags          [15]flags
	Posts          uint64
	Votes          uint64

	// Membership information.
	DateStarted string
	DateEnded   string
	MemberType  string
}

// Borrowed from boltdb readme itob returns an 8-byte big endian representation
// of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// Add a user to the users bucket.
// In the style of CreateUser() from the boltdb readme.
// Assumes username uniqueness already verified.
func (h *herus) initUser(user *User) error {
	return h.db.Update(func(tx *bolt.Tx) error {
		users := tx.Bucket(bucketUsers)
		encoded, err := json.Marshal(user)
		if err != nil {
			return err
		}
		// Stick bytes into users bucket.
		return users.Put(itob(user.ID), encoded)
	})
}

// Voting needs to be aware of what each user has voted for. Some archetecture
// that tracks when a user has voted for something and remembers which way the
// vote went. So, we can use a map. Page+item+user -> vote value. Which means
// you can't easily pull a full list of things a user has voted on without
// doing a bigger scan.

// [prefix][user][prefix][page][prefix][item] -> vote value. That allows you to
// scan the bucket to get a vote history. Upload history can be retrieved by
// scanning everything as well. If it becomes necessary to pull up those stats
// regularly, another thing can be added.
