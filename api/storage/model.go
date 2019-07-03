package storage

// QueueFromDB should store values from DB
type QueueFromDB struct {
	Domain   string
	Weight   int
	Priority int
}
