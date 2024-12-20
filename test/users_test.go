package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/ButterHost69/PKr-server/db"
)

var (
	userDetails		UserDetails
	dummyDetails	UserDetails
)

func TestMain(m *testing.M){
	var err error
	_, err = db.InitSQLiteDatabase(true, "../test_database.db")
	if err != nil {
		fmt.Printf("error occured in initiating database.\nError: %e\n", err)
		os.Exit(1)
	}
	err = db.InsertDummyData()
		if err != nil {
			fmt.Printf("error: Could Not Create Tables.\nError: %v", err)
			os.Exit(1)
		}	
	fmt.Println("database initiated...")

	dummyDetails.Username = "user#123"
	dummyDetails.Password = "password123"
	dummyDetails.Workspace_Name = "WorkspaceA"
	dummyDetails.Connection_Username = "user#456"
	
	code := m.Run()
	os.Exit(code)
}

type RegisterNewUserResp struct {
	Response	string	`json:"response"`	
	Username	string	`json:"username"`
}

type GenericResp struct {
	Response	string	`json:"response"`	
}

// Main struct to store details for all throughout tests
type UserDetails struct {
	Username			string	`json:"username"`
	Password			string	`json:"password"`
	Workspace_Name		string	`json:"workspace_name"`
	Connection_Username	string	`json:"connection_username"`
}

// FIXME: [X] Test Failing eventhough username Entry is present in the database. The error is in db.VerifyUsernameInUsersTable
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
	userDetails.Username = repsonse.Username
	userDetails.Password = user.Password
}	

// TODO: [ ] Write a Failing Test, for auth 
// /register/workspace
func TestRegisterWorkspace(t *testing.T){
	workspace_name := "WorkSpace1"

	url := "http://localhost:9069/register/workspace"
  	method := "POST"

  	payload := &bytes.Buffer{}
  	writer := multipart.NewWriter(payload)
  	err := writer.WriteField("username", userDetails.Username)
  	if err != nil {
		t.Fatalf("Error writing field: %v", err)
  	}

  	err = writer.WriteField("password", userDetails.Password)
  	if err != nil {
		t.Fatalf("Error writing field: %v", err)
  	}

	err = writer.WriteField("workspace_name", workspace_name)
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

	var repsonse GenericResp
	err = json.Unmarshal(body, &repsonse)
	if err != nil {
  	  t.Fatalf("Error failed to umarshall repsonse: %v", err)
  	}
	if resp.Status != "200 OK" && repsonse.Response == "success"{
		t.Fatalf("Error Expected Status: 200 OK  ||  Body: 'response':'success,\nreceived: Status: %s, Body: %s", resp.Status, string(body))
	}

	ifExists, err := db.CheckIfWorkspaceExists(userDetails.Username, workspace_name)
	if err != nil {
		t.Fatalf("Error failed to verify if workspace registered in db: %v", err)
	}

	if !ifExists {
		t.Fatalf("Workspace Entry not present in Database")
	}

	t.Logf("User Entry present in Database")
	userDetails.Workspace_Name = workspace_name
}

// TODO: [X] TEST ~ RegisterUserToWorkspace
// /register/user_to_workspace
func TestRegisterUserToWorkspace(t *testing.T){ 
	connection_username := "userWorkspace#123"

	url := "http://localhost:9069/register/user_to_workspace"
  	method := "POST"

	// username := ctx.PostForm("username")
	// password := ctx.PostForm("password")
	// workspace_name := ctx.PostForm("workspace_name")
	// connection_username := ctx.PostForm("connection_username")

  	payload := &bytes.Buffer{}
  	writer := multipart.NewWriter(payload)
  	err := writer.WriteField("username", userDetails.Username)
  	if err != nil {
		t.Fatalf("Error writing field: %v", err)
  	}

  	err = writer.WriteField("password", userDetails.Password)
  	if err != nil {
		t.Fatalf("Error writing field: %v", err)
  	}

	err = writer.WriteField("workspace_name", userDetails.Workspace_Name)
  	if err != nil {
		t.Fatalf("Error writing field: %v", err)
  	}
	 
	err = writer.WriteField("connection_username", connection_username)
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

	var repsonse GenericResp
	err = json.Unmarshal(body, &repsonse)
	if err != nil {
  	  t.Fatalf("Error failed to umarshall repsonse: %v", err)
  	}
	if resp.Status != "200 OK" && repsonse.Response == "success"{
		t.Fatalf("Error Expected Status: 200 OK  ||  Body: 'response':'success,\nreceived: Status: %s, Body: %s", resp.Status, string(body))
	}

	ifExists, err := db.VerifyConnectionUserExistsInWorkspaceConnectionTable(userDetails.Workspace_Name, userDetails.Username, connection_username)
	if err != nil {
		t.Fatalf("Error failed to verify if connection username attached to workspace in db: %v", err)
	}

	if !ifExists {
		t.Fatalf("User Connection to Workspace not present in Database")
	}

	t.Logf("User Connection to Workspace is present in Database")
	userDetails.Connection_Username = connection_username
}

// FIXME: [ ] Not Working 
// /whats/new
func TestGetAllMyConnectedWorkspaceInfo(t *testing.T){
	// connection_username := "userWorkspace#123"

	url := "http://localhost:9069/whats/new"
  	method := "GET"

	// username := ctx.PostForm("username")
	// password := ctx.PostForm("password")

  	payload := &bytes.Buffer{}
  	writer := multipart.NewWriter(payload)
  	err := writer.WriteField("username", dummyDetails.Username)
  	if err != nil {
		t.Fatalf("Error writing field: %v", err)
  	}

  	err = writer.WriteField("password", dummyDetails.Password)
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

}

func TestClose(t *testing.T){
	t.Cleanup(func() {
		db.CloseSQLiteDatabase()
	})
}