package main

import (
	"encoding/csv"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// sticky represents one note in the database
type sticky struct {
	ID    int    `json:"id"`
	Month int    `json:"month"`
	Year  int    `json:"year"`
	Text  string `json:"text"`
}

var stickies map[int]sticky

const hostAddress = "localhost:8484"

//
// MAIN FUNCTION
//

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: ./daffodil-data <filename>")
		return
	}

	filename := os.Args[1]
	file, err := os.OpenFile(filename, os.O_RDONLY, 0)

	if err != nil {
		log.Fatal("could not read file", filename, err.Error())
		return
	}

	err = initStickies(file)
	if err != nil {
		log.Fatal("error initializing stickies", err.Error())
		return
	}

	router := gin.Default()
	router.GET("/stickies", handleGetAllStickies)
	router.GET("/stickies/:id", handleGetStickyById)
	router.GET("/stickies/random", handleGetRandomSticky)

	router.Run(hostAddress)
}

//
// HANDLER FUNCTIONS
//

func handleGetAllStickies(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getAllStickies())
}

func handleGetStickyById(c *gin.Context) {
	idParam := c.Param("id")

	id, idParamErr := strconv.Atoi(idParam)
	if idParamErr != nil {
		c.IndentedJSON(http.StatusBadRequest, "id not numeric: "+idParam)
		return
	}

	fetchedSticky, okFetch := stickies[id]
	if !okFetch {
		c.IndentedJSON(http.StatusNotFound, "id not found: "+idParam)
		return
	}

	c.IndentedJSON(http.StatusOK, fetchedSticky)
}

func handleGetRandomSticky(c *gin.Context) {
	var randSticky sticky
	ok := false
	for !ok {
		randSticky, ok = stickies[rand.Intn(len(stickies))]
	}
	c.IndentedJSON(http.StatusOK, randSticky)
}

//
// MEMORY FUNCTIONS
//

func getAllStickies() (allStickies []sticky) {
	allStickies = make([]sticky, len(stickies))
	for i, s := range stickies {
		allStickies[i] = s
	}
	return
}

func getSticky(id int) sticky {
	return stickies[id]
}

func initStickies(file *os.File) error {
	stickies = make(map[int]sticky)

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal("could not read records", err.Error())
		return err
	}

	for _, record := range records {
		id, _ := strconv.Atoi(record[0])
		month, _ := strconv.Atoi(record[1])
		year, _ := strconv.Atoi(record[2])
		text := record[3]

		if id != len(stickies) || month < 1 || month > 12 || year < 0 {
			log.Fatal("invalid sticky file. row:", id, month, year, text)
		}

		newSticky := sticky{id, month, year, text}

		stickies[id] = newSticky
	}

	return nil
}
