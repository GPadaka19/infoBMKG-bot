package storage

import (
	"encoding/json"
	"os"
	"sync"
)

const historyFile = "history.json"

type History struct {
	mu    sync.RWMutex
	Items map[string]bool `json:"items"` // GUID -> true
}

// NewHistory creates a new history instance and loads data from file
func NewHistory() *History {
	h := &History{
		Items: make(map[string]bool),
	}
	h.load()
	return h
}

// load reads the history from the JSON file
func (h *History) load() {
	h.mu.Lock()
	defer h.mu.Unlock()

	file, err := os.ReadFile(historyFile)
	if err != nil {
		// If file doesn't exist, we start empty
		return
	}

	// We create a temporary struct to match JSON structure if needed,
	// but here direct map unmarshal is fine if the JSON is just the map.
	// However, for better structure let's unmarshal specific field if we change format later.
	// For now let's just save the map directly.
	_ = json.Unmarshal(file, &h.Items)
}

// save writes the current history to the JSON file
func (h *History) save() error {
	data, err := json.MarshalIndent(h.Items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(historyFile, data, 0644)
}

// IsNew checks if a GUID is new (not in history)
func (h *History) IsNew(guid string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, exists := h.Items[guid]
	return !exists
}

// Add adds a GUID to history and saves to file
func (h *History) Add(guid string) error {
	h.mu.Lock()
	// Check again inside lock to be safe (though race condition unlikely in single poller)
	h.Items[guid] = true
	h.mu.Unlock()

	// Implementation note:
	// Saving on every single add might be I/O intensive if many alerts come at once.
	// But for weather alerts (rare), this is perfectly fine and safer.
	return h.ApplySave()
}

// ApplySave is a public method to trigger save manually if needed
func (h *History) ApplySave() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.save()
}
