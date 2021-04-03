package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	gomniauthtest "github.com/stretchr/gomniauth/test"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	tsetChatUser := &chatUser{User: testUser}
	url, err := authAvatar.GetAvatarURL(tsetChatUser)
	if err != ErrNoAvatarURL {
		t.Error("If Value is not exists, AuthAvatar.GetAvatarURL should return ErrNoAvatarURL.")
	}
	// setting value.
	testUrl := "http://url-to-avatar/"
	testUser = &gomniauthtest.TestUser{}
	tsetChatUser.User = testUser
	testUser.On("AvatarURL").Return(testUrl, nil)
	url, err = authAvatar.GetAvatarURL(tsetChatUser)
	if err != nil {
		t.Error("If Value is exist, AuthAvatar.GetAvatarURL should not return any Error.")
	} else {
		if url != testUrl {
			t.Error("AuthAvatar.GetAvatarURL should return correct URL.")
		}
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := gravatarAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("GravatarAvatar should not return Error.")
	}
	if url != "//www.gravatar.com/avatar/abc" {
		t.Errorf("GravatarAvatar.GetAvatarURL returns %s ,this is incorrect value.", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	// create Avatar File for Test.
	filename := filepath.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer func() { os.Remove(filename) }()

	var fileSystemAvatar FileSystemAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURL should not return Error.")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURL returns invalid value [%s]", url)
	}
}
