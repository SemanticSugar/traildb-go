package main

import (
	"fmt"

	"github.com/lumost/traildb-go"
)

const SESSION_LIMIT = 30

type WikiEvent struct {
	TimeStamp int64  `tdb:"timestamp"`
	User      string `tdb:"user"`
	IP        string `tdb:"ip"`
	Title     string `tdb:"title"`
}

func main() {
	db, err := tdb.Open("wikipedia-history-small.tdb")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(db)
	fmt.Println(db.Version())
	countSessions(db)

}

func countSessions(db *tdb.TrailDB) {

	for i := 0; i < db.NumTrails; i++ {
		trail, err := tdb.NewTrail(db, i)
		if err != nil {
			panic(err.Error())
		}
		//Trail Scoped session counter for number of sessions per user
		numSessions := 0
		numEvents := 0
		prevTime := int64(0)
		for {
			evt := trail.NextEvent()
			if evt == nil {
				trail.Close()
				break
			}
			we := WikiEvent{}
			// cast to event to struct
			//  then assert that it is reflect.Value
			//  then use Interface() to get an interface for the value
			/// then type assert to a pointer of the original WikiEvent type
			parsedWikiEvent := evt.ToStruct(we).(*WikiEvent)
			if parsedWikiEvent.TimeStamp-prevTime > SESSION_LIMIT {
				numSessions++
			}
			numEvents++
			prevTime = parsedWikiEvent.TimeStamp
		}
		fmt.Printf("Number of Sessions: %d Number of Events: %d\n", numSessions, numEvents)
	}
}
