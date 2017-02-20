package main

import (
    zmq "github.com/pebbe/zmq4"
    rpn "github.com/irlndts/go-rpn"
    "bytes"
    "fmt"
    "strings"
    "strconv"
    "time"
    "math/rand"
)

func main() {
    context, _ := zmq.NewContext()
    
    //  Socket to receive messages on
    receiver, _ := context.NewSocket(zmq.PULL)
    defer receiver.Close()
    receiver.Connect("tcp://localhost:5557")

    //  Socket to send messages to task sink
    sender, _ := context.NewSocket(zmq.PUSH)
    defer sender.Close()
    sender.Connect("tcp://localhost:5558")

    fmt.Println("Server listening...")
    //  Process tasks forever
    for {
        msgbytes, _ := receiver.Recv(0)

        msgstr := string(msgbytes)
        req := strings.Split(msgstr, "\n")

        jobId, jobNumStr, jobInputs := req[0], req[1], req[2:]
        fmt.Println("ID", jobId, "| processing", jobNumStr, "RPNs")

        jobNum, _ := strconv.Atoi(jobNumStr)

        var response bytes.Buffer
        response.WriteString(jobId)

        for i := 0; i < jobNum; i++ {
            rpnResult, _ := rpn.Calc(jobInputs[i])
            rpnResultStr := strconv.FormatFloat(rpnResult, 'f', 10, 64)
            response.WriteString("\n")
            response.WriteString(rpnResultStr)

            fmt.Println("ID", jobId, "|", i, "|", rpnResult)
        }

        time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

        //  Send results to sink
        sender.Send(response.String(), 0)

    }
}