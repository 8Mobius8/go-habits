// Copyright © 2018 Mark Odell <m4rk.odell@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(isupCmd)
}

var isupCmd = &cobra.Command{
	Use:   "isup",
	Short: "Check if Habitica api is reachable.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(PrintHabiticaAPIStatus())
	},
}

func PrintHabiticaAPIStatus() string {
	res, err := http.Get("https://habitica.com/api/v3/status")

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var isUp StatusResponse
	err = json.Unmarshal(body, &isUp)
	if err != nil {
		log.Fatal(err)
	}

	if isUp.Data.Status == "up" {
		return "Habitica is reachable, GO catch all those pets!"
	} else {
		return ":( Habitica is unreachable."
	}
}

type StatusResponse struct {
	Success bool
	Data    struct {
		Status string
	}
}
