package http

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/moobu/moo/runtime"
)

type CreateArgs struct {
	Pod     *runtime.Pod
	Options *runtime.CreateOptions
}

func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	args := CreateArgs{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&args); err != nil {
		writeJSON(w, nil, err)
		return
	}
	defer r.Body.Close()

	// TODO: create a file for this session to output the log to
	options := args.Options
	opts := []runtime.CreateOption{
		runtime.Output(os.Stdout),
		runtime.Args(options.Args...),
		runtime.Env(options.Env...),
		runtime.Replicas(options.Replicas),
		runtime.Bundle(options.Bundle),
		runtime.GPU(options.GPU),
		runtime.CreateWithNamespace(options.Namespace),
	}
	writeJSON(w, nil, runtime.Create(args.Pod, opts...))
}

type DeleteArgs struct {
	Pod     *runtime.Pod
	Options *runtime.DeleteOptions
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	args := &DeleteArgs{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(args); err != nil {
		writeJSON(w, nil, err)
		return
	}
	defer r.Body.Close()

	options := args.Options
	err := runtime.Delete(args.Pod, runtime.DeleteWithNamespace(options.Namespace))
	writeJSON(w, nil, err)
}

func List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	options := runtime.ListOptions{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&options); err != nil {
		writeJSON(w, nil, err)
		return
	}
	defer r.Body.Close()

	opts := []runtime.ListOption{
		runtime.Name(options.Name),
		runtime.Tag(options.Tag),
		runtime.ListWithNamespace(options.Namespace),
	}
	pods, err := runtime.List(opts...)
	writeJSON(w, pods, err)
}
