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
	Type       string   `json:"type,omitempty"`
	Domain     string   `json:"domain,omitempty"`
	Spec       []string `json:"spec,omitempty"`
	MeEndpoint string   `json:"me_endpoint,omitempty"`
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
	EnableAllParameters bool        `json:"enableAllParameters,omitempty"`
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

func removeUnderscores(s string) string {
	var fragments []string
	for _, frag := range strings.Split(s, "_") {
		fragments = append(fragments, strings.Title(frag))
	}
	return strings.Join(fragments, "")
}

func insertPathParameters(path string) {
	// r := regexp.MustCompile(`\{(.*?)\}`)
	// matches := r.FindAllString(path, -1)
	// for _, match := range matches {
	// 	match := strings.Trim(match, "{")
	// 	trimmed := strings.Trim(match, "}")
	// 	re := regexp.MustCompile(`\{` + trimmed + `\}`)
	// 	path = re.ReplaceAllString(path, vars[trimmed])
	// }
	// fmt.Println(path)

	r := regexp.MustCompile(`\{(.*?)\}`)
	matches := r.FindAllString(path, -1)

	for _, match := range matches {
		m := regexp.MustCompile(match)
		i := m.FindStringIndex(path)
		fmt.Println(i)
	}

	//Note: return a list with each substring
	//before match
	//match replaced with

	//note to self: just do id
	//split path into "pre-id" and "post-id"
	//then return a list of strings
	//make it so it's "pre_id" + id + "post_id"

	//maybe seperate function for pre_id
	//seperate function for post id to make it easy to do

	//{{pre_id function}} + id + {{post_id_function}}
}

func populateTemplateFiles(source_dir string, dest_dir string, config *Config) {
	//Make destination directory
	dest_dir = "../" + dest_dir
	os.MkdirAll(dest_dir, 0755)

	//Add custom template functions
	funcMap := template.FuncMap{
		"ToUpper":           strings.ToUpper,
		"Title":             strings.Title,
		"RemoveUnderscores": removeUnderscores,
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
				buf, _ := ioutil.ReadFile(path)
				tmpl, _ := template.New("connectorFiles").Funcs(funcMap).Parse(string(buf))

				//Generate target file
				targetFile := filepath.Join(dest_dir, fileBase)
				f, err := os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE, 0755)
				if err != nil {
					return err
				}

				//Write to target file
				w := bufio.NewWriter(f)
				tmpl.Execute(w, config)
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
	// insertPathParameters("/projects/api/{version}/projects/{id}.json")
}
