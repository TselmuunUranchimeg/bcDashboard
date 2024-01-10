package services

import (
	"database/sql"
	"sync"
	"testing"
	"time"

	"bcDashboard/services"
)

func TestBackgroundTask(t *testing.T) {
	db, err := sql.Open("postgres", "user=postgres password=tselmuun100 dbname=bcDashboard port=5432 host=localhost sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(`DELETE FROM "check";`); err != nil {
		t.Fatal(err)
	}
	networks, err := services.InitNetworks(db)
	if err != nil {
		t.Fatal(err)
	}
	hashmaps, err := services.InitHashmaps(db, networks)
	if err != nil {
		t.Fatal(err)
	}
	for {
		var wg sync.WaitGroup
		wg.Add(len(networks))
		ch := make(chan error)
		for i := 0; i < len(networks); i++ {
			go func(i int) {
				defer wg.Done()
				services.BackgroundTask(db, hashmaps[i], 4)
			}(i)
		}
		go func() {
			wg.Wait()
			close(ch)
		}()
		for val := range ch {
			if val != nil {
				t.Fatal(val)
			}
		}
		time.Sleep(time.Second * 7)
	}
}
