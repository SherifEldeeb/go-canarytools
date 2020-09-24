package canarytools

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	log "github.com/sirupsen/logrus"
)

// SQSForwarder sends alerts to sqs
type SQSForwarder struct {
	queueName string
	queueURL  string
	sess      *session.Session
	svc       *sqs.SQS
	l         *log.Logger
}

// NewSQSForwarder creates a new SQSForwarder
func NewSQSForwarder(queueName string, l *log.Logger) (sqsforwarder *SQSForwarder, err error) {
	sqsforwarder = &SQSForwarder{}
	sqsforwarder.queueName = queueName
	sqsforwarder.l = l

	// check
	if queueName == "" {
		err = fmt.Errorf("queue name can't be empty")
		return
	}

	// create aws session
	sqsforwarder.sess, err = session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return
	}

	//TODO: Check if we have enough perm to do SQS

	// create SQS svc
	sqsforwarder.svc = sqs.New(sqsforwarder.sess)

	// check if queue exists
	// get Queue URL
	getQueueUrlOutput, err := sqsforwarder.svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		// QueueName: aws.String(*queue),
		QueueName: aws.String(queueName),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok { // is this sdk error?
			switch awsErr.Code() {
			case sqs.ErrCodeQueueDoesNotExist: // it does not exist
				l.WithField("SQSQueueName", queueName).Info("SQS queue does not exists")
				l.WithField("SQSQueueName", queueName).Info("attempting to create SQS queue")
				// create queue
				createQueueOutput, createQueueErr := sqsforwarder.svc.CreateQueue(&sqs.CreateQueueInput{
					QueueName: aws.String(queueName),
				})
				if createQueueErr != nil {
					return nil, fmt.Errorf("error creating SQS queue: %v", createQueueErr)
				}
				l.WithField("createQueueOutput", createQueueOutput.String()).Debug("SQS queue created successfully")
				sqsforwarder.queueURL = *createQueueOutput.QueueUrl
				// we cool now
				err = nil
			default:
				return
			}
		} else {
			return
		}
	} else {
		sqsforwarder.queueURL = *getQueueUrlOutput.QueueUrl
	}
	return
}

func (sqsf SQSForwarder) Forward(outChan <-chan []byte, incidentAckerChan chan<- []byte) {
	for i := range outChan {
		sqsf.l.WithFields(log.Fields{
			"source": "SQSForwarder",
			"stage":  "forward",
		}).Debug("SQS out incident")
		sendMessageOutput, err := sqsf.svc.SendMessage(&sqs.SendMessageInput{
			QueueUrl:    aws.String(sqsf.queueURL),
			MessageBody: aws.String(string(i)),
		})
		if err != nil {
			sqsf.l.WithFields(log.Fields{
				"source": "SQSForwarder",
				"stage":  "forward",
				"err":    err,
			}).Error("SQS error")
		}
		sqsf.l.WithFields(log.Fields{
			"source":            "SQSForwarder",
			"stage":             "forward",
			"sendMessageOutput": sendMessageOutput.String(),
		}).Trace("SQS")

		incidentAckerChan <- i
	}
}
