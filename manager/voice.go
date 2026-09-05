package manager

import (
	"context"
	"time"

	"github.com/TheTipo01/YADMB/constants"
	"github.com/TheTipo01/YADMB/embed"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

// JoinVC joins the voice channel if not already joined, returns true if joined successfully
func JoinVC(e *events.ApplicationCommandInteractionCreate, channelID snowflake.ID, server *Server, isDeferred chan struct{}) bool {
	if !server.VC.IsConnected() {
		// Join the voice channel
		err := server.VC.Join(channelID, e.Client())
		if err != nil {
			embed.SendAndDeleteEmbedInteraction(discord.NewEmbed().WithTitle(BotName).AddField(constants.ErrorTitle, constants.CantJoinVC, false).
				WithColor(0x7289DA), e, time.Second*5, isDeferred)
			return false
		}
	}
	return true
}

// QuitVC disconnects the bot from the voice channel after a short idle if nothing is playing
func (server *Server) QuitVC() {
	if !server.Queue.IsEmpty() {
		return
	}

	server.VC.Disconnect()

	// Also send Opcode 4 leave on the main gateway. Close() alone can no-op when
	// the voice WS is stuck reconnecting — GRUPIM then keeps the bot in the roster.
	if server.Clients != nil && server.Clients.Discord != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Clients.Discord.UpdateVoiceState(ctx, server.GuildID, nil, false, false)
	}
}

// FindUserVoiceState finds user current voice channel
func FindUserVoiceState(s *bot.Client, guildID, userID snowflake.ID) *discord.VoiceState {
	v, found := s.Caches.VoiceState(guildID, userID)

	if !found {
		return nil
	}

	return &v
}
