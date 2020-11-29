package main

import (
	"fmt"
	"github.com/AlekSi/alice"
)

type SessionState map[string]string

type ResponseWithState struct {
	// original alice.Response
	Response alice.ResponsePayload `json:"response"`
	Session  alice.ResponseSession `json:"session"`
	Version  string                `json:"version"`
	// New field
	SessionState SessionState `json:"session_state"`
}

type State struct {
	Session map[string]string `json:"session,omitempty"`
	User    map[string]string `json:"user,omitempty"`
}

type RequestWithState struct {
	// original alice.Request
	Version string               `json:"version"`
	Meta    alice.RequestMeta    `json:"meta"`
	Request alice.RequestPayload `json:"request"`
	Session alice.RequestSession `json:"session"`
	// New field
	State State `json:"state,omitempty"`
}

func audioUrlFromFileUID(id string) string {
	return fmt.Sprintf("<speaker audio=\"dialogs-upload/%s/%s.opus\">", dialog_ID, id)
}
