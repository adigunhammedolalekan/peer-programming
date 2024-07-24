package main

import "sync"

type Transaction struct {
	ID         int64
	Amount     int64
	UserId     int64
	ReceiverId int64
	Status     string
}

type TransactionStore interface {
	Create(userId, receiverId, amount int64) (*Transaction, error)
}

type transactionStore struct {
	store map[int64]*Transaction
	mu    sync.RWMutex
}

func (t *transactionStore) Create(userId, receiverId, amount int64) (*Transaction, error) {
	defer t.mu.Unlock()
	t.mu.Lock()

	id := t.nextId()
	tx := &Transaction{
		ID:         id,
		Amount:     amount,
		ReceiverId: receiverId,
		UserId:     userId,
		Status:     "CREATED",
	}
	t.store[id] = tx
	return tx, nil
}

func (t *transactionStore) nextId() int64 {
	return int64(len(t.store) + 1)
}

func NewTransactionStore() TransactionStore {
	return &transactionStore{
		store: make(map[int64]*Transaction),
		mu:    sync.RWMutex{},
	}
}
