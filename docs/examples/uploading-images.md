Uploading images as attachments or even to supply them in Embeds is quite easy,
images not used in the same post will be added as attachments
and images accessed using the `attachment://` scheme will be used on their respective locations.  
See the following example:

```go
f1, err := os.Open("myfavouriteimage.jpg")
if err != nil {
	panic(err)
}
defer f1.Close()
f2, err := os.Open("another.jpg")
if err != nil {
	panic(err)
}
defer f2.Close()

client.CreateMessage(channelID, &disgord.CreateMessageParams{
	Content: "This is my favourite image, and another in an embed!",
	Files: []disgord.CreateMessageFileParams{
		{f1, "myfavouriteimage.jpg", false},
		{f2, "another.jpg", false},
	},
	Embed: &disgord.Embed{
		Description: "Look here!",
		Image: &disgord.EmbedImage{
			URL: "attachment://another.jpg",
		},
	},
})
```