package mqtt

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/eclipse/paho.mqtt.golang"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {
	broker := context.GetInput("broker").(string)
	user := context.GetInput("user").(string)
	password := context.GetInput("password").(string)
	topic := context.GetInput("topic").(string)
	message := context.GetInput("message").(string)

	// do eval
	opts := mqtt.NewClientOptions()
	if strings.Contains(broker, ",") {
		splited := strings.Split(broker, ",")
		for _, val := range splited {
			opts.AddBroker("tcp://" + val)
		}
	} else {
		opts.AddBroker("tcp://" + broker)
	}
	opts.SetUsername(user)
	opts.SetPassword(password)
	opts.SetCleanSession(true)
	client := mqtt.NewClient(opts)
	defer client.Disconnect(250)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return false, token.Error()
	}
	token = client.Publish(topic, byte(2), false, message)
	sent := token.WaitTimeout(5000 * time.Millisecond)
	if !sent {
		log.Printf("Timeout occurred while trying to publish to topic '%s'", topic)
		return false, errors.New("Timeout occurred while trying to publish")
	}

	return true, nil
}
