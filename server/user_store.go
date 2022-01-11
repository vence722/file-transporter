package server

import (
	"net"
	"sync"
)

type UserStore struct {
	users map[string]net.Conn
	lock  *sync.Mutex
}

func (u *UserStore) AddUser(username string, conn net.Conn) bool {
	u.lock.Lock()
	defer u.lock.Unlock()
	if _, ok := u.users[username]; ok {
		return false
	}
	u.users[username] = conn
	return true
}

func (u *UserStore) GetUser(username string) (net.Conn, bool) {
	u.lock.Lock()
	defer u.lock.Unlock()
	if conn, ok := u.users[username]; ok {
		return conn, ok
	}
	return nil, false
}

func (u *UserStore) RemoveUser(username string) bool {
	u.lock.Lock()
	defer u.lock.Unlock()
	if _, ok := u.users[username]; !ok {
		return false
	}
	delete(u.users, username)
	return true
}

func (u *UserStore) ListUsers() []string {
	u.lock.Lock()
	defer u.lock.Unlock()
	var userList []string
	for user := range u.users {
		userList = append(userList, user)
	}
	return userList
}

func NewUserStore() *UserStore {
	return &UserStore{map[string]net.Conn{}, &sync.Mutex{}}
}
