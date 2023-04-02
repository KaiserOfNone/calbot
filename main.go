/* See license at the bottom of the file */
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string = ""
)

type CommandFunc func(s *discordgo.Session, i *discordgo.InteractionCreate)

type Command struct {
	Name        string
	Description string
	Cmd         CommandFunc
}

type Bot struct {
	session  *discordgo.Session
	commands []*discordgo.ApplicationCommand
	cmdmap   map[string]CommandFunc
}

func New(token string) (*Bot, error) {
	bot := Bot{
		cmdmap: map[string]CommandFunc{},
	}
	commands := []Command{
		{
			Name:        "help",
			Description: "gives you help",
			Cmd:         bot.help,
		},
	}
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	bot.session = dg
	err = bot.session.Open()
	if err != nil {
		return nil, err
	}
	err = bot.RegisterCommands(commands)
	if err != nil {
		return nil, err
	}
	return &bot, nil
}

func (b *Bot) RegisterCommands(commands []Command) error {
	for _, bcmd := range commands {
		appcmd := &discordgo.ApplicationCommand{
			Name:        bcmd.Name,
			Description: bcmd.Description,
		}
		cmd, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", appcmd)
		if err != nil {
			return err
		}
		b.commands = append(b.commands, cmd)
		b.cmdmap[bcmd.Name] = bcmd.Cmd
	}
	b.session.AddHandler(b.dispatchCommands)
	return nil
}

func (b *Bot) DeleteCommands() error {
	for _, cmd := range b.commands {
		err := b.session.ApplicationCommandDelete(b.session.State.User.ID, "", cmd.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) Stop() error {
	err := b.DeleteCommands()
	if err != nil {
		return err
	}
	b.session.Close()
	return nil
}

func (b *Bot) dispatchCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	recvCmd := i.ApplicationCommandData().Name
	if handler, ok := b.cmdmap[recvCmd]; ok {
		handler(s, i)
	}
}

func (b *Bot) help(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				"Hello digus",
			),
		},
	})
}

func main() {
	flag.StringVar(&Token, "t", "", "The secret token for the bot")
	flag.Parse()
	if Token == "" {
		flag.Usage()
		os.Exit(1)
	}

	bot, err := New(Token)
	if err != nil {
		log.Println("error creating Discord session: ", err)
		os.Exit(1)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
	err = bot.Stop()
	if err != nil {
		log.Println("error stopping bot: ", err)
		os.Exit(1)
	}
}

/*
   Copyright (C) 2023  Andres Villagra

   This program is free software; you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation; either version 2 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License along
   with this program; if not, write to the Free Software Foundation, Inc.,
   51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
*/
