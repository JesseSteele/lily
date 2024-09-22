// For YAML configs
// https://x.com/i/grok/share/ybc59WnxQoaxwZb7V6k1353yx

package main

import (
    "flag"
    "fmt"
    "github.com/spf13/viper"
    "log"
    "net/http"
    "os"
    "path/filepath"
)

func main() {

    // Config
    configFile := flag.String("config", "/etc/lily/master.conf", "Path to the configuration file")
    flag.Parse()

    viper.SetConfigFile(*configFile)
    err := viper.ReadInConfig()
    if err != nil {
        log.Fatalf("Error reading config file: %s", err)
    }

    // Start server
    server := &http.Server{
        Addr:    fmt.Sprintf(":%d", viper.GetInt("port")),
        Handler: setupRoutes(),
    }

    log.Println("Starting server...")
    log.Fatal(server.ListenAndServe())
}

func setupRoutes() http.Handler {
    // Placeholder for routing setup
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from Lily Server!")
    })
}

// For BASH only, does not include case switch for Go, Node, or Python
func execScript(scriptPath string, env map[string]string) ([]byte, error) {
    cmd := exec.Command("/bin/bash", scriptPath)
    cmd.Env = []string{}
    for k, v := range env {
        cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
    }
    output, err := cmd.CombinedOutput()
    return output, err
}
