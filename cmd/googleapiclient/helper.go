package googleapiclient

import (
	"context"
	"errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
)

func getCredentials(keypath, projectId string, scopes ...string) (*http.Client, error) {

	f, err := ioutil.ReadFile(keypath)
	if err != nil {
		return nil, err
	}

	cred, err := google.CredentialsFromJSON(context.Background(), f, scopes...)
	if err != nil {
		return nil, err
	}

	if cred.ProjectID == "" {
		return nil, errors.New("failed to authenticate")
	}

	if cred.ProjectID != projectId {
		return nil, errors.New("the project id found in the credentials file does not match the project id")
	}

	c := oauth2.NewClient(context.Background(), cred.TokenSource)

	return c, nil
}
