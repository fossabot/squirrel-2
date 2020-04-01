package mail

import (
	"errors"
	"fmt"
	"runtime/debug"
	"squirrel/config"
	"squirrel/log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dm"

	eParser "github.com/go-errors/errors"
)

var dmClient *dm.Client
var enabled bool

// Init inits aliyun mail config.
func Init(enableMail bool) {
	var err error

	enabled = enableMail
	if !enableMail {
		return
	}

	if err := config.LoadAliyunMailConfig(); err != nil {
		panic(err)
	}

	mailCfg := config.GetAliyunMailConfig()

	dmClient, err = dm.NewClientWithAccessKey(
		mailCfg.Region,
		mailCfg.AccessKeyID,
		mailCfg.AccessKeySecret)

	if err != nil {
		panic(err)
	}
}

// AlertIfErr Captures paniced error and send mail.
func AlertIfErr() {
	var err error
	if !enabled {
		return
	}

	if r := recover(); r != nil {
		switch t := r.(type) {
		case string:
			err = errors.New(t)
		case error:
			err = t
		default:
			err = errors.New("unknown error")
		}

		err = errors.New(eParser.Wrap(err, 0).ErrorStack())
		log.Error.Println(err)
		SendNotify("Error Detected", err.Error())
	}
}

// SendNotify sends mail to configured receivers.
func SendNotify(subject string, content string) {
	if !enabled {
		return
	}

	if content == "" {
		log.Printf("Mail content cannot be empty\n")
		debug.PrintStack()
		return
	}

	mailCfg := config.GetAliyunMailConfig()

	req := dm.CreateSingleSendMailRequest()
	req.AccountName = mailCfg.AccountName
	req.ReplyToAddress = requests.NewBoolean(false)
	req.AddressType = requests.NewInteger(1)
	if config.GetLabel() != "" {
		req.FromAlias = fmt.Sprintf("[%s]-sq", config.GetLabel())
	} else {
		req.FromAlias = "squirrel"
	}
	req.Subject = subject
	req.TextBody = content
	req.ToAddress = strings.Join(mailCfg.Receiver, ",")

	_, err := dmClient.SingleSendMail(req)

	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
