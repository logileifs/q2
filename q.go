package main

import "os"
import "fmt"
import "time"
import "bytes"
import "strings"
import "net/http"
import "encoding/json"
import "github.com/fatih/color"
import "github.com/jessevdk/go-flags"
import "github.com/briandowns/spinner"

var opts struct {
    // Slice of bool will append 'true' each time the option
    // is encountered (can be set multiple times, like -vvv)
    Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`

    // Example of a required flag
    Name string `short:"n" long:"name" description:"A name" required:"false"`
}
//{
//    "id": "chatcmpl-7Fr6c6rIuaqr5hGbEcqUDi94n8kIs",
//    "object":"chat.completion",
//    "created":1684013414,
//    "model":"gpt-3.5-turbo-0301",
//    "usage":{
//        "prompt_tokens":14,
//        "completion_tokens":5,
//        "total_tokens":19
//    },
//    "choices":[
//        {
//            "message":{
//                "role":"assistant",
//                "content":"This is a test!"
//            },
//            "finish_reason":"stop",
//            "index":0,
//        }
//    ]
//}
type MessageResponse struct {
    Role string    `json:"role"`
    Content string `json:"content"`
}
type ChoiceResponse struct {
    Message MessageResponse `json:"message"`
    FinishReason  string    `json:"finish_reason"`
    Index float64           `json:"index"`
}
type UsageResponse struct {
    PromptTokens     float64 `json:"prompt_tokens"`
    CompletionTokens float64 `json:"completion_tokens"`
    TotalTokens      float64 `json:"total_tokens"`
}
type OpenAIResponse struct {
    Id string                `json:"-"`
    Object string            `json:"-"`
    Created string           `json:"-"`
    Model string             `json:"-"`
    Usage UsageResponse      `json:"-"`
    Choices []ChoiceResponse `json:"choices"`
}

type OpenAIMessageRequest struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}
type OpenAIRequest struct {
    Model       string                 `json:"model"`
    Messages    []OpenAIMessageRequest `json:"messages"`
    Temperature float64                `json:"temperature"`
}

func makeRequest(prompt string, channel chan OpenAIResponse) {
    //fmt.Printf("Prompt: %s\n", prompt)
    var API_KEY = os.Getenv("OPENAI_API_KEY")
    var openAIRequest = OpenAIRequest{
        Model: "gpt-3.5-turbo",
        Messages: []OpenAIMessageRequest{
            {
                Role:    "user",
                Content: prompt,
            },
        },
        Temperature: 0.5,
    }
    postPayload, err := json.Marshal(openAIRequest)
    //fmt.Printf("json: %s\n", postPayload)
    payload := new(OpenAIResponse)
    req, err := http.NewRequest(
        "POST",
        "https://api.openai.com/v1/chat/completions",
        bytes.NewBuffer(postPayload),
    )
    if err != nil { panic(err) }
    req.Header.Set("authorization", fmt.Sprintf("Bearer %s", API_KEY))
    req.Header.Set("content-type", "application/json")
    client := &http.Client{}
    rsp, err := client.Do(req)
    if err != nil { panic(err) }
    defer rsp.Body.Close()

    json.NewDecoder(rsp.Body).Decode(&payload)
    channel <- *payload
}

func main() {
    args, err := flags.Parse(&opts)
    if err != nil {
        panic(err)
    }
    //file, err := os.Create("qhistory.json")
    //if err != nil {
    //    panic(err)
    //}
    //fmt.Printf("Verbosity: %v\n", opts.Verbose)
    var prompt = strings.Join(args, " ")
    //fmt.Printf("Remaining args: %s\n", prompt)

    s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)

    channel := make(chan OpenAIResponse)
    go makeRequest(prompt, channel)
    s.Start()

    response := <- channel
    s.Stop()
    //fmt.Printf("            \r\r")
    //fmt.Printf("%+v\n", response.Choices[0].Message.Content)
    //c := color.New(color.FgCyan)
    c := color.New(color.FgCyan, color.Bold)
    c.Printf("%s\n", response.Choices[0].Message.Content)
}
