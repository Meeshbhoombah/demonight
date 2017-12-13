/*

    demonight.go
    ENTRY POINT FOR DAAANK DEMO NIGHT APP

*/

package main

import (
        "log"
        "net/http"

        "github.com/googollee/go-socket.io"
)

func main() {
        server, err := socketio.NewServer(nil)
        if err != nil {
                log.Fatal(err)
        }

        server.On("connection", func(so socketio.Socket) {
                log.Println("on connection")
                so.Join("chat")

                // recieved phone # from client
                so.On("phone-number", func(msg string) {
                        log.Println(msg)

                        // send phone # back to client
                        emit := so.Emit("phone-number", msg)
                        log.Println("Emit: ", emit)
                        so.BroadcastTo("users", "phone-number", msg)
                })

                so.On("disconnection", func() {
                        log.Println("on disconnect")
                })
        })

        server.On("error", func(so socketio.Socket, err error) {
                log.Println("error:", err)
        })

        http.Handle("/demonight/socket/", server)

        /* START SERVER */
        if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("Could not start server: ", err)
	}

	log.Println("Listenting on PORT 3000")
}

