package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	"github.com/golang/glog"
	"github.com/zenazn/goji"
	gojiweb "github.com/zenazn/goji/web"
)

type lines []string

var dictionary map[string]string

func execute(scriptPath string) (string, error) {
	out, err := exec.Command(scriptPath).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func executeShellScriptHandler(c gojiweb.C, w http.ResponseWriter, r *http.Request) {
	actionName := c.URLParams["action_name"]

	scriptPath, ok := dictionary[actionName]
	if !ok {
		glog.Errorf("Error[%s] action[%s] from dictionary[%s] not found", "", actionName, dictionary)
		http.Error(w, fmt.Sprintf("Error[%s] action[%s] from dictionary[%s] not found", "", actionName, dictionary), http.StatusBadRequest)
		return
	}

	scriptContent, err := ioutil.ReadFile(scriptPath)
	if err != nil {
		glog.Errorf("Error[%s] get script content[%s] for action[%s] from dictionary[%s]", err, scriptContent, actionName, dictionary)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := execute(scriptPath)
	if err != nil {
		glog.Errorf("Error[%s] execute script content[%s] for action[%s] from dictionary[%s], cmd output[%s]", err, scriptContent, actionName, dictionary, out)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, _ := json.Marshal(map[string]string{
		"ok":  "1",
		"out": string(out),
	})

	fmt.Fprintf(w, "%s", string(result))
}

func main() {
	defer glog.Flush()
	glog.Info("Initializing shellscript executor server")

	port := flag.String("P", "6304", "port")
	flag.Parse()

	flag.Lookup("logtostderr").Value.Set("true")

	EnvSettingsInit()
	//read dictionary into memory
	dicData, err := ioutil.ReadFile(ProjectEnvSettings.DictionaryPath)
	if err != nil {
		FailOnError(err, fmt.Sprintf("Error[%s] reading dictionary from path[%]", err, ProjectEnvSettings.DictionaryPath))
	}

	glog.Infof("Dictionary[%s]", dicData)

	if err := json.Unmarshal(dicData, &dictionary); err != nil {
		FailOnError(err, fmt.Sprintf("Error[%s] decode json dictionary from path[%] with content[%s]", err, ProjectEnvSettings.DictionaryPath, dicData))
	}

	goji.Get("/shell/execute/:action_name", executeShellScriptHandler)

	flag.Set("bind", ":"+*port)
	goji.Serve()
}

func FailOnError(err error, msg string) {
	if err != nil {
		glog.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
