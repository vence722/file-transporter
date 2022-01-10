package server

import "sync"

type UserStore struct {
	users map[string]struct{}
	lock  *sync.Mutex
}

func (u *UserStore) AddUser(username string) bool {
	u.lock.Lock()
	defer u.lock.Unlock()
	if _, ok := u.users[username]; ok {
		return false
	}
	u.users[username] = struct{}{}
	return true
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
	return &UserStore{map[string]struct{}{}, &sync.Mutex{}}
}
