package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.etcd.io/bbolt"
)

func TestDB(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testdb-*.db")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name()) // clean up

	db, err := bbolt.Open(tmpFile.Name(), 0666, nil)
	require.NoError(t, err)
	defer db.Close()

	boltDB := &boltDB{
		DB: db,
	}

	err = boltDB.Write("jason", "21")
	require.NoError(t, err)
	ageAsString, err := boltDB.Read("jason")
	require.NoError(t, err)
	require.Equal(t, "21", ageAsString)
	_, err = boltDB.Read("jesus")
	require.Error(t, err, ErrUserNotFound)
	err = boltDB.Write("jason", "27")
	require.NoError(t, err)
	ageAsString, err = boltDB.Read("jason")
	require.NoError(t, err)
	require.Equal(t, "27", ageAsString)
	err = boltDB.Delete("jason")
	require.NoError(t, err)
	err = boltDB.Delete("jason")
	require.NoError(t, err)
	_, err = boltDB.Read("jesus")
	require.Error(t, err)

}
