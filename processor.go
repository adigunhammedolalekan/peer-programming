package main

import (
	"context"
	"log"
	"sync"
)

type Processor interface {
	SubmitUser(user *User)
	SubmitTransaction(transaction *Transaction)
	Start()
	Stop()
}

type processor struct {
	userChan    chan *User
	txChan      chan *Transaction
	workerCount int

	userStore UserStore
	txStore   TransactionStore

	wg         sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func (p *processor) SubmitUser(user *User) {
	p.userChan <- user
}

func (p *processor) SubmitTransaction(transaction *Transaction) {
	p.txChan <- transaction
}

func (p *processor) Start() {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

func (p *processor) worker(id int) {
	defer p.wg.Done()

	log.Printf("Worker %d is running", id)
	for {
		select {
		case <-p.ctx.Done():
			log.Println("DONE is called")
			return
		case user := <-p.userChan:
			p.processUser(user)
		case tx := <-p.txChan:
			p.processTransaction(tx)
		}
	}
}

func (p *processor) processUser(u *User) {
	log.Printf("Processing user: %d", u.ID)
	u.Verified = true
	p.userStore.Update(u)

}

func (p *processor) processTransaction(tx *Transaction) {
	log.Printf("Processing transaction: %d", tx.ID)
	user, err := p.userStore.GetUser(tx.UserId)
	if err != nil {
		log.Println("user not found", err)
		return
	}

	receiver, err := p.userStore.GetUser(tx.ReceiverId)
	if err != nil {
		log.Println("receiver not found", err)
		return
	}

	if !receiver.Verified || !user.Verified {
		log.Printf("Sender or Receiver is not verified: sender_verified=%v id=%d, receiver_verified=%v id=%d",
			user.Verified, user.ID, receiver.Verified, receiver.ID)
		return
	}

	user.Balance -= tx.Amount
	receiver.Balance += tx.Amount

	p.userStore.Update(user)
	p.userStore.Update(receiver)
}

func (p *processor) Stop() {
	log.Println("Stopping task processor...")
	p.cancelFunc()
	p.wg.Wait()
	close(p.txChan)
	close(p.userChan)
}

func NewProcessor(userStore UserStore, store TransactionStore, workerCount int) Processor {
	ctx, cancel := context.WithCancel(context.Background())
	return &processor{
		userChan:    make(chan *User),
		txChan:      make(chan *Transaction),
		cancelFunc:  cancel,
		ctx:         ctx,
		userStore:   userStore,
		txStore:     store,
		workerCount: workerCount,
	}
}
