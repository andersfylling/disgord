package main

import (
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
)

var (
	exampleButtonsIntCreateHandler     chan *disgord.InteractionCreate = make(chan *disgord.InteractionCreate)
	exampleSelectMenuIntCreateHandler  chan *disgord.InteractionCreate = make(chan *disgord.InteractionCreate)
	exampleModalIntCreateHandler       chan *disgord.InteractionCreate = make(chan *disgord.InteractionCreate)
	exampleModalSubmitIntCreateHandler chan *disgord.InteractionCreate = make(chan *disgord.InteractionCreate)
)

func exampleButtonsHandler() {
	for {
		intCreate, active := <-exampleButtonsIntCreateHandler
		if !active {
			log.Debug("exampleButtonsIntCreateHandler no longer active")
		}
		if intCreate.ApplicationID != appID {
			continue
		}
		if intCreate.Data.Name != "example_buttons" {
			continue
		}
		if intCreate.Data.Type != disgord.ApplicationCommandChatInput {
			continue
		}

		buttons := []*disgord.MessageComponent{
			{
				Type:     disgord.MessageComponentActionRow,
				CustomID: "0",
				Components: []*disgord.MessageComponent{
					{
						Type:     disgord.MessageComponentButton,
						Style:    disgord.Primary,
						CustomID: "1",
						Label:    "my smiley button",
						Emoji: &disgord.Emoji{
							Name: "🙂",
						},
					},
					{
						Type:     disgord.MessageComponentButton,
						Style:    disgord.Primary,
						CustomID: "2",
						Label:    "my frowny button",
						Emoji: &disgord.Emoji{
							Name: "🙁",
						},
					},
					{
						Type:     disgord.MessageComponentButton,
						Style:    disgord.Primary,
						CustomID: "3",
						Label:    "my melon button",
						Emoji: &disgord.Emoji{
							Name: "🍉",
						},
					},
				},
			},
		}
		resp := &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{
					{
						Title:       "My Buttons",
						Description: "Don't push them, please.",
					},
				},
				Components: buttons,
			},
		}
		err := client.SendInteractionResponse(context.Background(), intCreate, resp)
		if err != nil {
			log.Debug(err)
		}
	}
}

func exampleSelectMenuHandler() {
	for {
		intCreate, active := <-exampleSelectMenuIntCreateHandler
		if !active {
			log.Debug("exampleSelectMenuIntCreateHandler no longer active")
		}
		if intCreate.ApplicationID != appID {
			continue
		}
		if intCreate.Data.Name != "example_select_menu" {
			continue
		}
		if intCreate.Data.Type != disgord.ApplicationCommandChatInput {
			continue
		}

		selectMenu := []*disgord.MessageComponent{
			{
				Type:     disgord.MessageComponentActionRow,
				CustomID: "0",
				Components: []*disgord.MessageComponent{
					{
						Type:        disgord.MessageComponentSelectMenu,
						CustomID:    "1",
						MinValues:   1,
						MaxValues:   1,
						Placeholder: "Select the best fruit",
						Options: []*disgord.SelectMenuOption{
							{
								Label: "Watermelon",
								Value: "watermelon",
								Emoji: &disgord.Emoji{
									Name: "🍉",
								},
							},
							{
								Label: "Mango",
								Value: "mango",
								Emoji: &disgord.Emoji{
									Name: "🥭",
								},
							},
							{
								Label: "Potato",
								Value: "potato",
								Emoji: &disgord.Emoji{
									Name: "🥔",
								},
							},
						},
					},
				},
			},
		}
		resp := &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{
					{
						Title:       "The Best Fruit",
						Description: "Don't steal them, please.",
					},
				},
				Components: selectMenu,
			},
		}
		err := client.SendInteractionResponse(context.Background(), intCreate, resp)
		if err != nil {
			log.Debug(err)
		}
	}
}

func exampleModalHandler() {
	for {
		intCreate, active := <-exampleModalIntCreateHandler
		if !active {
			log.Debug("exampleModalIntCreateHandler no longer active")
		}
		if intCreate.ApplicationID != appID {
			continue
		}
		if intCreate.Data.Name != "example_modal" {
			continue
		}
		if intCreate.Data.Type != disgord.ApplicationCommandChatInput {
			continue
		}

		modal := []*disgord.MessageComponent{
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type:        disgord.MessageComponentTextInput,
						Style:       disgord.TextInputStyleShort,
						CustomID:    "1",
						Label:       "Name",
						MinValues:   1,
						MaxValues:   1000,
						Placeholder: "e.g. Jon",
						Required:    true,
					},
				},
			},
		}
		resp := &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackModal,
			Data: &disgord.CreateInteractionResponseData{
				Title:      "My Modal",
				CustomID:   "my_modal",
				Components: modal,
			},
		}
		err := client.SendInteractionResponse(context.Background(), intCreate, resp)
		if err != nil {
			log.Debug(err)
		}
	}
}

func exampleModalSubmitHandler() {
	for {
		intCreate, active := <-exampleModalSubmitIntCreateHandler
		if !active {
			log.Debug("exampleModalSubmitIntCreateHandler no longer active")
		}
		if intCreate.ApplicationID != appID {
			continue
		}
		if intCreate.Type != disgord.InteractionModalSubmit {
			continue
		}

		resp := &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{
					{
						Title:       "Modal Submit Message",
						Description: fmt.Sprintf("Submitted Answer: %s", intCreate.Data.Components[0].Components[0].),
					},
				},
			},
		}
		err := client.SendInteractionResponse(context.Background(), intCreate, resp)
		if err != nil {
			log.Debug(err)
		}
	}
}
