package main

import (
	"errors"
)

// ErrNoAvatar is Error in case of not returning Avatar URL by Avatar Instance.
var ErrNoAvatarURL = errors.New("chat: Can't receive Avatar's URL.")

// Avatar is Type of user's profile image.
type Avatar interface {
	// GetAvatarURL returns specified Avatar's URL.
	// In case of Error, return Error. Especially return ErrNoAvatarURL if user can't receive URL.
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "/avatars/" + useridStr + ".jpeg", nil
		}
	}
	return "", ErrNoAvatarURL
}
