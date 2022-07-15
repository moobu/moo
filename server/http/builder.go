package http

import (
	"encoding/json"
	"net/http"

	"github.com/moobu/moo/builder"
)

type BuildArgs struct {
	Source  *builder.Source
	Options *builder.BuildOptions
}

type CleanArgs struct {
	Bundle  *builder.Bundle
	Options *builder.CleanOptions
}

func Build(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	args := BuildArgs{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&args); err != nil {
		writeJSON(w, nil, err)
		return
	}
	defer r.Body.Close()

	opts := args.Options
	bundle, err := builder.Build(args.Source,
		builder.BuildWithNamespace(opts.Namespace),
		builder.Dir(opts.Dir),
		builder.Ref(opts.Ref))

	writeJSON(w, bundle, err)
}

func Clean(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	args := CleanArgs{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&args); err != nil {
		writeJSON(w, nil, err)
		return
	}
	defer r.Body.Close()

	opts := args.Options
	err := builder.Clean(args.Bundle,
		builder.CleanWithNamespace(opts.Namespace))

	writeJSON(w, nil, err)
}
