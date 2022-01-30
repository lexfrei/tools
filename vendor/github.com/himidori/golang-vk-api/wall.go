package vkapi

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

type Wall struct {
	Count    int         `json:"count"`
	Posts    []*WallPost `json:"items"`
	Profiles []*User     `json:"profiles"`
	Groups   []*Group    `json:"groups"`
}

type WallPost struct {
	ID           int                  `json:"id"`
	FromID       int                  `json:"from_id"`
	OwnerID      int                  `json:"owner_id"`
	ToID         int                  `json:"to_id"`
	Date         int64                `json:"date"`
	MarkedAsAd   int                  `json:"marked_as_ads"`
	IsPinned     int                  `json:"is_pinned"`
	PostType     string               `json:"post_type"`
	CopyPostDate int64                `json:"copy_post_date"`
	CopyPostType string               `json:"copy_post_type"`
	CopyOwnerID  int                  `json:"copy_owner_id"`
	CopyPostID   int                  `json:"copy_post_id"`
	CopyHistory  []*WallPost          `json:"copy_history"`
	CreatedBy    int                  `json:"created_by"`
	Text         string               `json:"text"`
	CanDelete    int                  `json:"can_delete"`
	CanPin       int                  `json:"can_pin"`
	Attachments  []*MessageAttachment `json:"attachments"`
	PostSource   *Source              `json:"post_source"`
	Comments     *Comment             `json:"comments"`
	Likes        *Like                `json:"likes"`
	Reposts      *Repost              `json:"reposts"`
	Online       int                  `json:"online"`
	ReplyCount   int                  `json:"reply_count"`
	SignerID     int                  `json:"signer_id"`
}

type Comment struct {
	Count   int `json:"count"`
	CanPost int `json:"can_post"`
}

type Like struct {
	Count      int `json:"count"`
	UserLikes  int `json:"user_likes"`
	CanLike    int `json:"can_like"`
	CanPublish int `json:"can_publish"`
}

type Repost struct {
	Count        int `json:"count"`
	UserReposted int `json:"user_reposted"`
}

type Source struct {
	Type string `json:"type"`
}

func (client *VKClient) WallGet(id interface{}, count int, params url.Values) (*Wall, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("count", strconv.Itoa(count))

	switch id.(type) {
	case int:
		params.Set("owner_id", strconv.Itoa(id.(int)))
	case string:
		params.Set("domain", id.(string))
	default:
		return nil, errors.New("unexpected type in id filed")
	}

	resp, err := client.MakeRequest("wall.get", params)
	if err != nil {
		return nil, err
	}

	var wall *Wall
	json.Unmarshal(resp.Response, &wall)

	return wall, nil
}

func (client *VKClient) WallGetByID(id string, params url.Values) (*Wall, error) {
	if params == nil {
		params = url.Values{}
	}

	params.Set("posts", id)

	resp, err := client.MakeRequest("wall.getById", params)
	if err != nil {
		return nil, err
	}

	wall := &Wall{}
	if params.Get(`extended`) == `1` {
		json.Unmarshal(resp.Response, &wall)
	} else {
		json.Unmarshal(resp.Response, &wall.Posts)
	}
	wall.Count = len(wall.Posts)
	return wall, nil
}

func (client *VKClient) WallPost(ownerID int, message string, params url.Values) (int, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("owner_id", strconv.Itoa(ownerID))
	params.Set("message", message)

	resp, err := client.MakeRequest("wall.post", params)
	if err != nil {
		return 0, err
	}
	m := map[string]int{}
	if err = json.Unmarshal(resp.Response, &m); err != nil {
		return 0, err
	}

	return m["post_id"], nil
}

func (client *VKClient) WallPostComment(ownerID int, postID int, message string, params url.Values) (int, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("owner_id", strconv.Itoa(ownerID))
	params.Set("post_id", strconv.Itoa(postID))
	params.Set("message", message)

	resp, err := client.MakeRequest("wall.createComment", params)
	if err != nil {
		return 0, err
	}
	m := map[string]int{}
	if err = json.Unmarshal(resp.Response, &m); err != nil {
		return 0, err
	}

	return m["comment_id"], nil

}
