package main

import "errors"

type APiService struct {
	userStore        UserStore
	transactionStore TransactionStore
	processor        Processor
}

func NewApiService(userS UserStore, txStore TransactionStore, processor Processor) *APiService {
	return &APiService{
		userStore:        userS,
		transactionStore: txStore,
		processor:        processor,
	}
}

func (api *APiService) CreateUser(name string) (*User, error) {
	user, err := api.userStore.Create(name)
	if err != nil {
		return nil, err
	}

	api.processor.SubmitUser(user)
	return user, nil
}

func (api *APiService) GetUsers() ([]*User, error) {
	return api.userStore.GetUsers()
}

func (api *APiService) CreateTransaction(userId, receiverId, amount int64) (*Transaction, error) {
	user, err := api.userStore.GetUser(userId)
	if err != nil {
		return nil, err
	}
	_, err = api.userStore.GetUser(receiverId)
	if err != nil {
		return nil, err
	}

	if user.Balance < amount {
		return nil, errors.New("insufficient funds")
	}

	tx, err := api.transactionStore.Create(userId, receiverId, amount)
	if err != nil {
		return nil, err
	}

	api.processor.SubmitTransaction(tx)
	return tx, nil
}
