package engine

type ProcessEventParameter struct {
	Name          string `json:"name"`
	ParameterType string `json:"type"`
	Value         string `json:"value"`
}

type ProcessEvent struct {
	EventType  string                  `json:"type"`
	Name       string                  `json:"name"`
	Parameters []ProcessEventParameter `json:"parameters"`
}

type ProcessEventsPrePost struct {
	Pre  []ProcessEvent `json:"pre"`
	Post []ProcessEvent `json:"post"`
}

type ProcessTransaction struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	To          string               `json:"to"`
	Events      ProcessEventsPrePost `json:"events"`
}

type TimeSpan struct {
	Unit  string `json:"unit"`
	Value int    `json:"value"`
}

type ProcessTrigger struct {
	Name     string         `json:"name"`
	TimeSpan TimeSpan       `json:"timeSpan"`
	Events   []ProcessEvent `json:"event"`
}

type ProcessNode struct {
	Name     string               `json:"name"`
	NodeType string               `json:"version"`
	Events   ProcessEventsPrePost `json:"events"`
}

type ProcessMap struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Nodes       []ProcessNode `json:"nodes"`
}
