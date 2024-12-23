package main

import (
	"fmt"
	"log/slog"
	"reflect"

	"github.com/bedrock-gophers/provider/provider"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
)

func main() {
	chat.Global.Subscribe(chat.StdoutSubscriber{})
	c := server.DefaultConfig()
	c.Players.SaveData = false

	conf, err := c.Config(slog.Default())
	if err != nil {
		panic(err)
	}

	provider.NewProvider(&conf, provider.DefaultSettings())
	srv := conf.New()
	srv.CloseOnProgramEnd()

	srv.Listen()

	for p := range srv.Accept() {
		for _, i := range p.Inventory().Clear() {
			v, _ := i.Value("test")
			fmt.Printf("%v %s\n", v, reflect.TypeOf(v))
		}
		p.Inventory().AddItem(item.NewStack(item.Apple{}, 1).WithValue("test", float64(8.00)))
		p.SetGameMode(world.GameModeCreative)
	}
}
