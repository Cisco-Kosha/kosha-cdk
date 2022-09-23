package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	// "golang.org/x/text/cases"
	// "golang.org/x/text/language"
)

// var nonGoFiles = []string{"/go.mod"}

type Config struct {
	ConnectorName string      `json:"connector_name,omitempty"`
	Auth          Auth        `json:"auth,omitempty"`
	Prometheus    Prometheus  `json:"prometheus,omitempty"`
	Endpoints     []Endpoints `json:"endpoints,omitempty"`
}

type Auth struct {
	Type       string     `json:"type,omitempty"`
	Domain     string     `json:"domain,omitempty"`
	Spec       []string   `json:"spec,omitempty"`
	MeEndpoint MeEndpoint `json:"me_endpoint,omitempty"`
}

type MeEndpoint struct {
	Path   string `json:"path,omitempty"`
	Method string `json:"method,omitempty"`
}

type Prometheus struct {
	Enabled string `json:"enabled,omitempty"`
}

type Endpoints struct {
	Name                string      `json:"name,omitempty"`
	Description         string      `json:"description,omitempty"`
	ApiPath             string      `json:"api_path,omitempty"`
	Method              string      `json:"method,omitempty"`
	KoshaPath           string      `json:"kosha_path,omitempty"`
	QueryParameters     []string    `json:"query_parameters,omitempty"`
	PathParameters      []string    `json:"path_parameters,omitempty"`
	EnableAllParameters bool        `json:"enable_all_parameters,omitempty"`
	Tag                 string      `json:"tag,omitempty"`
	Body                interface{} `json:"body,omitempty"`
}

func copyFile(source_path string, dest_path string) {
	bytesRead, err := ioutil.ReadFile(source_path)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(dest_path, bytesRead, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func leadingPath(s string) string {
	splitted := strings.Split(s, "{")
	return splitted[0]
}

func trailingPath(path string) string {
	splitted := strings.Split(path, "}")
	return splitted[1]
}

func isPathParam(path string) bool {
	return strings.Contains(path, "{")
}

func removeUnderscores(s string) string {
	var fragments []string
	for _, frag := range strings.Split(s, "_") {
		fragments = append(fragments, strings.Title(frag))
	}
	return strings.Join(fragments, "")
}

func populateTemplateFiles(source_dir string, dest_dir string, config *Config) {
	//Make destination directory
	dest_dir = "../" + dest_dir
	os.MkdirAll(dest_dir, 0755)

	//Add custom template functions
	funcMap := template.FuncMap{
		"ToUpper":           strings.ToUpper,
		"ToLower":           strings.ToLower,
		"Title":             strings.Title,
		"LeadingPath":       leadingPath,
		"RemoveUnderscores": removeUnderscores,
		"TrailingPath":      trailingPath,
		"isPathParam":       isPathParam,
	}

	err := filepath.Walk("./"+source_dir, func(path string, info os.FileInfo, err error) error {
		fileInfo, _ := os.Stat(path)
		re := regexp.MustCompile(source_dir + "\\s*(.*?)$")
		fileBase := re.FindStringSubmatch(path)[1]
		//If it is a directory, copy it
		if fileInfo.IsDir() {
			os.MkdirAll(filepath.Join(dest_dir, fileBase), 0755)
		} else {
			//If it is a template file, read and populate tags
			if strings.Contains(path, ".tmpl") {
				re = regexp.MustCompile(source_dir + "\\s*(.*?)\\s*.tmpl")
				fileBase = re.FindStringSubmatch(path)[1]
				//Read in file and execute template
				buf, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				tmpl, err := template.New("connectorFiles").Funcs(funcMap).Parse(string(buf))
				if err != nil {
					return err
				}

				//Generate target file
				targetFile := filepath.Join(dest_dir, fileBase)
				f, err := os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE, 0755)
				if err != nil {
					return err
				}

				//Write to target file
				w := bufio.NewWriter(f)
				err = tmpl.Execute(w, config)
				if err != nil {
					return err
				}
				w.Flush()
			} else {
				//Not a template file -- copy as is
				copyFile(path, filepath.Join(dest_dir, fileBase))
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	var config Config
	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(configFile, &config)
	destination_dir := config.ConnectorName + "-connector"
	populateTemplateFiles("repo-structure", destination_dir, &config)
	// fmt.Println(getPostIdString("/projects/{id}.json"))
}
