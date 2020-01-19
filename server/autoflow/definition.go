package autoflow

type ActionParameter = map[string]string

type FlowStep struct {
    Id string
    Action string
    Parameter ActionParameter
    Next string
    OnError string
}

type AutoflowData struct {
    Steps []FlowStep
}