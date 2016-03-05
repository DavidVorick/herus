package main

import (
	"encoding/json"
  "encoding/binary"

	"github.com/boltdb/bolt"
)

type membership struct {
  DateStarted   string
  DateEnded     string
  MemberType    string
}

type payment struct {
  Processing  bool
  Processed   bool
}

type flags struct {
	Date			string
	FlaggerID	string
}

type User struct {
  ID          		int
  Username   			string
  Password    		[]byte
  Posts       		int
  Email     		  string
	Flags						[15]flags
	Payment					payment
	MembershipData	membership
}

// Borrowed from boltdb readme
// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, uint64(v))
    return b
}

func (h *herus) userHandler(w http.Response)

// Add a user to the users bucket.
// In the style of CreateUser() from the boltdb readme.
// Assumes username uniqueness already verified.
func (h *herus) initUser(user *User) error {
	return h.db.Update(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte("Users"))

		// Borrowed from boltdb readme
		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		// id, _ = users.NextSequence()
		// user.ID = int(id)
		encoded, err := json.Marshal(user)
		if err != nil {
			return err
		}
		// Stick bytes into users bucket
		return users.Put(itob(user.ID),encoded)
	})
}
