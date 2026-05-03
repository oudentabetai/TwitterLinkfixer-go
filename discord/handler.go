package discord

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// 各サービスのプレフィックスと変換先のマップ
var conversionRules = []struct {
	pattern    string
	prefix     string
	replaceTo  string
	exceptions []string
}{
	{
		pattern:    `(https://twitter\.com|https://x\.com)`,
		prefix:     "https://",
		replaceTo:  "https://fxtwitter.com",
		exceptions: []string{"https://twitter.com/", "https://x.com/"},
	},
	{
		pattern:    `https://www\.instagram\.com`,
		prefix:     "https://www.instagram.com",
		replaceTo:  "https://www.uuinstagram.com",
		exceptions: []string{"https://www.instagram.com/"},
	},
	{
		pattern:    `https://pixiv\.net`,
		prefix:     "https://pixiv.net",
		replaceTo:  "https://phixiv.net",
		exceptions: []string{"https://pixiv.net/"},
	},
	{
		pattern:    `https://soundcloud\.com`,
		prefix:     "https://soundcloud.com",
		replaceTo:  "https://fxcloud.ofton.dev",
		exceptions: []string{"https://soundcloud.com/"},
	},
	{
		pattern:    `https://open\.spotify\.com`, // Goに合わせて元の文字列を使用する想定
		prefix:     "https://open.spotify.com",
		replaceTo:  "https://open.fxspotify.com",
		exceptions: []string{}, // 元のコードの exceptions に合わせる場合ここで指定
	},
}

// ConvertMessage はメッセージ内のURLを条件に応じて変換する
func ConvertMessage(msg string) (string, bool) {
	for _, rule := range conversionRules {
		// 例外条件に完全に一致する場合は処理をスキップ
		isException := false
		for _, ex := range rule.exceptions {
			if msg == ex {
				isException = true
				break
			}
		}
		if isException {
			continue
		}

		// URLが含まれているかチェック
		if strings.Contains(msg, rule.prefix) {
			re := regexp.MustCompile(rule.pattern)
			converted := re.ReplaceAllString(msg, rule.replaceTo)
			return converted, true
		}
	}

	return "", false // 変換が行われなかった場合
}

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// メッセージがボット自身のものであれば無視
	if m.Author.ID == s.State.User.ID {
		return
	}
	content := m.Content
	converted, changed := ConvertMessage(content)
	if changed {
		// 変換されたURLを含むメッセージを送信
		_, err := s.ChannelMessageSend(m.ChannelID, converted)
		if err != nil {
			// エラー処理（例: ログに記録）
			return
		}
	}
}
