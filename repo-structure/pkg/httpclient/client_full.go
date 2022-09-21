package httpclient

// import (
// 	{{if eq .Auth.Type "Basic Authentication"}}
// 	"encoding/base64"
// 	{{end}}
// 	"io/ioutil"
// 	"net/http"
// 	"net/url"
// 	// "fmt"
// 	// "strconv"
// )

// {{if eq .Auth.Type "Basic Authentication"}}
// func basicAuth(username string, password string) string {
// 	auth := username + ":" + password
// 	return base64.StdEncoding.EncodeToString([]byte(auth))
// }
// {{end}}

// func makeHttpReq(username string, password string, req *http.Request, params url.Values) []byte {
// 	{{if eq .Auth.Type "Basic Authentication"}}
// 	req.Header.Add("Authorization", "Basic "+basicAuth(username, password))
// 	{{end}}
// 	req.URL.RawQuery = params.Encode()
// 	client := &http.Client{}

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil
// 	}

// 	defer resp.Body.Close()
// 	bodyBytes, _ := ioutil.ReadAll(resp.Body)

// 	return bodyBytes
// }

// {{range .Endpoints}}
// {{if eq .Name "account"}}
// func GetAccount(url string, username string, password string, params url.Values) interface{} {
// 	req, err := http.NewRequest("POST", url+"{{ .Path }}", nil)
// 	if err != nil {
// 		return nil
// 	}

// 	res := makeHttpReq(username, password, req, params)
// 	return res
// }
// {{end}}
// {{end}}

// {{range .Endpoints}}
// {{ if gt (len .PathParameters) 0 }}
// func {{ .Name | Title }}(vars map[string]string, config *Config, params url.Values) interface{} {
// {{ else if or (eq .Method "POST") (eq .Method "PUT") }}
// func {{ .Name | Title }}(config *Config, body interface{}, params url.Values) interface{} {
// {{else}}
// func {{ .Name | Title }}(config *Config, params url.Values) interface{} {
// {{end}}
// 	{{ if gt (len .PathParameters) 0 }}
// 	req, err := http.NewRequest("{{ .Method }}", "{{ $.Auth.Domain }}"+"{{ .ApiPath }}", nil)
// 	{{else if or (eq .Method "POST") (eq .Method "PUT") }}
// 	jsonReq, err := json.Marshal(body)
// 	if err != nil {
// 		return "", err
// 	}
// 	req, err := http.NewRequest("{{ .Method }}", "{{ $.Auth.Domain }}"+"{{ .ApiPath }}", nil)
// 	{{else}}
// 	req, err := http.NewRequest("{{ .Method }}", "{{ $.Auth.Domain }}"+"{{ .ApiPath }}", nil)
// 	{{end}}
// 	req.Header.Set("Content-Type", "application/json; charset=utf-8")

// 	return string(makeHttpReq(username, password, req, params)), nil
// }
// {{end}}
