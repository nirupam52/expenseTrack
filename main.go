package main
import (
	"fmt"
	"log"
	"net/http"

	"github.com/nirupam52/expenseTrack/internal/config"
)

func main() {
	appConfig := config.LoadConfig()

	fmt.Println("App Config:", appConfig)
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", appConfig.Port), nil))
}