package core

import "github.com/bwmarrin/discordgo"

type ChannelCategory struct {
	channel  *discordgo.Channel
	channels []*discordgo.Channel
}

func (cc *ChannelCategory) Position() int {
	return cc.Position()
}

func (cc *ChannelCategory) Channels() []Channel {
	c := make([]Channel, len(cc.channels))
	for i, chann := range cc.channels {
		c[i] = &TextChannel{
			channel: chann,
		}
	}
	return c
}

func (cc *ChannelCategory) ID() string {
	return cc.channel.ID
}

func (cc *ChannelCategory) Channel() *discordgo.Channel {
	return cc.channel
}

func (cc *TextChannel) ID() string {
	return cc.channel.ID
}

func (cc *TextChannel) Position() int {
	return cc.channel.Position
}

type Category interface {
	Position() int
	ID() string
	Channels() []Channel
}

type TextChannel struct {
	channel *discordgo.Channel
}

type Channel interface {
	Position() int
	ID() string
}

func SortChannels(channels []*discordgo.Channel) []*ChannelCategory {
	var categories = []*ChannelCategory{}

	for _, channel := range channels {
		if channel.ParentID == "" {
			categories = append(categories, &ChannelCategory{
				channel: channel,
			})
		}
	}

	for _, category := range categories {
		for _, channel := range channels {
			if channel.ParentID == category.ID() {
				category.channels = append(category.channels, channel)
			}
		}
	}

	return categories
}
