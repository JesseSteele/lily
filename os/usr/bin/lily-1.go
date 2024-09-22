// For JSON configs
// https://x.com/i/grok/share/rFGC1btuCsHHua2K1DGhCNp6I

package main

import (
    "crypto/tls"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os/exec"
    "strings"
)

type SiteConfig struct {
    Name     string `json:"name"`
    Root     string `json:"root"`
    CertFile string `json:"certFile"`
    KeyFile  string `json:"keyFile"`
    Language string `json:"language"`
}

type Config struct {
    Sites []SiteConfig `json:"sites"`
}

var config Config

func main() {
    // Config (but must be written in JSON; no such config exists at initial stage of development)
    configFile := flag.String("config", "/etc/lily/master.conf", "Path to the configuration file")
    flag.Parse()

    viper.SetConfigFile(*configFile)
    err := viper.ReadInConfig()
    if err != nil {
        log.Fatalf("Error reading config file: %s", err)
    }

    configData, err := ioutil.ReadFile("server_config.json")
    if err != nil {
        log.Fatalf("Failed to read config: %v", err)
    }
    if err := json.Unmarshal(configData, &config); err != nil {
        log.Fatalf("Failed to parse config: %v", err)
    }

    for _, site := range config.Sites {
        go func(s SiteConfig) {
            http.HandleFunc("/"+s.Name+"/", func(w http.ResponseWriter, r *http.Request) {
                handleRequest(w, r, s)
            })
            log.Println("Serving site:", s.Name)
        }(site)
    }

    // Non-SSL server for redirection or other purposes
    go http.ListenAndServe(":80", http.HandlerFunc(redirectToHTTPS))

    // Start SSL servers
    for _, site := range config.Sites {
        go startSecureServer(site)
    }

    select {} // Keep the main goroutine alive
}

func startSecureServer(site SiteConfig) {
    server := &http.Server{
        Addr: ":443",
        TLSConfig: &tls.Config{
            MinVersion: tls.VersionTLS12,
        },
    }
    log.Fatal(server.ListenAndServeTLS(site.CertFile, site.KeyFile))
}

func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}

func handleRequest(w http.ResponseWriter, r *http.Request, site SiteConfig) {
    scriptPath := site.Root + strings.TrimPrefix(r.URL.Path, "/"+site.Name+"/")
    
    var cmd *exec.Cmd
    switch site.Language {
    case "go":
        cmd = exec.Command("go", "run", scriptPath+".go")
    case "node":
        cmd = exec.Command("node", scriptPath+".js")
    case "python":
        cmd = exec.Command("python3", scriptPath+".py")
    case "bash":
        cmd = exec.Command("/bin/bash", scriptPath+".bash")
    default:
        http.Error(w, "Unsupported language", http.StatusInternalServerError)
        return
    }

    // Execute the Bash script
    cmd := exec.Command("/bin/bash", scriptPath)
    cmd.Env = append(os.Environ(), 
        fmt.Sprintf("REQUEST_METHOD=%s", r.Method),
        fmt.Sprintf("QUERY_STRING=%s", r.URL.RawQuery),
        fmt.Sprintf("CONTENT_LENGTH=%d", r.ContentLength),
        // Add more environment variables as needed
    )

    // If POST, pass the body content
    if r.Method == "POST" {
        body, _ := ioutil.ReadAll(r.Body)
        cmd.Stdin = strings.NewReader(string(body))
    }

    output, err := cmd.CombinedOutput()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error executing script: %v", err), http.StatusInternalServerError)
        return
    }

    // Write the output to the response
    w.Write(output)
}
