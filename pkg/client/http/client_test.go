package client

import (
	"Desktop/golangProjects/CRUD/pkg"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

type testSetup struct {
	mux *mux.Router
}

func setupDeps(t *testing.T) testSetup {
	createHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("val", "create")
		w.WriteHeader(http.StatusOK)
	}
	userHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Add("val", "read")
			user := &pkg.HttpData{
				Name: "jason",
				Age:  21,
			}
			msg, err := json.Marshal(user)
			require.NoError(t, err)
			_, err = w.Write(msg)
			require.NoError(t, err)

		case http.MethodPatch:
			w.Header().Add("val", "update")
			w.WriteHeader(http.StatusOK)

		case http.MethodDelete:
			w.Header().Add("val", "delete")
			w.WriteHeader(http.StatusOK)
		}
	}

	mux := mux.NewRouter()

	mux.HandleFunc(pkg.CREATEADDRROUTE, createHandler)
	mux.HandleFunc(fmt.Sprintf("%s%s", pkg.USERADDROUTE, "jason"), userHandler)

	return testSetup{
		mux: mux,
	}
}

func TestClient(t *testing.T) {
	testSetup := setupDeps(t)
	mux := testSetup.mux
	testServer := httptest.NewServer(mux)

	defer testServer.Close()

	client := New(testServer.URL)

	resp, err := client.CreateUser("jason", 21)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf(createResponse, "jason", 21), resp)

	resp, err = client.ReadUser("jason")
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf(readResponse, "jason", 21), resp)

	resp, err = client.UpdateUser("jason", 22)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf(updateResponse, "jason", 22), resp)

	resp, err = client.DeleteUser("jason")
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf(deleteResponse, "jason"), resp)

}
