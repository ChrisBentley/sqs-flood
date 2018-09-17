package main

import (
	"flag"
	"log"
	"io/ioutil"
	"os"
 	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/satori/go.uuid"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	src := flag.String("src", "", "source file")
	dest := flag.String("dest", "eu-west-1", "destination queue")
	awsRegion := flag.String("region", "", "aws region")
	awsProfile := flag.String("profile", "", "aws profile")
	messageGroupId := flag.String("messageGroupId", "", "message group id for fifo queues only")
	flag.Parse()

	if *src == "" || *dest == "" || *awsProfile == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *awsRegion == "" {
		*awsRegion = "eu-west-1"
	}

	log.Printf("source file : %v", *src)
	log.Printf("destination queue : %v", *dest)
	log.Printf("region : %v", *awsRegion)
	log.Printf("profile : %v", *awsProfile)
	log.Printf("messageGroupId : %v", *messageGroupId)

	client := sqs.New(session.Must(session.NewSessionWithOptions(session.Options{
	    Config:            aws.Config{Region: awsRegion},
	    Profile:           *awsProfile,
	    SharedConfigState: session.SharedConfigEnable,
	    })))

	// Read the src file into memory
	data, err := ioutil.ReadFile(*src)
	check(err)

	// Add an array wrapper to the messages
	dataAsString := string(data[:])
	dataAsString = "[" + dataAsString + "]"

	// Parse the data as json
	jsonParsed, err := gabs.ParseJSON([]byte(dataAsString))

	// Get the messages from the parsed json
	messages, _ := jsonParsed.Children()

	log.Printf("Read in %v messages...", len(messages))

	for i, message := range messages {
		m := message.String()

		fmt.Println("Sending message", i+1, "to the queue...")

		if *messageGroupId == "" {
			smi := sqs.SendMessageInput{
				MessageBody:       			aws.String(m),
				QueueUrl:          			dest,
			}

			_, err := client.SendMessage(&smi)

			if err != nil {
				log.Printf("ERROR sending message to destination %v", err)
				return
			}

		} else { // Add MessageGroupId and MessageDeduplicationId for fifo
			msgDeduplicationId, err1 := uuid.NewV4()
			if err1 != nil {
				panic(err1)
			}

			smi := sqs.SendMessageInput{
				MessageBody:       			aws.String(m),
				QueueUrl:          			dest,
				MessageGroupId:    			messageGroupId,
				MessageDeduplicationId:		aws.String(msgDeduplicationId.String()),
			}

			_, err := client.SendMessage(&smi)

			if err != nil {
				log.Printf("ERROR sending message to destination %v", err)
				return
			}
		}

	}
}
