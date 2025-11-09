package vk2tg

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	vkapi "github.com/SevereCloud/vksdk/v3/api"
	vkObject "github.com/SevereCloud/vksdk/v3/object"
	"github.com/cockroachdb/errors"
	tb "gopkg.in/telebot.v4"
)

// Moscow
//
//nolint:mnd // just a time
var zone = time.FixedZone("UTC+3", 3*60*60)

type VTClinent struct {
	config     *config
	tgClient   *tb.Bot
	vkClient   *vkapi.VK
	LastUpdate time.Time
	StartTime  time.Time
	WG         *sync.WaitGroup
	ticker     *time.Ticker
	chVKPosts  chan *vkObject.WallWallpost
	logger     *log.Logger
	storage    storage
}

type config struct {
	LastPostDate int           `yaml:"lastPostDate"`
	LastPostID   int           `yaml:"lastPostId"`
	Paused       bool          `yaml:"paused"`
	Period       time.Duration `yaml:"period"`
	Silent       bool          `yaml:"silent"`
	TGToken      string        `yaml:"tgToken"`
	TGUser       int64         `yaml:"tgUser"`
	VKToken      string        `yaml:"vkToken"`

	// Storage
	StorageEnabled bool `yaml:"storageEnabled"`

	// Hidden items
	serviceName string
}

func NewVTClient(tgToken, vkToken string, tgRecepient int64, period time.Duration) *VTClinent {
	vtcli := new(VTClinent)
	vtcli.config = new(config)
	vtcli.config.TGToken = tgToken
	vtcli.config.VKToken = vkToken
	vtcli.config.TGUser = tgRecepient
	vtcli.chVKPosts = make(chan *vkObject.WallWallpost, 10)
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

func (vtCli *VTClinent) Start() error {
	vtCli.logger.Println("Starting...")

	var err error

	vtCli.vkClient = vkapi.NewVK(vtCli.config.VKToken)

	vtCli.tgClient, err = tb.NewBot(
		tb.Settings{
			Token:  vtCli.config.TGToken,
			Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		},
	)
	if err != nil {
		return errors.Wrap(err, "Can't longin to TG")
	}

	vtCli.tgClient.Handle("/status", vtCli.status)
	vtCli.tgClient.Handle("/pause", vtCli.pause)
	vtCli.tgClient.Handle("/mute", vtCli.mute)

	err = vtCli.tgClient.SetCommands(
		[]tb.Command{
			{Text: "mute", Description: "(Un)mute bot"},
			{Text: "pause", Description: "(Un)pause bot"},
			{Text: "status", Description: "Show current status"},
		},
	)
	if err != nil {
		return errors.Wrap(err, "can't set commands")
	}

	go vtCli.tgClient.Start()

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

		vkWall, err := vtCli.vkClient.WallGet(
			vkapi.Params{
				"owner_id": -57692133,
				"count":    10,
			},
		)
		if err != nil {
			vtCli.logger.Printf("failed to fetch posts: %s", err)

			continue
		}

		if vkWall.Items[0].ID == vtCli.config.LastPostID {
			continue
		}

		for index := vkWall.Count - 1; index >= 0; index-- {
			vtCli.logger.Printf("Post %d: Processing", vkWall.Items[index].ID)

			if vtCli.config.LastPostID >= vkWall.Items[index].ID {
				vtCli.logger.Printf("Post %d: Not a new post, skipped", vkWall.Items[index].ID)

				continue
			}

			vtCli.logger.Printf("Post %d: Selected as latest", vkWall.Items[index].ID)
			vtCli.config.LastPostDate = vkWall.Items[index].Date
			vtCli.config.LastPostID = vkWall.Items[index].ID

			if vtCli.config.StorageEnabled {
				vtCli.storage.SetLastPost(vkWall.Items[index].ID)
			}

			if !strings.Contains(vkWall.Items[index].Text, "#поиск") {
				vtCli.logger.Printf("Post %d: Post does not contain required substring, skipping", vkWall.Items[index].ID)

				continue
			}

			vtCli.logger.Printf("Post %d: Sending to TG", vkWall.Items[index].ID)

			vtCli.chVKPosts <- &vkWall.Items[index]
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
				var maxSize float64

				for sizeIndex := range post.Attachments[attachmentsIndex].Photo.Sizes {
					//nolint:lll // whis can't be shorter
					if maxSize < post.Attachments[attachmentsIndex].Photo.Sizes[sizeIndex].Width*post.Attachments[attachmentsIndex].Photo.Sizes[sizeIndex].Height {
						maxSize = post.Attachments[attachmentsIndex].Photo.Sizes[sizeIndex].Width *
							post.Attachments[attachmentsIndex].Photo.Sizes[sizeIndex].Height
						url = &post.Attachments[attachmentsIndex].Photo.Sizes[sizeIndex].URL
					}
				}

				album = append(album, &tb.Photo{
					File: tb.FromURL(*url),
				})
			}
		}

		if len(album) > 0 {
			_, err := vtCli.tgClient.SendAlbum(&tb.User{ID: vtCli.config.TGUser}, album)
			if err != nil {
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

func (vtCli *VTClinent) sendMessage(u *tb.User, options ...any) error {
	_, err := vtCli.tgClient.Send(u, options)
	if err != nil {
		return errors.Wrap(err, "error on sending message")
	}

	return nil
}

func (vtCli *VTClinent) status(tbContext tb.Context) error {
	msg := fmt.Sprintf("I'm fine\nLast post date:\t%s\nReceived in:\t%s\nUptime:\t%s\nPaused:\t%t\nSound:\t%t",
		time.Unix(int64(vtCli.config.LastPostDate), 0).In(zone).Format(time.RFC822),
		vtCli.LastUpdate.In(zone).Format(time.RFC822),
		time.Since(vtCli.StartTime).Round(time.Second),
		vtCli.config.Paused,
		!vtCli.config.Silent,
	)

	_, err := vtCli.tgClient.Send(tbContext.Sender(), msg)
	if err != nil {
		return errors.Wrap(err, "error on sending message")
	}

	return nil
}

func (vtCli *VTClinent) pause(tbContext tb.Context) error {
	if !vtCli.config.Paused {
		vtCli.Pause()

		err := vtCli.sendMessage(tbContext.Sender(), "Paused! Send /pause to continue")
		if err != nil {
			vtCli.logger.Println(err)
		}
	} else {
		vtCli.Resume()

		err := vtCli.sendMessage(tbContext.Sender(), "Unpaused! Send /pause to stop")
		if err != nil {
			vtCli.logger.Println(err)
		}
	}

	return nil
}

func (vtCli *VTClinent) mute(tbContext tb.Context) error {
	if !vtCli.config.Silent {
		vtCli.Mute()

		err := vtCli.sendMessage(tbContext.Sender(), "Muted! Send /mute to go loud")
		if err != nil {
			vtCli.logger.Println(err)
		}
	} else {
		vtCli.Unmute()

		err := vtCli.sendMessage(tbContext.Sender(), "Unmuted! Send /mute to go silent")
		if err != nil {
			vtCli.logger.Println(err)
		}
	}

	return nil
}

func (vtCli *VTClinent) generateOptionsForPost(post *vkObject.WallWallpost) *tb.SendOptions {
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
