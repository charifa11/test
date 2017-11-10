
package main 
 import "github.com/charifa11/test/packageone"

func main() {
    // this is really fast, so we will need the microsecs in the logs to see something :D
    log.SetFlags(log.LstdFlags | log.Lmicroseconds)

    // run pipeline
    packageone.RunPipeline()
}