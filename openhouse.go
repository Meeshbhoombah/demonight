/*
main.go

Backend for Daaank's Open House application
*/

package main

import (
        "log"
        "net/http"

        "github.com/googollee/go-socket.io"
)

func main() {

        // create new connection
        server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

        server.On("connection", func(so socketio.Socket) {
                log.Println("CONNECTED")

                so.Join("demonight-realtime")
                so.On("location", func(loc string) {
                    log.Println("EVENT 'location': ", loc)
                    so.BroadcastTo("demonight-realtime", "location", loc)
                })

		so.On("disconnection", func() {
                    log.Println("DISCONNECTED")
		})
	})

        server.On("location", func(loc string) {
                log.Println("EVENT 'location' contains: ", loc)
        })

        server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

        /* ROUTES */
        http.Handle("/demonight/socket/", server)

        /* START SERVER */
        log.Println("STARTING OPEN HOUSE SERVER @ http://localhost:3000")
        if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("Could not start server: ", err)
	}
}

