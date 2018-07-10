package bot

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/thecodeah/Gopher/src/commands"
	"github.com/thecodeah/Gopher/src/storage"
)

// Configuration contains settings loaded in from environment
// variable (.env) files.
type Configuration struct {
	Token    string `required:"true"`
	Prefix   string `default:"!"`
	Cooldown int    `default:"3"`

	MongoHost string `split_words:"true" default:"127.0.0.1" required:"true"`
	MongoPort string `split_words:"true" default:"27017" required:"true"`
	MongoName string `split_words:"true" default:"gopher" required:"true"`
	MongoUser string `split_words:"true" required:"false"`
	MongoPass string `split_words:"true" required:"false"`
}

// Bot contains information that's necessary for the bot.
type Bot struct {
	config         Configuration
	session        *discordgo.Session
	commandHandler *commands.CommandHandler
	database       *storage.Database
}

// New initializes the bot as well as all commands.
func New(config Configuration) (bot *Bot, err error) {
	bot = &Bot{
		config: config,
	}

	bot.database, err = storage.Connect(storage.Config{
		Host: config.MongoHost,
		Port: config.MongoPort,
		Name: config.MongoName,
		User: config.MongoUser,
		Pass: config.MongoPass,
	})
	if err != nil {
		return
	}

	bot.session, err = discordgo.New("Bot " + config.Token)
	if err != nil {
		return
	}

	err = bot.session.Open()
	if err != nil {
		return
	}

	bot.commandHandler = commands.New(commands.Config{
		Prefix:   config.Prefix,
		Cooldown: time.Duration(config.Cooldown) * time.Second,
	})
	bot.registerCommands()

	bot.session.AddHandler(bot.messageCreate)

	return
}

// Close closes the bot
func (bot Bot) Close() {
	bot.session.Close()
}

func (bot Bot) messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	bot.commandHandler.Process(session, message)
}

func (bot Bot) registerCommands() {
	bot.commandHandler.Register("ping", commands.PingCommand)
}
