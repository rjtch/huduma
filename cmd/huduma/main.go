package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/huduma/cmd/huduma/handlers"

	"github.com/huduma/internal/mongo"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	server *http.Server
	db     *mongo.BooksDB
)

var versions string

var (
	rootCommand = &cobra.Command{
		Use:   "huduma",
		Short: "huduma is an awesome server",
		Long: ` 
		Huduma is a http restfull webserver that can manage incoming request to a mongodb database.
		It was implemented to learn how to organize code in Golang. For e.g how organize packages.  
		The architecture and pattern used in this project and also more than 50% of code is 
		from William kennedy's classes.A Golang expert and software architect.
		`,
		Run: RunServer,
	}

	versionCmd = cobra.Command{
		Run: version,
		Use: "version",
	}
)

func version(cmd *cobra.Command, args []string) {
	fmt.Println("1.0.0", versions)
}

//initialization loop
//Init initialize the rootcommand
func Init() {
	cobra.OnInitialize()
	rootCommand.PersistentFlags().StringP("config", "C", "", "An explicit config file to use")
	rootCommand.Flags().IntP("port", "p", 0, "the port to do things on")
	rootCommand.AddCommand(&versionCmd)

}

//RunServer is a our cobra command
func RunServer(cmd *cobra.Command, args []string) {

	//Init()
	log.Info("Huduma is starting ...")
	log.Info("Start initializing mongo")
	dbDialTimeout := 200 * time.Second
	//shutDownTimeout := 5 * time.Second
	host := os.Getenv("HOST")
	if host == "" {
		host = ":3000"
	}

	db, err := mongo.NewCollection(host, dbDialTimeout)
	if err != nil {
		log.Fatalf("error when registring DB: %v", err)
	}
	defer db.Close()

	//address := fmt.Sprintf(":%v", globalConfig.Port)
	server = &http.Server{
		Addr:           host,
		Handler:        handlers.API(db),
		ReadTimeout:    3 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		log.Printf("Start : Listening %s", host)
		log.Printf("shutdown : Listener closed %v", server.ListenAndServe())
		wg.Done()
	}()

	log.Info("Huduma started successfully")

	wg.Wait()
	log.Println("main : Conpleted")

	/**
	I will add later more logic hier to enable huduma to handle requests on incoming TLS connection
	**/
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}
}
