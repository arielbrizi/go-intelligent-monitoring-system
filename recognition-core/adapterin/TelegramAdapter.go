package recognitionadapterin

import (
	"go-intelligent-monitoring-system/domain"
	"os"
	"time"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

//TelegramAdapter ...
type TelegramAdapter struct {
	bot         *tb.Bot
	redisClient *redis.Client
}

func (ta *TelegramAdapter) Recipient() string {
	return os.Getenv("TELEGRAM_CHANNEL")
}

func (ta *TelegramAdapter) NotifyInitializedSystem() error {
	// The message is sent to the telegram user/channel defined in ta.Recipient() func
	_, err := ta.bot.Send(ta, "Input System initialized")

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

//GetLastRecognizedImage return last recognized face
func (ta *TelegramAdapter) GetLastRecognizedImage() string {
	val, err := ta.redisClient.Get("lastRecognized").Result()
	if err != nil {
		log.WithError(err).Error("error getting redis value: lastRecognized")
	}
	return val
}

//GetLastNotRecognizedImage return last not recognized face
func (ta *TelegramAdapter) GetLastNotRecognizedImage() string {
	val, err := ta.redisClient.Get("lastNotRecognized").Result()
	if err != nil {
		log.WithError(err).Error("error getting redis value: lastNotRecognized")
	}
	return val
}

//ActivateRecognition activate face recognition
func (ta *TelegramAdapter) ActivateRecognition() error {
	errRedis := ta.redisClient.Set("statusRecognition", "ON", 0).Err()
	if errRedis != nil {
		log.WithError(errRedis).Error("error saving in redis")
	}
	return errRedis
}

//DeactivateRecognition deactivate face recognition
func (ta *TelegramAdapter) DeactivateRecognition() error {
	errRedis := ta.redisClient.Set("statusRecognition", "OFF", 0).Err()
	if errRedis != nil {
		log.WithError(errRedis).Error("error saving in redis")
	}
	return errRedis
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

	/*
		To be able to get the list of commands from Telegram App. Chat BotFather:
		- /mybots
		- Select "ims_config_bot"
		- Edit Bot
		- Edit Commands

		hello - HealthCheck
		last_not_recognized - Get Last Not Recognized
		last_recognized - Get Last Recognized
		activate_recognition - Activate Recognition
		deactivate_recognition - Deactivate Recognition
	*/

	b.Handle("/menu", func(m *tb.Message) {
		b.Send(m.Sender, "ups...under construction")
		/* Read
		-  https://github.com/tucnak/telebot/tree/v2.3.5
		-  https://core.telegram.org/bots#inline-keyboards-and-on-the-fly-updating
		*/
	})

	telBot := &TelegramAdapter{
		bot: b,
	}

	telBot.redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis-server:6379",
		Password: os.Getenv("REDIS_PASS"),
		DB:       0, // use default DB
	})

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hello "+os.Getenv("TELEGRAM_USER"))
	})

	b.Handle("/last_not_recognized", func(m *tb.Message) {
		b.Send(m.Sender, "The last person not recognized was: "+telBot.GetLastNotRecognizedImage())
	})

	b.Handle("/last_recognized", func(m *tb.Message) {
		b.Send(m.Sender, "The last person recognized was: "+telBot.GetLastRecognizedImage())
	})

	b.Handle("/activate_recognition", func(m *tb.Message) {
		err := telBot.ActivateRecognition()
		if err != nil {
			b.Send(m.Sender, "Recognition could not be Activated ")
		} else {
			b.Send(m.Sender, "Recognition Activated ")

		}
	})

	b.Handle("/deactivate_recognition", func(m *tb.Message) {
		err := telBot.DeactivateRecognition()
		if err != nil {
			b.Send(m.Sender, "Recognition could not be Deactivated ")
		} else {
			b.Send(m.Sender, "Recognition Deactivated ")

		}
	})

	go b.Start()

	defer telBot.NotifyInitializedSystem()

	defer telBot.ActivateRecognition()

	return telBot
}
