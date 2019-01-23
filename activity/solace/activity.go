package solace

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"pack.ag/amqp"
)

var log = logger.GetLogger("activity-solacepublisher")

var senderLock sync.Mutex

var storedHostIP string

var storedVpnName string

var storedTopicName string

var storedClient *amqp.Client

var storedSession *amqp.Session

var storedSender *amqp.Sender

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
func (a *MyActivity) Eval(actCtx activity.Context) (done bool, err error) {

	// do eval
	hostIP := actCtx.GetInput("hostIP").(string)
	vpnName := actCtx.GetInput("vpnName").(string)
	topicName := actCtx.GetInput("topicName").(string)
	username := actCtx.GetInput("username").(string)
	password := actCtx.GetInput("password").(string)
	data := actCtx.GetInput("data").(string)

	// 1. if many data send in same time, it will create many connections that may over solace connection number limit, you can run the activity_test and switch the TEST_ROUTINE_NUM to 600 or 700
	// 2. so we use mutex lock temporarily now
	// TODO 3. check which amqp object(client, session, sender) are concurrent safe, and use github.com/jolestar/go-commons-pool to pooling it's children
	senderLock.Lock()
	defer senderLock.Unlock()
	sender, err := getSender(hostIP, vpnName, topicName, username, password)

	if err != nil {
		actCtx.SetOutput("publishSuccess", "false")
		return false, err
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer sender.Close(ctx)
	defer cancel()
	if strings.HasPrefix(data, "{") {
		data = "[" + data + "]"
	}
	err = sender.Send(ctx, amqp.NewMessage([]byte(data)))

	if err != nil {
		sender.Close(ctx)
		storedSender = nil
		clearStore()

		log.Errorf("AMQP Send error: [%s]", err)
		actCtx.SetOutput("publishSuccess", "false")
		return false, err
	}

	log.Debugf("Solace publisher send successfully by topicName [%s] and data [%s]", topicName, data)
	actCtx.SetOutput("publishSuccess", "true")
	return true, nil
}

func getSender(hostIP, vpnName, topicName, username, password string) (*amqp.Sender, error) {
	if hostIP != storedHostIP || vpnName != storedVpnName || topicName != storedTopicName || storedClient == nil || storedSession == nil || storedSender == nil {
		err := createSender(hostIP, vpnName, topicName, username, password)
		if err != nil {
			return nil, err
		}
	}
	return storedSender, nil
}

func createSender(hostIP, vpnName, topicName, username, password string) error {
	var client *amqp.Client
	var err error
	if username != "" {
		client, err = amqp.Dial("amqp://"+hostIP+"/"+vpnName, amqp.ConnSASLPlain(username, password))
	} else {
		client, err = amqp.Dial("amqp://" + hostIP + "/" + vpnName)
	}
	if err != nil {
		log.Errorf("AMQP Dial fail by hostIP [%s] and vpnName [%s] and error is [%s]", hostIP, vpnName, err)
		return err
	}
	session, err := client.NewSession()
	if err != nil {
		log.Error("AMQP NewSession fail and error is [%s]", err)
		return err
	}
	sender, err := session.NewSender(amqp.LinkTargetAddress("topic://" + topicName))
	if err != nil {
		log.Errorf("AMQP NewSender fail by topicName [%s] and error is [%s]", topicName, err)
		return err
	}
	storedHostIP = hostIP
	storedVpnName = vpnName
	storedTopicName = topicName
	storedClient = client
	storedSession = session
	storedSender = sender
	return nil
}

func clearStore() {
	storedHostIP = ""
	storedVpnName = ""
	storedTopicName = ""
	if storedSession == nil {
		storedSession.Close(context.Background())
	}
	if storedClient == nil {
		storedClient.Close()
	}
}
