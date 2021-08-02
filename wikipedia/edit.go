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
	"net/url"
)

func UpdatePage(contents string) error {
	form := url.Values{}
	form.Add("action", "edit")
	form.Add("title", "User:ClueBot NG/ReviewInterface/Stats")
	form.Add("summary", "Uploading Stats")
	form.Add("token", "+\\")
	form.Add("format", "json")
	form.Add("text", contents)

	resp, _ := http.PostForm("https://en.wikipedia.org/w/api.php", form)
	body, _ := ioutil.ReadAll(resp.Body)
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
		return fmt.Errorf("API error: %s", data.Edit.Result)
	}
	return nil
}
