package server

import (
	"Desktop/golangProjects/CRUD/pkg"
	"Desktop/golangProjects/CRUD/pkg/database"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

type testSetup struct {
	tmpFileName string
	db          database.Database
}

func setupTestDeps(t *testing.T) testSetup {
	tmpFile, err := os.CreateTemp("", "testdb-*.db")
	require.NoError(t, err)
	db, err := database.New(tmpFile.Name(), 0666, nil)
	require.NoError(t, err)
	return testSetup{
		tmpFileName: tmpFile.Name(),
		db:          db,
	}
}

func createRequest(t *testing.T, name string, age int, url string, route string, method string) *http.Request {
	user := &pkg.HttpData{
		Name: name,
		Age:  age,
	}
	body, err := json.Marshal(user)
	require.NoError(t, err)
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", url, route), bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Add("Content-Type", "application/json")
	return req
}

func readHelper(t *testing.T, url string, name string) pkg.HttpData {
	req := createRequest(t, name, 0, url, fmt.Sprintf("%s%s", pkg.USERADDROUTE, name), http.MethodGet)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var user pkg.HttpData
	err = json.Unmarshal(body, &user)
	require.NoError(t, err)
	return user
}

func TestServer(t *testing.T) {
	testSetup := setupTestDeps(t)
	defer testSetup.db.Close()
	defer os.Remove(testSetup.tmpFileName)
	serv := New(testSetup.db)
	mux := mux.NewRouter()
	mux.HandleFunc(pkg.CREATEADDRROUTE, serv.HandleCreate)
	mux.HandleFunc(fmt.Sprintf("%s%s", pkg.USERADDROUTE, "{name}"), serv.HandleUsers)
	server := httptest.NewServer(mux)
	defer server.Close()

	// create
	req := createRequest(t, "jason", 21, server.URL, pkg.CREATEADDRROUTE, http.MethodPost)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)

	//read
	user := readHelper(t, server.URL, "jason")
	require.Equal(t, "jason", user.Name)
	require.Equal(t, 21, user.Age)

	// update
	req = createRequest(t, "jason", 29, server.URL, fmt.Sprintf("%s%s", pkg.USERADDROUTE, "jason"), http.MethodPatch)
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)

	// read again
	user = readHelper(t, server.URL, "jason")
	require.Equal(t, "jason", user.Name)
	require.Equal(t, 29, user.Age)

	// delete
	req = createRequest(t, "jason", 0, server.URL, fmt.Sprintf("%s%s", pkg.USERADDROUTE, "jason"), http.MethodDelete)
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)

	// try read one more time, but it should fail
	req = createRequest(t, "jason", 0, server.URL, fmt.Sprintf("%s%s", pkg.USERADDROUTE, "jason"), http.MethodGet)
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, string(body), "jason does not exist in our database")
}
