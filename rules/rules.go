package rules

// Rules set rules to be worked with when executing gomon.
// this rules can be executed from a file default file is: gomon.json
type Rules struct {
	Restartable string   `json:"restartable"` // Default rs.
	Ignore      []string `json:"ignore"`      // Dirs to be ignored.
	Verbose     bool     `json:"verbose"`     // Verbose stdout.
	ExecCommand string   `json:"execCommand"` // First command to be executed default: go run main.go.
	Events      []Event  `json:"events"`      // Set Event(command) on [start, restart, exit].
	Watch       []string `json:"watch"`       // Dirs to be watched.
	Ext         []string `json:"ext"`         // Extension to be watched.
}

// Event is what to be executed in one of this events
type Event struct {
	OnStart   string `json:"onStart"`   // On start.
	OnRestart string `json:"onRestart"` // On restart.
	OnExit    string `json:"onExit"`    // On exit.
}

// NewRules create a new Rules.
func NewRules() *Rules {
	return &Rules{
		Restartable: "rs",
		Ignore:      []string{"vendor", ".git"},
		Verbose:     true,
		ExecCommand: "go run main.go",
		Watch:       []string{"*"},
		Ext:         []string{"go"},
		Events:      []Event{},
	}
}
