package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mikerybka/cmd-server/pkg/util"
)

type CmdServer struct {
	Dir string
}

func (s *CmdServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := &struct {
		Dir string
		Cmd string
	}{}
	json.NewDecoder(r.Body).Decode(req)
	c := strings.Split(req.Cmd, " ")
	cmd := exec.Command(c[0], c[1:]...)
	cmd.Dir = filepath.Join(s.Dir, req.Dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("ERROR:", string(out))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := &struct {
		Output   string
		ExitCode int
	}{
		Output:   string(out),
		ExitCode: cmd.ProcessState.ExitCode(),
	}
	json.NewEncoder(w).Encode(res)
}

func main() {
	util.Serve(&CmdServer{
		Dir: util.EnvVar("DIR", "."),
	})
}
