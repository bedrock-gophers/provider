package main

import (
	"github.com/bedrock-gophers/provider/provider"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.InfoLevel

	chat.Global.Subscribe(chat.StdoutSubscriber{})
	c := server.DefaultConfig()
	c.Players.SaveData = false

	conf, err := c.Config(log)
	if err != nil {
		log.Fatalln(err)
	}

	conf.PlayerProvider = provider.NewProvider(provider.DefaultSettings())
	srv := conf.New()
	srv.CloseOnProgramEnd()

	srv.Listen()

	for srv.Accept(func(p *player.Player) {
		p.SetGameMode(world.GameModeCreative)
	}) {

	}
}
