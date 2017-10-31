package user

import (
	"errors"
	"sync"
	"time"
)

type DataSource struct {
	Users map[int64]Model
	mu	sync.RWMutex
}

func NewDataSource() *DataSource {
	return &DataSource{
		Users:make(map[int64]Model),
	}
}

func (d *DataSource) GetBy(query func(Model) bool) (user Model, found bool) {
	d.mu.RLock()
	for _, user = range d.Users {
		found = query(user)
		if found {
			break
		}
	}
	d.mu.RUnlock()
	return
}

func (d *DataSource) GetByID(id int64) (Model, bool) {
	return d.GetBy(func(u Model) bool {
		return u.ID == id
	})
}

func (d *DataSource) GetByUsername(username string) (Model, bool) {
	return d.GetBy(func(u Model) bool {
		return u.Username == username
	})
}

func (d *DataSource) getLastID() (lastID int64) {
	d.mu.RLock()
	for id := range d.Users {
		if id > lastID {
			lastID = id
		}
	}
	d.mu.RUnlock()
	return lastID
}

func (d *DataSource) InsertOrUpdate(user Model) (Model, error) {
	hashedPassword, err := GeneratePassword(user.password)
	if err != nil {
		return user, err
	}
	user.HashedPassword = hashedPassword
	if id := user.ID; id > 0 {
		_, found := d.GetByID(id)
		if !found {
			return user, errors.New("ID should be zero or a valid one that maps to an existing User")				
		}
		d.mu.Lock()
		d.Users[id] = user
		d.mu.Unlock()
		return user, nil
	}

	id := d.getLastID() + 1
	user.ID = id
	d.mu.Lock()
	user.CreateAt = time.Now()
	d.Users[id] = user
	d.mu.Unlock()
	return user, nil
}