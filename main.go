package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/oudentabetai/twitterlinkfixer-go/discord"
	"github.com/oudentabetai/twitterlinkfixer-go/storage"
)

var (
	dgs *discordgo.Session
)

func main() {
	sessionManager := &discord.DiscordSessionManager{}
	dgs = sessionManager.InitializeSession(storage.Envs.DISCORD_BOT_TOKEN)
	dgs.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsAll | discordgo.PermissionSendMessages
	if err := dgs.Open(); err != nil {
		log.Fatalf("Discordセッションのオープンに失敗: %v", err)
	}
	dgs.AddHandler(discord.OnMessageCreate)
	//dgs.AddHandler(discord.OnInteractionCreate)
	defer dgs.Close()
	log.Println("ボットが起動しました。Ctrl+Cで終了します。")

	//deleteAllGlobalCommands(dgs, os.Getenv("APPLICATION_ID"))
	//SyncCommands(dgs, "", storage.Envs.APPLICATION_ID)
	waitForExitSignal()
}

func waitForExitSignal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func SyncCommands(s *discordgo.Session, guildID string, appID string) {
	commands := []*discordgo.ApplicationCommand{}
	_, err := s.ApplicationCommandBulkOverwrite(appID, guildID, commands)
	if err != nil {
		log.Panicf("コマンドの同期に失敗しました: %v", err)
	}
	log.Println("コマンドを更新しました")
}
