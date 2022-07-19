package wikipedia

// MIT License
//
// Copyright (c) 2021 Damian Zaremba
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func getLoginToken(httpClient *http.Client) (string, error) {
	req, err := http.NewRequest("GET", "https://en.wikipedia.org/w/api.php?action=query&meta=tokens&type=login&format=json", nil)
	if err != nil {
		return "", nil
	}
	req.Header.Set("User-Agent", "ClueBot NG Review NG/1.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}
	if err := resp.Body.Close(); err != nil {
		return "", nil
	}

	data := struct {
		Query struct {
			Tokens struct {
				LoginToken string `json:"logintoken"`
			} `json:"tokens"`
		} `json:"query"`
	}{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", nil
	}
	return data.Query.Tokens.LoginToken, nil
}

func getCsrfToken(httpClient *http.Client) (string, error) {
	req, err := http.NewRequest("GET", "https://en.wikipedia.org/w/api.php?action=query&meta=tokens&format=json", nil)
	if err != nil {
		return "", nil
	}
	req.Header.Set("User-Agent", "ClueBot NG Review NG/1.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}
	if err := resp.Body.Close(); err != nil {
		return "", nil
	}

	data := struct {
		Query struct {
			Tokens struct {
				CsrfToken string `json:"csrftoken"`
			} `json:"tokens"`
		} `json:"query"`
	}{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", nil
	}
	return data.Query.Tokens.CsrfToken, nil
}

func login(httpClient *http.Client, username, password, token string) error {
	form := url.Values{}
	form.Add("action", "login")
	form.Add("lgname", username)
	form.Add("lgpassword", password)
	form.Add("lgtoken", token)
	form.Add("format", "json")

	req, err := http.NewRequest("POST", "https://en.wikipedia.org/w/api.php", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "ClueBot NG Review NG/1.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}

	data := struct {
		Login struct {
			Result string
		}
	}{}
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	if data.Login.Result != "Success" {
		return fmt.Errorf("API error: %s", body)
	}
	return nil
}

func UpdatePage(contents string) error {
	httpClient := &http.Client{}
	return updatePage(httpClient, contents)
}

func UpdatePageWithCredentials(contents, username, password string) error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	httpClient := &http.Client{Jar: jar}

	loginToken, err := getLoginToken(httpClient)
	if err != nil {
		panic(err)
	}

	if err := login(httpClient, username, password, loginToken); err != nil {
		panic(err)
	}

	return updatePage(httpClient, contents)
}

func updatePage(httpClient *http.Client, contents string) error {
	csrfToken, err := getCsrfToken(httpClient)
	if err != nil {
		return nil
	}
	form := url.Values{}
	form.Add("action", "edit")
	form.Add("title", "User:ClueBot NG/ReviewInterface/Stats")
	form.Add("summary", "Uploading Stats")
	form.Add("token", csrfToken)
	form.Add("format", "json")
	form.Add("text", contents)

	req, err := http.NewRequest("POST", "https://en.wikipedia.org/w/api.php", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "ClueBot NG Review NG/1.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}

	data := struct {
		Edit struct {
			Result string
			PageId int `json:"pageid"`
			Title  string
			OldId  int `json:"oldrevid"`
			NewId  int `json:"newrevid"`
		}
	}{}
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	if data.Edit.Result != "Success" {
		return fmt.Errorf("API error: %s", body)
	}
	return nil
}
