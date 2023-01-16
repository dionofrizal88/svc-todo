package services

import(
	"sync"
)

type UserLoginStore interface {
    Save(userLogin *UserLogin) error
    Find(username string) (*UserLogin, error)
}

type InMemoryUserLoginStore struct {
    mutex sync.RWMutex
    users map[string]*UserLogin
}

func NewInMemoryUserLoginStore() *InMemoryUserLoginStore {
    return &InMemoryUserLoginStore{
        users: make(map[string]*UserLogin),
    }
}

func (store *InMemoryUserLoginStore) Save(userLogin *UserLogin) error {
    store.mutex.Lock()
    defer store.mutex.Unlock()

    if store.users[userLogin.Username] != nil {
        return nil
    }

    store.users[userLogin.Username] = userLogin.Clone()
    return nil
}

func (store *InMemoryUserLoginStore) Find(username string) (*UserLogin, error) {
    store.mutex.RLock()
    defer store.mutex.RUnlock()

    user := store.users[username]
    if user == nil {
        return nil, nil
    }

    return user.Clone(), nil
}

