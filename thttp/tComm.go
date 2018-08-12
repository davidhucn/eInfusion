package thttp

type reqMsg struct {
	ID     string `json:"ID"`
	Action string `json:"Action"`
	// Action string `json:"-"`
}

var cliMsg []reqMsg
