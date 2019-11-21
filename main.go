package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    _ "gopkg.in/yaml.v2"
    "io/ioutil"
    "time"
    "strconv"
    "fmt"
)

type config struct {
    PeerPort int64 `yaml:peerPort`
    AliveCheckInterval int64 `yaml:aliveCheckInterval`
}

// Check if all peers are still alive
func checkAlive(peers *map[string]bool, aliveCheckDuration time.Duration, c config) {
    var aliveCount int
    for {
        aliveCount = len(*peers)
        fmt.Printf("\nChecking if (%v) peers are alive\n", aliveCount)
        for ip := range *peers {
            resp, err := http.Get("http://" + ip + ":" + strconv.FormatInt(c.PeerPort, 10))
            if err != nil {
                delete(*peers, ip)
            } else {
                defer resp.Body.Close()
                body, err := ioutil.ReadAll(resp.Body)
                if err != nil {

                } else {
                    fmt.Printf("%v", body)
                    // if body.dancing == "false" {
                    //     delete(peers, ip)
                    // }
                }
            }
        }
        fmt.Printf("%v peers died\n", (aliveCount - len(*peers)))
        time.Sleep(aliveCheckDuration)
    }
}

func main() {
    // TODO: Get config from file
    var c config
    c.PeerPort = 6666
    c.AliveCheckInterval = 60

    formattedAliveCheckInterval := (strconv.FormatInt(c.AliveCheckInterval, 10)+ "s")
    aliveCheckDuration, _ := time.ParseDuration(formattedAliveCheckInterval)

    // Keep track of alive peers.
    var peers map[string]bool
    peers = make(map[string]bool)

    // Set the router as the default one shipped with Gin
    router := gin.Default()

    disco := router.Group("/disco")
    {
        disco.GET("/", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H {
                "nodes": peers,
            })
        })

        disco.GET("/join", func(c *gin.Context) {
            peers[c.ClientIP()] = true

            c.JSON(http.StatusOK, gin.H {
                "status": "You've joined the disco",
            })
        })

        disco.GET("/leave", func(c *gin.Context) {
            delete(peers, c.ClientIP())

            c.JSON(http.StatusOK, gin.H {
                "status": "You've left the disco",
            })
        })
    }
    go checkAlive(&peers, aliveCheckDuration, c)
    router.Run(":80")
}
