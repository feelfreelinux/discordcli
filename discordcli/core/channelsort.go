package core

import "github.com/bwmarrin/discordgo"

/*
ChannelCategory keeps reference to all of its child channels
*/
type ChannelCategory struct {
	Channel  *discordgo.Channel
	Channels []*discordgo.Channel
}

/*
SortChannels sorts given channels and maps them to []*ChannelCategory
*/
func SortChannels(channels []*discordgo.Channel) []*ChannelCategory {
	var categories = []*ChannelCategory{}

	for _, channel := range channels {
		if channel.ParentID == "" {
			categories = append(categories, &ChannelCategory{
				Channel: channel,
			})
		}
	}

	for _, category := range categories {
		for _, channel := range channels {
			if channel.ParentID == category.Channel.ID {
				category.Channels = append(category.Channels, channel)
			}
		}
	}

	return categories
}
