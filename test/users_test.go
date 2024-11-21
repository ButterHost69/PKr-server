package test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/ButterHost69/PKr-server/db"
)


func TestMain(m *testing.M){
	db.InitSQLiteDatabase(true)

	code := m.Run()
	os.Exit(code)
}

type RegisterNewUserResp struct {
	Response	string	`json:"response"`	
	Username	string	`json:"username"`
}


// FIXME: [ ] Test Failing eventhough username Entry is present in the database. The error is in db.VerifyUsernameInUsersTable
// Test for POST /register/user
func TestRegisterNewUser(t *testing.T) {
	user := struct{
		Username	string 	`json:"username"`
		Password	string	`json:"password"`
	}{
		Username: "user1", 
		Password: "pass1",
	}

	url := "http://localhost:9069/register/user"
  	method := "POST"

  	payload := &bytes.Buffer{}
  	writer := multipart.NewWriter(payload)
  	err := writer.WriteField("username", user.Username)
  	if err != nil {
		t.Fatalf("Error writing field: %v", err)
  	}

  	err = writer.WriteField("password", user.Password)
  	if err != nil {
		t.Fatalf("Error writing field: %v", err)
  	}

  	err = writer.Close()
  	if err != nil {
  	  t.Fatalf("Error in closing writer: %v", err)
  	  return
  	}

  	client := &http.Client {
  	}
  	req, err := http.NewRequest(method, url, payload)
  	if err != nil {
  	  t.Fatalf("Error failed to create request: %v", err)
  	}
  	req.Header.Set("Content-Type", writer.FormDataContentType())
  	
	resp, err := client.Do(req)
  	if err != nil {
  	  t.Fatalf("Error failed to make send request: %v", err)
  	}
  	defer resp.Body.Close()

  	body, err := io.ReadAll(resp.Body)
  	if err != nil {
  	  t.Fatalf("Error failed to read from the response: %v", err)
  	}

  	t.Logf("Response status: %v", resp.Status)
	t.Logf("Response body: %s", body)

	var repsonse RegisterNewUserResp
	err = json.Unmarshal(body, &repsonse)
	if err != nil {
  	  t.Fatalf("Error failed to umarshall repsonse: %v", err)
  	}
	if resp.Status != "200 OK" && repsonse.Response == "success"{
		t.Fatalf("Error Expected Status: 200 OK  ||  Body: 'response':'success,\nreceived: Status: %s, Body: %s", resp.Status, string(body))
	} 

	t.Logf("username by server : %s", repsonse.Username)
	ifUserExists, err := db.VerifyUserExistsInUsersTable(string(repsonse.Username))
	if err != nil {
  	  t.Fatalf("Error failed to check if user entry created in Database: %v", err)
  	}
	if !ifUserExists {
		t.Fatalf("User Entry not present in Database")
	}

	t.Logf("User Entry present in Database")
}	