package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dangerousgameofpool/discord-wordle-bot/words"
	"github.com/enescakir/emoji"
	"github.com/joho/godotenv"
)

func main() {
	botToken := envLoad()
	discord, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatal("Error creating Discord session")
		return
	}

	discord.AddHandler(messageCreate)
	discord.Open()

	fmt.Println("Bot is online!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

	discord.Close()
}

func envLoad() string {
	godotenv.Load()
	return os.Getenv("BOT_TOKEN")
}

var (
	gameStarted bool
	w           wordle
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// TODO slash commands for this??

	// TODO support custom word length
	if strings.HasPrefix(m.Content, "!play") {
		if gameStarted {
			s.ChannelMessageSend(m.ChannelID, "There's already an active game running! Use `!end` to kill it and start anew.")
			return
		}

		w = NewWordle(5)
		gameStarted = true
		s.ChannelMessageSend(m.ChannelID, "Wordle has started! Send a word preceded by `!guess` to play.")
	}

	if gameStarted {
		// TODO game board should be displayed via embed
		if strings.HasPrefix(m.Content, "!guess") {
			args := strings.Split(m.Content, " ")
			w.processGuess(args[1])
			s.ChannelMessageSendEmbed(m.ChannelID, w.embedBoard())
		}

		// TODO remove this later
		if strings.HasPrefix(m.Content, "!answer") {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The answer is: ||%s||", w.Answer()))
		}

		// TODO no response on this command atm, get it working
		if strings.HasPrefix(m.Content, "!history") {
			s.ChannelMessageSend(m.ChannelID, w.History("\n"))
		}

		if strings.HasPrefix(m.Content, "!end") {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Thanks for playing! The answer was ||%s||", w.answer))
			w.endGame()
			gameStarted = false
		}
	}
}

// A wordle represents a wordle puzzle.
type wordle struct {
	board   string           // Shows the current state of the wordle game
	turns   int              // Tracks the number of turns that have passed.
	isMatch bool             // Checks if the user's guess is a match to answer
	play    bool             // Controls if the game continues or ends
	dict    words.Dictionary // Keeps a wordle's available wordlist and gives random words
	answer  string           // The answer to a wordle puzzle.
	history []string         // Keeps a history of the player's guesses
}

// NewWordle creates and returns a new wordle struct.
func NewWordle(l int) wordle {
	w := wordle{
		turns:   0,
		isMatch: false,
		play:    true,
		dict:    words.NewDictionary(l),
	}
	// Not sure if there's a more elegant way to do this
	// syntactically. Putting it inside the struct wouldn't work.
	w.answer = w.dict.RandomWord()
	return w
}

// embedBoard returns a MessageEmbed containing a wordle's board string.
func (w wordle) embedBoard() *discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x6AAA64,
		Description: w.board,
		Title:       "Wordle",
	}
	return &embed
}

// Answer returns the string held by a wordle's answer field.
func (w wordle) Answer() string {
	return w.answer
}

func (w *wordle) processGuess(g string) {
	clue := ""
	for i := 0; i < len(g); i++ {
		if g[i] == w.answer[i] {
			clue += emoji.GreenSquare.String()
		} else if strings.Contains(g, string(w.answer[i])) {
			clue += emoji.YellowSquare.String()
		} else {
			clue += emoji.WhiteLargeSquare.String()
		}
	}
	w.history = append(w.history, g)
	clue += "\n"
	w.updateBoard(clue)
}

// appendHistory appends a user's guess to a wordle's history slice.
func (w *wordle) appendHistory(g string) {
	w.history = append(w.history, g)
}

// History returns a string representation of a wordle's history slice.
// Receives an argument for a delimiter to use.
func (w wordle) History(delimiter string) string {
	return strings.Join(w.history, delimiter)
}

func (w *wordle) updateBoard(s string) {
	s += "\n"
	w.board += s
}

// endGame sets a Wordle struct's "play" field to false, ending the current game.
// Allows the user to end a game prematurely, in addition to terminating it when
// their guess is an exact match for the answer.
func (w *wordle) endGame() {
	w.play = false
}
