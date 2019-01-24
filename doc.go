/*
NPM-GO is a Go client for the NPM HTTP API based on rabbithole (https://github.com/michaelklishin/rabbit-hole)

All HTTP API operations are accessible via `npmgo.Client`, which
should be instantiated with `npmgo.NewClient`.

        // URI, username, password
        rmqc, _ = NewClient("http://127.0.0.1:15672", "guest", "guest")
*/
package npmgo
