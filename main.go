package main

import (
	"path/filepath"
	"os"
	"text/template"
	"strings"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"bufio"

	// conf "github.com/kosha/kosha-cdk/config"
	// "github.com/kosha/kosha-cdk/utils"
)

var nonGoFiles = []string{"/go.mod"}

type Config struct {
	ConnectorName string `json:"connector_name,omitempty"`
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

func contains(elems []string, v string) bool {
    for _, s := range elems {
        if v == s {
            return true
        }
    }
    return false
}

func populateTemplateFiles(source_dir string, dest_dir string, config *Config) {
	//Make destination directory
	dest_dir = "../" + dest_dir
	os.MkdirAll(dest_dir, 0755)

	err := filepath.Walk("./" + source_dir, func(path string, info os.FileInfo, err error) error {
		fileInfo, _ := os.Stat(path)
		re := regexp.MustCompile(source_dir + "\\s*(.*?)$")
		fileBase := re.FindStringSubmatch(path)[1]
		fmt.Println(fileBase)
		//If it is a directory, copy it
		if fileInfo.IsDir() {
			fmt.Println("I AM IN IS DIR" + path)
			os.MkdirAll(filepath.Join(dest_dir, fileBase), 0755)
		} else {
			fmt.Println("I AM NOT" + path)
			//If it is a template file, read and populate tags
			if strings.Contains(path, ".tmpl") {
				re = regexp.MustCompile(source_dir + "\\s*(.*?)\\s*.tmpl")
				fileBase = re.FindStringSubmatch(path)[1]
				fmt.Println(fileBase)
				//Read in file and execute template
				buf, _ := ioutil.ReadFile(path)
				tmpl, _ := template.New("connectorFiles").Parse(string(buf))
				
				//Generate target file
				targetFile := filepath.Join(dest_dir, fileBase)
				// if ! utils.contains(conf.nonGoFiles, fileBase) {
				// 	targetFile = targetFile + ".go"
				// }
				if ! (contains(nonGoFiles, fileBase)) {
					targetFile = targetFile + ".go"
				}
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
}














// func ParseTemplates() []string {
// 	var filepaths []string
//     err := filepath.Walk("./repo-structure", func(path string, info os.FileInfo, err error) error {
//         if strings.Contains(path, ".tmpl") {
//             filepaths = append(filepaths, path)
//             if err != nil {
//             	fmt.Println(err)
//             }
//         }
//     })

//     if err != nil {
//         panic(err)
//     }

//     return filepaths
// }

// func ParseDirectory(baseDir string) {
// 	t, err := template.New("templ").ParseGlob(path + "/templates/*.html")
// }


// func main() {
// 	connectorData := Input{"teamwork"}
// 	filePaths := ParseTemplates()
// 	for _, file := range filePaths {

// 	}

//     t := template.New("main.tmpl").ParseFiles(filePaths)
// 	err := t.Execute(os.Stdout, connectorData)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func main() {
// 	connectorData := Input{"teamwork"}
// 	if err := os.Mkdir("connector", os.ModePerm); err != nil {
//         log.Fatal(err)
//     }
// 	err := filepath.Walk("./repo-structure", func(path string, info os.FileInfo, err error) error {
//         file, err := os.Open(path)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fileInfo, err := file.Stat()
// 		if err != nil {
// 			panic(err)
// 		}
// 		if fileInfo.IsDir() {

// 		}
//     })
// }
