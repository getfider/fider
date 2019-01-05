package stripe

import "encoding/json"

// Application describes the properties for an Application.
type Application struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UnmarshalJSON handles deserialization of an Application.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (a *Application) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		a.ID = id
		return nil
	}

	type application Application
	var v application
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = Application(v)
	return nil
}
