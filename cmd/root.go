// Copyright © 2016 Anirudh Ramanathan anirudh4444@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	neturl "net/url"
	"os"
	"regexp"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var cfgFile string
var depth int
var filterEx string
var outputEx string
var seenUrls map[string]bool = make(map[string]bool)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "crawler <url>",
	Short: "Web page crawl and scrape utility",
	Long:  `Crawls web pages to a particular depth and allows filtering by a particular regex, and a depth.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("See options by running \"crawler --help\"")
			return
		}
		doCrawl([]string{args[0]}, depth)
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scraper.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().IntVar(&depth, "depth", 1, "Depth of crawl.")
	RootCmd.Flags().StringVar(&filterEx, "urlregex", "", "URL regex to filter how we crawl.")
	RootCmd.Flags().StringVar(&outputEx, "outregex", "^[^\\d|\\D]$", "Regex to capture items from page.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".scraper") // name of config file (without extension)
	viper.AddConfigPath("$HOME")    // adding home directory as first search path
	viper.AutomaticEnv()            // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func doCrawl(url []string, depth int) {
	if depth == 0 {
		return
	}

	var newUrls []string
	for _, u := range url {
		fmt.Println("Fetching URL: ", u)
		page, err := getPage(u)
		if err != nil {
			glog.Warningf("Page fetch failed :", u)
		} else {
			newUrls = append(newUrls, crawlPage(page, u)...)
		}
	}
	doCrawl(newUrls, depth-1)
}

func crawlPage(page string, url string) []string {
	rFull := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
	rProtocol := regexp.MustCompile(`https?:\/\/.*?`)
	rShort := regexp.MustCompile(`href=["'](.*?)["']`)
	rOutput := regexp.MustCompile(outputEx)

	u, err := neturl.Parse(url)
	var baseurl string
	if err == nil {
		baseurl = u.Scheme + "://" + u.Host
	}
	var lidx int = strings.LastIndex(url, "/")
	trimmed_url := url[:lidx]
	trimmed_url = strings.TrimRight(trimmed_url, "/")

	matches := rFull.FindAllString(page, -1)
	stubs := rShort.FindAllStringSubmatch(page, -1)
	for _, s := range stubs {
		if rProtocol.MatchString(s[1]) {
			matches = append(matches, s[1])
		} else if s[1][0] == '/' {
			trimmed := strings.TrimLeft(s[1], "/")
			matches = append(matches, (baseurl + "/" + trimmed))
		} else {
			trimmed := strings.TrimLeft(s[1], "/")
			matches = append(matches, (trimmed_url + "/" + trimmed))
		}

	}


	m := regexp.MustCompile(filterEx)
	var matches_filtered []string
	for _, v := range matches {
		if _, ok := seenUrls[v]; !ok && m.MatchString(v) {
			matches_filtered = append(matches_filtered, v)
			seenUrls[v] = true
		}
	}

	output := rOutput.FindAllString(page, -1)
	fmt.Println(output)
	return matches_filtered
}

func getPage(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	s := buf.String()
	return s, nil
}
