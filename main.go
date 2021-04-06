package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

var joinedServers []string

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + "ODI2MTgxMTc0NjczNjcwMTQ0.YGIvLQ.SKLXUWuiZ9ZqvddDZ5iTkbIssrg")
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	_ = dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if voiceMessageHandler(s, m, "bruh", "sounds/bruh.mp3") {
		return
	}
	if voiceMessageHandler(s, m, "moan13", "sounds/moan13.mp3") {
		return
	}
	if voiceMessageHandler(s, m, "bann", "sounds/ban_den_weg.mp3") {
		return
	}
	if voiceMessageHandler(s, m, "jamoin", "sounds/ja_moin.mp3") {
		return
	}
	if voiceMessageHandler(s, m, "megalovania", "sounds/megalovania.mp3") {
		return
	}
	if voiceMessageHandler(s, m, "ough", "sounds/ough.mp3") {
		return
	}
	if voiceMessageHandler(s, m, "yooo", "sounds/yooooooooooo.mp3") {
		return
	}
	if voiceMessageHandler(s, m, "haha", "sounds/haha.mp3") {
		return
	}
	if voiceMessageHandler(s, m, "letsgo", "sounds/letsgo.mp3") {
		return
	}

	// If the message is "bing" reply with "Bong!"
	if strings.ToLower(m.Content) == "bing" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Bong!")

		return
	}

	// If the message is "bong" reply with "Bing!"
	if strings.ToLower(m.Content) == "bong" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Bing!")

		return
	}

}

func voiceMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate, message string, filename string) bool {
	m.Content = strings.ToLower(m.Content)
	if m.Content == message {
		playSound(s, m, filename, true)
		return true
	}
	if strings.Contains(m.Content, message) {
		playSound(s, m, filename, false)
		return true
	}
	return false
}

func playSound(s *discordgo.Session, m *discordgo.MessageCreate, filename string, sendErrMsg bool) {

	// s.Guild() funktioniert hier nicht, weil die VoiceStates nur in "state-cached guilds" verfügbar sind,
	// deshalb s.State.Guild()
	st, _ := s.State.Guild(m.GuildID)

	vc := func() (vc *discordgo.Channel) {
		for _, state := range st.VoiceStates {
			if state.UserID == m.Author.ID {
				channel, _ := s.State.Channel(state.ChannelID)
				return channel
			}
		}
		return nil
	}()
	if vc == nil {
		if sendErrMsg {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Bischte dumm oder was? Du muss schon in nem Channel sein kek alda")
		}
		return
	}

	if contains(joinedServers, m.GuildID) {
		return
	}
	joinedServers = append(joinedServers, m.GuildID)

	dvc, err := s.ChannelVoiceJoin(vc.GuildID, vc.ID, false, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	dgvoice.PlayAudioFile(dvc, filename, make(chan bool))
	_ = dvc.Disconnect()

	joinedServers = remove(joinedServers, m.GuildID)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func remove(s []string, e string) []string {
	for i, v := range s {
		if v == e {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
