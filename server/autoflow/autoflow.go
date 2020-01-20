package autoflow

import (
    "fmt"
    "log"
    "encoding/json"
    "io/ioutil"
    "time"
    "github.com/google/uuid"
    "github.com/go-redis/redis/v7"
)


type Autoflow struct {
    redisClient *redis.Client
}

func (f *Autoflow) CreateSession(flowName string) string {
    log.Printf("CreateSession for flow %s\n", flowName)

    // load data from file
    file, _ := ioutil.ReadFile(fmt.Sprintf("./data/%s.json", flowName))
    data := AutoflowData{}
    loadDataErr := json.Unmarshal([]byte(file), &data)
    if loadDataErr != nil {
        log.Printf(loadDataErr.Error())
        return ""
    }

    for _, step := range data.Steps {
        stepBytes, _ := json.Marshal(step)
        value := string(stepBytes)
        key := string(step.Id)
        f.redisClient.HSet(flowName, key, value)
    }

    id, _ := uuid.NewUUID()
    sessionId := id.String()
    
    f.redisClient.HSet(sessionId, "flowname", flowName)
    f.setCurrentFlowStepId(sessionId, "0")

    f.getFlowStep(flowName, "0")

    timeout := time.Duration(3600)*time.Second
    f.redisClient.Expire(sessionId, timeout)
    f.redisClient.Expire(flowName, timeout)

    return sessionId
}

func (f *Autoflow) QueryNextStep(sessionId string, clientParams map[string]interface{}) (string, ActionParameter) {
    log.Printf("QueryStep for session %s\n", sessionId)
    flowname, currentStep, succeed := f.getCurrentStep(sessionId)
    if !succeed {
        return "", nil
    }

    nextStep, succeed := f.getFlowStep(flowname, currentStep.Next)
    if !succeed {
        return "", nil
    }

    f.setCurrentFlowStepId(sessionId, nextStep.Id)

    if len(nextStep.ServerAction) > 0 {
        // perform serverAction if exists
        if !invokeServerAction(nextStep.ServerAction, clientParams, nextStep.Parameter) {
            log.Printf("Error: server side action failed")
            return f.OnError(sessionId)
        }
        return nextStep.Action, nil
    }

    return nextStep.Action, nextStep.Parameter
}

func (f *Autoflow) OnError(sessionId string) (string, ActionParameter) {
    log.Printf("Received Error for session %s\n", sessionId)

    flowname, currentStep, succeed:= f.getCurrentStep(sessionId)
    if !succeed {
        return "", nil
    }

    if currentStep.OnError == "" {
        log.Printf("no error handling for current step %s. stop execution", currentStep.Id)
        return "", nil
    }

    nextStep, succeed := f.getFlowStep(flowname, currentStep.OnError)
    if !succeed {
        return "", nil
    }

    f.setCurrentFlowStepId(sessionId, nextStep.Id)

    return nextStep.Action, nextStep.Parameter
}

func (f *Autoflow) Stop(sessionId string) {
    log.Printf("stop session %s\n", sessionId)
    f.redisClient.Del(sessionId)
    return
}


func (f *Autoflow) getCurrentStep(sessionId string) (string, FlowStep, bool) {
    flowstep := FlowStep{}

    currentStepId, succeed := f.getCurrentFlowStepId(sessionId)
    if !succeed {
        return "", flowstep, false
    }

    log.Printf("current step id: %s\n", currentStepId)
    flowname := f.redisClient.HGet(sessionId, "flowname").Val()
    currentStep, succeed := f.getFlowStep(flowname, currentStepId)
    if !succeed {
        return flowname, flowstep, false
    }

    return flowname, currentStep, true
}


func (f *Autoflow) getCurrentFlowStepId(sessionId string) (string, bool) {
    log.Printf("get current step id for session: %s\n", sessionId)
    value := f.redisClient.HGet(sessionId, "current_id")
    if value.Err() != nil {
        log.Printf("Error: %s\n", value.Err().Error())
        return "", false
    }

    return value.Val(), true
}

func (f *Autoflow) setCurrentFlowStepId(sessionId string, stepId string) bool {
    value := f.redisClient.HSet(sessionId, "current_id", stepId)
    if value.Err() != nil {
        log.Printf("Error: %s\n", value.Err().Error())
        return false
    }

    return true
}


func (f *Autoflow) getFlowStep(flowName string, stepId string) (FlowStep, bool) {
    log.Printf("getting step %s from %s\n", stepId, flowName)
    flowStep := FlowStep{}    
    stepData := f.redisClient.HGet(flowName, stepId)
    if stepData.Err() != nil {
        log.Printf("Error fetching redis key (%s, %s): %v\n", flowName, stepId, stepData.Err())
        return flowStep, false
    }

    err := json.Unmarshal([]byte(stepData.Val()), &flowStep)
    if err != nil {
        log.Printf(err.Error())
        return flowStep, false
    }

    return flowStep, true
}


func NewAutoflow(redisAddr string) Autoflow {
    flow := Autoflow{}
    flow.redisClient = redis.NewClient(&redis.Options{
                            Addr:     redisAddr,
                            Password: "", 
                            DB:       0, 
                        })
    return flow
}

var autoflowInstance = NewAutoflow("127.0.0.1:6379")

func GetInstance() *Autoflow {
    return &autoflowInstance
}