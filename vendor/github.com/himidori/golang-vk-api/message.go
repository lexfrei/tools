package vkapi

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

const (
	ActivityTypeTyping   = "typing"
	ActivityTypeAudioMsg = "audiomessage"
)

type Dialog struct {
	Count    int     `json:"count"`
	Messages []*Item `json:"items"`
}

type Message struct {
	Count    int              `json:"count"`
	Messages []*DialogMessage `json:"items"`
}

type Item struct {
	Message *DialogMessage `json:"message"`
	InRead  int            `json:"in_read"`
	OutRead int            `json:"out_read"`
}

type DialogMessage struct {
	MID               int                  `json:"id"`
	Date              int64                `json:"date"`
	Out               int                  `json:"out"`
	UID               int                  `json:"user_id"`
	ReadState         int                  `json:"read_state"`
	Title             string               `json:"title"`
	Body              string               `json:"body"`
	RandomID          int                  `json:"random_id"`
	ChatID            int64                `json:"chat_id"`
	ChatActive        string               `json:"chat_active"`
	PushSettings      *Push                `json:"push_settings"`
	UsersCount        int                  `json:"users_count"`
	AdminID           int                  `json:"admin_id"`
	Photo50           string               `json:"photo_50"`
	Photo100          string               `json:"photo_100"`
	Photo200          string               `json:"photo_200"`
	ForwardedMessages []*ForwardedMessage  `json:"fwd_messages"`
	Attachments       []*MessageAttachment `json:"attachments"`
}

type Push struct {
	Sound         int   `json:"sound"`
	DisabledUntil int64 `json:"disabled_until"`
}

type ForwardedMessage struct {
	UID               int                  `json:"user_id"`
	Date              int64                `json:"date"`
	Body              string               `json:"body"`
	Attachments       []*MessageAttachment `json:"attachments"`
	ForwardedMessages []*ForwardedMessage  `json:"fwd_messages"`
}

type MessageAttachment struct {
	Type     string             `json:"type"`
	Audio    *AudioAttachment   `json:"audio"`
	Video    *VideoAttachment   `json:"video"`
	Photo    *PhotoAttachment   `json:"photo"`
	Document *DocAttachment     `json:"doc"`
	Link     *LinkAttachment    `json:"link"`
	Wall     *WallPost          `json:"wall"`
	Sticker  *StickerAttachment `json:"sticker"`
}

type StickerAttachment struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	Photo64   string `json:"photo_64"`
	Photo128  string `json:"photo_128"`
	Photo256  string `json:"photo_256"`
	Photo352  string `json:"photo_352"`
	Photo512  string `json:"photo_512"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

type HistoryAttachment struct {
	Attachments []HistoryAttachmentItem `json:"items"`
	NextFrom    string                  `json:"next_from"`
}

type HistoryAttachmentItem struct {
	MID        int                `json:"message_id"`
	Attachment *MessageAttachment `json:"attachment"`
}

type AudioAttachment struct {
	ID        int    `json:"id"`
	OwnerID   int    `json:"owner_id"`
	Artist    string `json:"artist"`
	Title     string `json:"title"`
	Duration  int    `json:"duration"`
	URL       string `json:"url"`
	Performer string `json:"performer"`
}

type LinkAttachment struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Target      string `json:"target"`
}

type Keyboard struct {
	OneTime 	bool		`json:"one_time"`
	Buttons		[][]Button 	`json:"buttons"`
	Inline		bool 		`json:"inline"`
}

type Button struct {
	Action		map[string]string	`json:"action"`
	Color		string				`json:"color,omitempty"`
}

func (client *VKClient) DialogsGet(count int, params url.Values) (*Dialog, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Add("count", strconv.Itoa(count))

	resp, err := client.MakeRequest("messages.getDialogs", params)
	if err != nil {
		return nil, err
	}

	var dialog *Dialog
	json.Unmarshal(resp.Response, &dialog)

	return dialog, nil
}

func (client *VKClient) GetHistoryAttachments(peerID int, mediaType string, count int, params url.Values) (*HistoryAttachment, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Add("count", strconv.Itoa(count))
	params.Add("media_type", mediaType)
	params.Add("peer_id", strconv.Itoa(peerID))

	resp, err := client.MakeRequest("messages.getHistoryAttachments", params)
	if err != nil {
		return nil, err
	}

	var att *HistoryAttachment
	json.Unmarshal(resp.Response, &att)
	return att, nil
}

func (client *VKClient) MessagesGet(count int, chatID int, isDialog bool, params url.Values) (int, []*DialogMessage, error) {
	if params == nil {
		params = url.Values{}
	}
	if isDialog {
		chatID += 2000000000
	}

	params.Add("user_id", strconv.Itoa(chatID))
	params.Add("count", strconv.Itoa(count))

	resp, err := client.MakeRequest("messages.getHistory", params)
	if err != nil {
		return 0, nil, err
	}

	var message *Message
	json.Unmarshal(resp.Response, &message)

	return message.Count, message.Messages, nil
}

func (client *VKClient) MessagesGetByID(message_ids []int, params url.Values) (int, []*DialogMessage, error) {
	if params == nil {
		params = url.Values{}
	}
	s := ArrayToStr(message_ids)
	params.Add("message_ids", s)

	resp, err := client.MakeRequest("messages.getById", params)
	if err != nil {
		return 0, nil, err
	}

	var message *Message
	json.Unmarshal(resp.Response, &message)

	return message.Count, message.Messages, nil
}

func (client *VKClient) MessagesSend(peerOrDomain interface{}, message string, params url.Values) (APIResponse, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Add("message", message)

	switch peerOrDomain.(type) {
	case int:
		params.Add("peer_id", strconv.Itoa(peerOrDomain.(int)))
	case string:
		params.Add("domain", peerOrDomain.(string))
	}

	resp, err := client.MakeRequest("messages.send", params)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (client *VKClient) MessagesDelete(ids []int, spam int, deleteForAll int) (int, error) {
	params := url.Values{}
	s := ArrayToStr(ids)
	params.Add("message_ids", s)
	params.Add("spam", strconv.Itoa(spam))
	params.Add("delete_for_all", strconv.Itoa(deleteForAll))

	resp, err := client.MakeRequest("messages.delete", params)
	if err != nil {
		return 0, err
	}

	delCount := 0
	var idMap map[string]int
	reader := strings.NewReader(string(resp.Response))
	err = json.NewDecoder(reader).Decode(&idMap)
	if err != nil {
		return 0, err
	}

	for _, v := range idMap {
		if v == 1 {
			delCount++
		}
	}

	return delCount, nil
}

func (client *VKClient) MessagesSetActivity(user int, params url.Values) error {
	if params == nil {
		params = url.Values{}
	}

	params.Add("user_id", strconv.Itoa(user))

	_, err := client.MakeRequest("messages.setActivity", params)
	if err != nil {
		return err
	}

	return nil
}
