package vk2tg

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	vkapi "github.com/himidori/golang-vk-api"
	"github.com/pkg/errors"
	tb "gopkg.in/telebot.v3"
)

// Moscow
//nolint:gomnd // just a time
var zone = time.FixedZone("UTC+3", 3*60*60)

type VTClinent struct {
	config     *config
	tgClient   *tb.Bot
	vkClient   *vkapi.VKClient
	LastUpdate time.Time
	StartTime  time.Time
	WG         *sync.WaitGroup
	ticker     *time.Ticker
	chVKPosts  chan *vkapi.WallPost
	logger     *log.Logger
	stateFile  string
}

type config struct {
	TGUser       int64         `yaml:"tgUser"`
	Paused       bool          `yaml:"paused"`
	Silent       bool          `yaml:"silent"`
	TGToken      string        `yaml:"tgToken"`
	VKToken      string        `yaml:"vkToken"`
	Period       time.Duration `yaml:"period"`
	LastPostID   int           `yaml:"lastPostId"`
	LastPostDate int64         `yaml:"lastPostDate"`
}

func NewVTClient(tgToken, vkToken string, tgRecepient int64, period time.Duration) *VTClinent {
	vtcli := new(VTClinent)
	vtcli.config = new(config)
	vtcli.config.TGToken = tgToken
	vtcli.config.VKToken = vkToken
	vtcli.config.TGUser = tgRecepient
	vtcli.chVKPosts = make(chan *vkapi.WallPost)
	vtcli.WG = &sync.WaitGroup{}
	vtcli.config.Silent = false
	vtcli.config.Paused = false
	vtcli.StartTime = time.Now()
	vtcli.config.Period = period
	vtcli.ticker = time.NewTicker(period)
	vtcli.logger = log.New(io.Discard, "vk2tg: ", log.Ldate|log.Ltime|log.Lshortfile)

	return vtcli
}

func (vtCli *VTClinent) WithLogger(logger *log.Logger) *VTClinent {
	vtCli.logger = logger

	return vtCli
}

func (vtCli *VTClinent) WithConfig(path string) *VTClinent {
	vtCli.stateFile = path

	return vtCli
}

func (vtCli *VTClinent) Start() error {
	vtCli.logger.Println("Starting...")

	var err error

	vtCli.vkClient, err = vkapi.NewVKClientWithToken(vtCli.config.VKToken, nil, true)
	if err != nil {
		return errors.Wrap(err, "Can't longin to VK")
	}

	vtCli.tgClient, err = tb.NewBot(
		tb.Settings{
			Token:  vtCli.config.TGToken,
			Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		},
	)
	if err != nil {
		return errors.Wrap(err, "Can't longin to TG")
	}

	var commands []tb.Command

	vtCli.tgClient.Handle("/status", vtCli.status)

	commands = append(commands,
		tb.Command{Text: "status", Description: "Show status"})

	vtCli.tgClient.Handle("/pause", vtCli.pause)

	commands = append(commands,
		tb.Command{Text: "pause", Description: "(Un)pause bot"})

	vtCli.tgClient.Handle("/mute", vtCli.mute)

	commands = append(commands,
		tb.Command{Text: "mute", Description: "(Un)mute bot"})

	err = vtCli.tgClient.SetCommands(commands)

	if err != nil {
		return errors.Wrap(err, "can't set commands")
	}

	go vtCli.tgClient.Start()

	//nolint:gomnd // workers amount
	vtCli.WG.Add(2)

	go vtCli.VKWatcher()
	go vtCli.TGSender()

	return nil
}

func (vtCli *VTClinent) Pause() {
	vtCli.logger.Println("Watcher paused")
	vtCli.ticker.Stop()
	vtCli.config.Paused = true
}

func (vtCli *VTClinent) Resume() {
	vtCli.logger.Println("Watcher unpaused")
	vtCli.ticker.Reset(vtCli.config.Period)
	vtCli.config.Paused = false
}

func (vtCli *VTClinent) Mute() {
	vtCli.config.Silent = true
}

func (vtCli *VTClinent) Unmute() {
	vtCli.config.Silent = false
}

func (vtCli *VTClinent) Wait() {
	vtCli.WG.Wait()
}

func (vtCli *VTClinent) VKWatcher() {
	vtCli.WG.Add(1)
	defer vtCli.WG.Done()

	for range vtCli.ticker.C {
		vtCli.LastUpdate = time.Now()

		vkWall, err := vtCli.vkClient.WallGet("cosplay_second", 10, nil)
		if err != nil {
			vtCli.logger.Printf("failed to fetch posts: %s", err)

			continue
		}

		if vkWall.Posts[0].ID == vtCli.config.LastPostID {
			continue
		}

		for index := len(vkWall.Posts) - 1; index >= 0; index-- {
			vtCli.logger.Printf("Post %d: Processing", vkWall.Posts[index].ID)

			if vtCli.config.LastPostID > vkWall.Posts[index].ID {
				vtCli.logger.Printf("Post %d: Not a new post, skipped", vkWall.Posts[index].ID)

				continue
			} else {
				vtCli.logger.Printf("Post %d: Selected as latest", vkWall.Posts[index].ID)
				vtCli.config.LastPostDate = vkWall.Posts[index].Date
				vtCli.config.LastPostID = vkWall.Posts[index].ID
			}

			if !strings.Contains(vkWall.Posts[index].Text, "#поиск") {
				vtCli.logger.Printf("Post %d: Post does not contain required substring, skipping", vkWall.Posts[index].ID)

				continue
			}

			vtCli.logger.Printf("Post %d: Sending to TG", vkWall.Posts[index].ID)
			vtCli.chVKPosts <- vkWall.Posts[index]
		}
	}
}

func (vtCli *VTClinent) TGSender() {
	vtCli.WG.Add(1)
	defer vtCli.WG.Done()
	defer vtCli.logger.Println("Sender: done")

	for post := range vtCli.chVKPosts {
		var (
			album tb.Album
			url   *string
		)

		for attachmentsIndex := range post.Attachments {
			if post.Attachments[attachmentsIndex].Type == "photo" {
				var maxSize int

				for sizeIndex := range post.Attachments[attachmentsIndex].Photo.Sizes {
					//nolint:lll // whis can't be shorter
					if maxSize < post.Attachments[attachmentsIndex].Photo.Sizes[sizeIndex].Width*post.Attachments[attachmentsIndex].Photo.Sizes[sizeIndex].Height {
						maxSize = post.Attachments[attachmentsIndex].Photo.Sizes[sizeIndex].Width *
							post.Attachments[attachmentsIndex].Photo.Sizes[sizeIndex].Height
						url = &post.Attachments[attachmentsIndex].Photo.Sizes[sizeIndex].Url
					}
				}

				album = append(album, &tb.Photo{
					File: tb.FromURL(*url),
				})
			}
		}

		if len(album) > 0 {
			if _, err := vtCli.tgClient.SendAlbum(&tb.User{ID: vtCli.config.TGUser}, album); err != nil {
				vtCli.logger.Printf("Can't send album: %s\n", err)
			}
		}

		_, err := vtCli.tgClient.Send(
			&tb.User{ID: vtCli.config.TGUser},
			post.Text,
			vtCli.generateOptionsForPost(post),
		)
		if err != nil {
			vtCli.logger.Printf("Can't send message: %s\n", err)

			return
		}

		chat, err := vtCli.tgClient.ChatByID(vtCli.config.TGUser)
		if err != nil {
			vtCli.logger.Printf("Error on fetching user info: %s", err)

			return
		}

		vtCli.logger.Printf("Post %d: Sent to %s", post.ID, chat.FirstName)
	}
}

func (vtCli *VTClinent) sendMessage(u *tb.User, options ...interface{}) error {
	if _, err := vtCli.tgClient.Send(u, options); err != nil {
		return errors.Wrap(err, "error on sending message")
	}

	return nil
}

func (vtCli *VTClinent) status(tbContext tb.Context) error {
	msg := fmt.Sprintf("I'm fine\nLast post date:\t%s\nReceived in:\t%s\nUptime:\t%s\nPaused:\t%t\nSound:\t%t",
		time.Unix(vtCli.config.LastPostDate, 0).In(zone).Format(time.RFC822),
		vtCli.LastUpdate.In(zone).Format(time.RFC822),
		time.Since(vtCli.StartTime).Round(time.Second),
		vtCli.config.Paused,
		!vtCli.config.Silent,
	)

	if _, err := vtCli.tgClient.Send(tbContext.Sender(), msg); err != nil {
		return errors.Wrap(err, "error on sending message")
	}

	return nil
}

func (vtCli *VTClinent) pause(tbContext tb.Context) error {
	if !vtCli.config.Paused {
		vtCli.Pause()

		if err := vtCli.sendMessage(tbContext.Sender(), "Paused! Send /pause to continue"); err != nil {
			vtCli.logger.Println(err)
		}
	} else {
		vtCli.Resume()

		if err := vtCli.sendMessage(tbContext.Sender(), "Unpaused! Send /pause to stop"); err != nil {
			vtCli.logger.Println(err)
		}
	}

	return nil
}

func (vtCli *VTClinent) mute(tbContext tb.Context) error {
	if !vtCli.config.Silent {
		vtCli.Mute()

		if err := vtCli.sendMessage(tbContext.Sender(), "Muted! Send /mute to go loud"); err != nil {
			vtCli.logger.Println(err)
		}
	} else {
		vtCli.Unmute()

		if err := vtCli.sendMessage(tbContext.Sender(), "Unmuted! Send /mute to go silent"); err != nil {
			vtCli.logger.Println(err)
		}
	}

	return nil
}

func (vtCli *VTClinent) generateOptionsForPost(post *vkapi.WallPost) *tb.SendOptions {
	return &tb.SendOptions{
		ReplyTo: &tb.Message{},
		ReplyMarkup: &tb.ReplyMarkup{
			InlineKeyboard: [][]tb.InlineButton{
				{
					tb.InlineButton{
						Text: "🌎 К посту",
						URL:  "https://vk.com/wall-57692133_" + strconv.Itoa(post.ID),
					},
					tb.InlineButton{
						Text: "✍️ Написать",
						URL:  "vk.com/write" + strconv.Itoa(post.SignerID),
					},
				},
			},
		},
		DisableNotification: vtCli.config.Silent,
	}
}
