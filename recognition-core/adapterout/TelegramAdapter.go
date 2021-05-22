package recognitionadapterout

import (
	"go-intelligent-monitoring-system/domain"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

//TelegramAdapter ...
type TelegramAdapter struct {
	bot *tb.Bot
}

func (ta *TelegramAdapter) Recipient() string {
	return os.Getenv("TELEGRAM_CHANNEL")
}

func (ta *TelegramAdapter) NotifyInitializedSystem() error {
	// The message is sent to the telegram user/channel defined in ta.Recipient() func
	_, err := ta.bot.Send(ta, "System initialized")

	if err != nil {
		log.WithError(err).Error("Error sending Initialized System Notification")
	} else {
		log.Info("Telegram Bot initialized")
	}

	return err
}

//NotifyUnauthorizedFace ...
func (ta *TelegramAdapter) NotifyUnauthorizedFace(notification domain.Notification) error {

	// The message is sent to the telegram user defined in ta.Recipient() func
	msg := notification.Message + "\n \n " + notification.Image.URL + "\n \n "
	_, err := ta.bot.Send(ta, msg)

	return err
}

//NewTelegramAdapter initializes a TelegramAdapter object.
func NewTelegramAdapter() *TelegramAdapter {

	b, err := tb.NewBot(tb.Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".
		//URL: "http://195.129.111.17:8012",

		Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return nil
	}

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hello "+os.Getenv("TELEGRAM_USER"))
	})

	go b.Start()

	telBot := &TelegramAdapter{
		bot: b,
	}

	defer telBot.NotifyInitializedSystem()

	return telBot
}
