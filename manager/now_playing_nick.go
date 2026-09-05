package manager

import (
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/lit"
	"github.com/disgoorg/disgo/discord"
)

const nowPlayingNickMax = 32

// setNowPlayingNick updates the bot's guild nick so GRUPIM (and Discord) can
// show a car-CD-style “now playing” marquee under the bot in the voice member list.
// Pass empty title to clear the nick when playback stops.
//
// Uses PATCH /guilds/{id}/members/{botId} (not @me) so it works on GRUPIM
// deployments that already support Modify Member by user id.
func (server *Server) setNowPlayingNick(title string) {
	if server.Clients == nil || server.Clients.Discord == nil {
		return
	}

	var nick *string
	if trimmed := strings.TrimSpace(title); trimmed != "" {
		formatted := formatNowPlayingNick(trimmed)
		nick = &formatted
	} else {
		empty := ""
		nick = &empty
	}

	botID := server.Clients.Discord.ApplicationID
	go func() {
		_, err := server.Clients.Discord.Rest.UpdateMember(
			server.GuildID,
			botID,
			discord.MemberUpdate{Nick: nick},
		)
		if err != nil {
			lit.Debug("failed to update now-playing nick for guild %s: %v", server.GuildID, err)
		}
	}()
}

func formatNowPlayingNick(title string) string {
	prefix := "♪ "
	budget := nowPlayingNickMax - utf8.RuneCountInString(prefix)
	runes := []rune(title)
	if len(runes) > budget {
		if budget > 1 {
			runes = append(runes[:budget-1], '…')
		} else {
			runes = runes[:budget]
		}
	}
	return prefix + string(runes)
}
