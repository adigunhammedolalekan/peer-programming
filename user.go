package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

type User struct {
	ID       int64
	Name     string
	Balance  int64
	Verified bool
}

type UserStore interface {
	Create(name string) (*User, error)
	Update(user *User) error
	GetUser(id int64) (*User, error)
	GetUsers() ([]*User, error)
}

type userStore struct {
	store map[int64]*User
	mu    sync.RWMutex
}

func (u *userStore) Create(name string) (*User, error) {
	u.mu.Lock()

	id := u.nextId()
	user := &User{
		ID:       id,
		Name:     name,
		Balance:  int64(rand.Intn(100)),
		Verified: false,
	}
	u.store[id] = user
	u.mu.Unlock()
	return user, nil
}

func (u *userStore) Update(user *User) error {
	defer u.mu.Unlock()
	u.mu.Lock()

	u.store[user.ID] = user
	return nil
}

func (u *userStore) GetUser(id int64) (*User, error) {
	if user, ok := u.store[id]; ok {
		return user, nil
	}
	return nil, errors.New(fmt.Sprintf("user %d was not found", id))
}

func (u *userStore) GetUsers() ([]*User, error) {
	data := make([]*User, 0, len(u.store))
	for _, v := range u.store {
		data = append(data, v)
	}
	return data, nil
}

func (u *userStore) nextId() int64 {
	return int64(len(u.store) + 1)
}

func NewUserStore() UserStore {
	return &userStore{store: make(map[int64]*User), mu: sync.RWMutex{}}
}
