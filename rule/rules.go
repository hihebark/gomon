package rule

import "encoding/json"

// Rules set rules to be worked with when executing gomon.
// this rules can be executed from a file default file is: gomon.json
type Rule struct {
	Restartable string   `json:"restartable"` // Default rs.
	Ignore      []string `json:"ignore"`      // Dirs to be ignored.
	Verbose     bool     `json:"verbose"`     // Verbose stdout.
	ExecCommand string   `json:"execCommand"` // First command to be executed default: go run main.go.
	Events      Event    `json:"events"`      // Set Event(command) on [start, restart, exit].
	Watch       string   `json:"watch"`       // Dirs to be watched.
	Ext         []string `json:"ext"`         // Extension to be watched.
}

// Event is what to be executed in one of this events
type Event struct {
	OnStart   string `json:"onStart"`   // On start.
	OnRestart string `json:"onRestart"` // On restart.
	OnExit    string `json:"onExit"`    // On exit.
}

// NewRules create a new Rules.
func NewRule() *Rule {
	return &Rule{
		Restartable: "rs",
		Ignore:      []string{"vendor", ".git"},
		Verbose:     true,
		ExecCommand: "go run main.go",
		Watch:       "*.*",
		Ext:         []string{"go"},
		Events:      Event{},
	}
}
func NewRuleFromFile(body string) (*Rule, error) {
	rule := &Rule{}
	err := json.Unmarshal([]byte(body), rule)
	if err != nil {
		return nil, err
	}
	return rule, nil
}
