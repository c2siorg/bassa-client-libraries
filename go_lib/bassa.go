//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package bassa

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/gojektech/heimdall"
	"github.com/gojektech/heimdall/httpclient"
	"github.com/hokaccha/go-prettyjson"
)

// Bassa : Bassa Go object
type Bassa struct {
	apiURL     string
	token      string
	timeout    int
	retryCount int
	httpClient *httpclient.Client
}

var (
	errBadFormat        = errors.New("invalid format")
	errIncompleteParams = errors.New("Some fields are not valid or empty")
	emailRegexp         = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// validateFormat : Helper function to validate email address
func validateFormat(email string) error {
	if !emailRegexp.MatchString(email) {
		return errBadFormat
	}
	return nil
}

// Init : Initialization of Bassa
func (b *Bassa) Init(apiURL string, timeout int, retryCount int) {
	if apiURL == "" || timeout == 0 {
		panic(errIncompleteParams)
	}
	u, err := url.Parse(apiURL)
	if err != nil {
		fmt.Println(u)
		panic(err)
	} else {
		b.apiURL = apiURL
		b.timeout = timeout
		b.retryCount = retryCount
		b.token = ""
		timeout := time.Duration(timeout) * time.Millisecond
		httpClient := httpclient.NewClient(
			httpclient.WithHTTPTimeout(timeout),
			httpclient.WithRetryCount(retryCount),
			httpclient.WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(10*time.Millisecond, 50*time.Millisecond))),
		)
		b.httpClient = httpClient
	}
}

// Login : Function to login as a user
func (b *Bassa) Login(userName string, password string) {
	if userName == "" || password == "" {
		panic(errIncompleteParams)
	}
	endpoint := "/api/login"
	apiURL := b.apiURL + endpoint

	form := url.Values{}
	form.Add("user_name", userName)
	form.Add("password", password)

	response, err := http.PostForm(apiURL, form)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	b.token = response.Header["Token"][0]

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(string(respBody))
		panic(err)
	}
}

// AddRegularUserRequest : Function to login as a user
func (b *Bassa) AddRegularUserRequest(userName string, password string, email string) {
	if userName == "" || password == "" || email == "" {
		panic(errIncompleteParams)
	}

	err := validateFormat(email)
	if err != nil {
		panic(err)
	}

	endpoint := "/api/regularuser"
	apiURL := b.apiURL + endpoint

	requestBody, err := json.Marshal(map[string]string{
		"user_name": userName,
		"password":  password,
		"email":     email})
	if err != nil {
		panic(err)
	}
	request, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}
	request.Header.Set("token", b.token)
	response, err := b.httpClient.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(string(respBody))
		panic(err)
	}
}

// AddUserRequest : Function to login as a user
func (b *Bassa) AddUserRequest(userName string, password string, email string, authLevel int) {
	if userName == "" || password == "" || email == "" {
		panic(errIncompleteParams)
	}

	err := validateFormat(email)
	if err != nil {
		panic(err)
	}

	endpoint := "/api/user"
	apiURL := b.apiURL + endpoint

	requestBody := []byte(fmt.Sprintf("{user_name:\"%s\", password: \"%s\", email: \"%s\", auth: %d}", userName, password, email, authLevel))

	request, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}
	request.Header.Set("token", b.token)
	response, err := b.httpClient.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	var r interface{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		panic(err)
	}
	out, err := prettyjson.Marshal(r)
	fmt.Println(string(out))
}

// RemoveUserRequest : Function to remove user
func (b *Bassa) RemoveUserRequest(userName string) string {
	if userName == "" {
		panic(errIncompleteParams)
	}

	endpoint := "/api/user" + "/" + userName
	apiURL := b.apiURL + endpoint

	request, err := http.NewRequest("DELETE", apiURL, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("token", b.token)
	response, err := b.httpClient.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	var r interface{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		panic(err)
	}
	out, err := prettyjson.Marshal(r)
	return string(out)
}

// UpdateUserRequest : Function to update user request
func (b *Bassa) UpdateUserRequest(userName string, newUserName string, password string, authLevel int, email string) {
	if userName == "" || password == "" || email == "" || newUserName == "" {
		panic(errIncompleteParams)
	}

	err := validateFormat(email)
	if err != nil {
		panic(err)
	}

	endpoint := "/api/user"
	apiURL := b.apiURL + endpoint + "/" + userName

	requestBody := []byte(fmt.Sprintf("{user_name:\"%s\", password: \"%s\", email: \"%s\", auth_level: %d}", newUserName, password, email, authLevel))

	request, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}
	request.Header.Set("token", b.token)
	response, err := b.httpClient.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	var r interface{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		panic(err)
	}
	out, err := prettyjson.Marshal(r)
	fmt.Println(string(out))
}

// GetUserRequest : Function to get user request
func (b *Bassa) GetUserRequest() string {

	endpoint := "/api/user"
	apiURL := b.apiURL + endpoint

	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("token", b.token)
	response, err := b.httpClient.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	var r interface{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		panic(err)
	}
	out, err := prettyjson.Marshal(r)
	fmt.Println(string(out))
	return string(out)
}

// GetUserSignupRequests : Function to get user signup requests
func (b *Bassa) GetUserSignupRequests() string {

	endpoint := "/api/user/requests"
	apiURL := b.apiURL + endpoint

	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("token", b.token)
	response, err := b.httpClient.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	var r interface{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		panic(err)
	}
	out, err := prettyjson.Marshal(r)
	fmt.Println(string(out))
	return string(out)
}

// ApproveUserRequest : Function to approve user request
func (b *Bassa) ApproveUserRequest(userName string) {
	if userName == "" {
		panic(errIncompleteParams)
	}
	endpoint := "/api/user/approve"
	apiURL := b.apiURL + endpoint + "/" + userName

	request, err := http.NewRequest("POST", apiURL, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("token", b.token)
	response, err := b.httpClient.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	var r interface{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		panic(err)
	}
	out, err := prettyjson.Marshal(r)
	fmt.Println(string(out))
}

// GetBlockedUserRequests : Function to get blocked user requests
func (b *Bassa) GetBlockedUserRequests() string {

	endpoint := "/api/user/blocked"
	apiURL := b.apiURL + endpoint

	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("token", b.token)
	response, err := b.httpClient.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	var r interface{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		panic(err)
	}
	out, err := prettyjson.Marshal(r)
	fmt.Println(string(out))
	return string(out)
}

// BlockUserRequest : Function to block user request
func (b *Bassa) BlockUserRequest(userName string) {
	if userName == "" {
		panic(errIncompleteParams)
	}
	endpoint := "/api/user/blocked"
	apiURL := b.apiURL + endpoint + "/" + userName

	request, err := http.NewRequest("POST", apiURL, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("token", b.token)
	response, err := b.httpClient.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	var r interface{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		panic(err)
	}
	out, err := prettyjson.Marshal(r)
	fmt.Println(string(out))
}

// UnBlockUserRequest : Function to unblock user request
func (b *Bassa) UnBlockUserRequest(userName string) {
	if userName == "" {
		panic(errIncompleteParams)
	}
	endpoint := "/api/user/blocked"
	apiURL := b.apiURL + endpoint + "/" + userName

	request, err := http.NewRequest("DELETE", apiURL, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("token", b.token)
	response, err := b.httpClient.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	var r interface{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		panic(err)
	}
	out, err := prettyjson.Marshal(r)
	fmt.Println(string(out))
}

// GetDownloadUserRequests : Function to get download user requests
func (b *Bassa) GetDownloadUserRequests(limit int) string {
	if limit == 0 {
		limit = 1
	}
	endpoint := "/api/user/downloads"
	apiURL := b.apiURL + endpoint + "/" + string(limit)

	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("token", b.token)
	response, err := b.httpClient.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	var r interface{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		panic(err)
	}
	out, err := prettyjson.Marshal(r)
	fmt.Println(string(out))
	return string(out)
}

// GetToptenHeaviestUsers : Function to get top ten heaviest users
func (b *Bassa) GetToptenHeaviestUsers() string {

	endpoint := "/api/user/heavy"
	apiURL := b.apiURL + endpoint

	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("token", b.token)
	response, err := b.httpClient.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	var r interface{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		panic(err)
	}
	out, err := prettyjson.Marshal(r)
	fmt.Println(string(out))
	return string(out)
}
