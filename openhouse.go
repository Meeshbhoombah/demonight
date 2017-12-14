/*
main.go

Backend for Daaank's Open House application
*/

package main

import (
        "log"
        "net/http"

        "github.com/googollee/go-socket.io"
	"gopkg.in/mgo.v2"
)

const (
	SERVER = "localhost"
	DBNAME = "daaank-open-house"
	COLLEC = "points"
)

type Point struct {
        PointString string
}

func main() {
        /* CONNECT TO DB */
	session, err := mgo.Dial(SERVER)
	if err != nil {
		panic(err)
		log.Println("Could not create connection: " + err.Error())
	}

	defer session.Close()

        session.SetMode(mgo.Monotonic, true)
        c := session.DB(DBNAME).C(COLLEC)
        log.Println(c)

        /* SOCKET */
        server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

        server.On("connection", func(so socketio.Socket) {
                log.Println("CONNECTED")

                so.Join("demonight")
                so.On("location", func(loc string) {
                    log.Println("EVENT 'location': ", loc)

                    var p Point
                    p.PointString = loc

                    /* INSERT */
                    if err := c.Insert(p); err != nil {
                            e := err.Error()

                            log.Println("Database error: " + e)
                            return
                    }

                    log.Println("EMIT: ", so.Emit("demonight", loc))
                })

		so.On("disconnection", func() {
                    log.Println("DISCONNECTED")
		})
	})

        server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

        /* ROUTES */
        http.Handle("/demonight/socket/", server)

        /* START SERVER */
        log.Println("OPEN HOUSE SERVER LISTENING ON PORT 3000...")
        if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("Could not start server: ", err)
	}
}

