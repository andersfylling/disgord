package main

import "github.com/andersfylling/disgord"

var (
	testcmdIntCreateHandler chan *disgord.InteractionCreate = make(chan *disgord.InteractionCreate)
)

func testcmdHandler() {
	for {
		intCreate, active := <-testcmdIntCreateHandler
		if !active {
			log.Debug("testcmdIntCreateHandler no longer active")
		}
		if intCreate.ApplicationID != appID {
			continue
		}
		buttons := []*disgord.MessageComponent{
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type:  disgord.MessageComponentButton,
						Style: disgord.Primary,
						Label: "my smiley button",
						Emoji: &disgord.Emoji{
							Name: "🙂",
						},
					},
					{
						Type:  disgord.MessageComponentButton,
						Style: disgord.Primary,
						Label: "my frowny button",
						Emoji: &disgord.Emoji{
							Name: "🙁",
						},
					},
					{
						Type:  disgord.MessageComponentButton,
						Style: disgord.Primary,
						Label: "my melon button",
						Emoji: &disgord.Emoji{
							Name: "🍉",
						},
					},
				},
			},
		}
		buttonMsg := &disgord.CreateMessage{
			Embeds: []*disgord.Embed{
				{
					Title:       "My Buttons",
					Description: "Don't push them, please.",
				},
			},
			Components: buttons,
		}
		selectMenu := []*disgord.MessageComponent{
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type: disgord.MessageComponentSelectMenu,
						Options: []*disgord.SelectMenuOption{
							{
								Label: "Watermelon",
								Value: "1",
								Emoji: &disgord.Emoji{
									Name: "🍉",
								},
							},
							{
								Label: "Mango",
								Value: "2",
								Emoji: &disgord.Emoji{
									Name: "🥭",
								},
							},
							{
								Label: "Potato",
								Value: "3",
								Emoji: &disgord.Emoji{
									Name: "🥔",
								},
							},
						},
					},
				},
			},
		}
		menuMsg := &disgord.CreateMessage{
			Embeds: []*disgord.Embed{
				{
					Title:       "The Best Fruit",
					Description: "Don't steal them, please.",
				},
			},
			Components: selectMenu,
		}
		modal := []*disgord.MessageComponent{
			{
				Title:    "My Modal",
				CustomID: "my_modal",
				Components: []*disgord.MessageComponent{
					{
						Type: disgord.MessageComponentActionRow,
						Components: []*disgord.MessageComponent{
							{
								Type:        disgord.MessageComponentTextInput,
								Style:       1,
								Label:       "Name",
								MinValues:   1,
								MaxValues:   1000,
								Placeholder: "e.g. Jon",
								Required:    true,
							},
						},
					},
				},
			},
		}
		modalMsg := &disgord.CreateMessage{
			Components: modal,
		}
		client.Channel(intCreate.ChannelID).CreateMessage(buttonMsg)
		client.Channel(intCreate.ChannelID).CreateMessage(menuMsg)
		client.Channel(intCreate.ChannelID).CreateMessage(modalMsg)
	}
}
