package main

import (
	"bufio"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// change to heading map
var headingSyntaxMap = map[string]string{
	"SMALL_HEADING":  "### ",
	"MEDIUM_HEADING": "## ",
	"LARGE_HEADING":  "# ",
}

var codeLanguageMap = map[string]bool{
	"go":   true,
	"java": true,
	"json": true,
}

type HttpErrorMessage struct {
	MESSAGE string `json:"message" binding:"required"`
}

type MarkDownRequest struct {
	MDTYPE  string   `json:"mdtype" binding:"required"`
	MDVALUE []string `json:"mdvalue" binding:"required"`
}

type ReadmeResponse struct {
	NAME   string   `json:"name" binding:"required"`
	VALUES []string `json:"values" binding:"required"`
}

type AddHeaderRequest struct {
	HEADER_TYPE string `json:"header_type" binding:"required"`
	VALUE       string `json:"value" binding:"required"`
}

type AddCodeRequest struct {
	CODE_LANGUAGE string `json:"code_language" binding:"required"`
	VALUE         string `json:"value" binding:"required"`
}

type AddLinkRequest struct {
	DESCRIPTION string `json:"description" binding:"required"`
	LINK        string `json:"link" binding:"required"`
}

type AddTableRequest struct {
	COLUMN_NAMES  []string            `json:"column_names" binding:"required"`
	COLUMN_VALUES map[string][]string `json:"column_values" binding:"required"`
}

var readmeDB = make(map[string][]string)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//tables
//definition lists
//paragraph
func setupRouter() *gin.Engine {
	router := gin.New()

	router.POST("/readme", createReadme)
	router.GET("/readme/:id", getReadme)
	router.PUT("/readme/:id/header", addHeader)
	router.PUT("/readme/:id/paragraph", addParagraph)
	router.PUT("/readme/:id/code", addCode)
	router.PUT("/readme/:id/blockquote", addBlockquote)
	router.PUT("/readme/:id/link", addLink)
	router.PUT("/readme/:id/image", addImage)
	router.PUT("/readme/:id/table", addTable)
	router.POST("/readme/:id/file", createReadmeFile)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}

// @title           Survey Voting API
// @version         1.0
// @description     This is a Survey Voting API
// @contact.name   Matthew Croft
// @contact.url    https://www.linkedin.com/in/matthew-croft-44a5a5b3/
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /
func main() {
	router := setupRouter()

	router.Run("localhost:8080")
}

// create constants
// write test
// produce file

func createReadme(c *gin.Context) {
	readmeName := c.Query("name")

	readmeDB[readmeName] = append(readmeDB[readmeName], "")
}

func createReadmeFile(c *gin.Context) {
	readmeId := c.Param("id")
	f, err := os.Create("/tmp/readme.md")
	check(err)

	//write buffer
	wr := bufio.NewWriter(f)

	var lines = readmeDB[readmeId]

	for _, line := range lines {
		if _, err := wr.Write([]byte(line)); err != nil {
			panic(err)
		}
	}

	if err = wr.Flush(); err != nil {
		panic(err)
	}
}

func getReadme(c *gin.Context) {
	readmeId := c.Param("id")

	if len(readmeDB[readmeId]) < 1 {
		c.IndentedJSON(http.StatusNotFound, HttpErrorMessage{MESSAGE: "could not find readme"})
		return
	}

	c.IndentedJSON(http.StatusOK, readmeDB[readmeId])
}

func addHeader(c *gin.Context) {
	readmeId := c.Param("id")
	var addHeaderRequest AddHeaderRequest

	if err := c.BindJSON(&addHeaderRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, HttpErrorMessage{MESSAGE: "incorrect request body, should be AddHeaderRequest body"})
		return
	}

	if len(readmeDB[readmeId]) < 1 {
		c.IndentedJSON(http.StatusNotFound, HttpErrorMessage{MESSAGE: "could not find readme"})
		return
	}

	headerMarkdown := headingSyntaxMap[addHeaderRequest.HEADER_TYPE]

	createdString := headerMarkdown + addHeaderRequest.VALUE + "\n"

	readmeDB[readmeId] = append(readmeDB[readmeId], createdString)

	c.IndentedJSON(http.StatusOK, gin.H{"message": createdString})
}

func addParagraph(c *gin.Context) {
	readmeId := c.Param("id")
	paragraph := c.Query("paragraph")

	if len(readmeDB[readmeId]) < 1 {
		c.IndentedJSON(http.StatusNotFound, HttpErrorMessage{MESSAGE: "could not find readme"})
		return
	}

	if strings.TrimSpace(paragraph) == "" {
		c.IndentedJSON(http.StatusBadRequest, HttpErrorMessage{MESSAGE: "paragraph cannot be empty"})
	}

	paragraph = paragraph + "\n"

	readmeDB[readmeId] = append(readmeDB[readmeId], paragraph)

}

func addCode(c *gin.Context) {
	readmeId := c.Param("id")
	var addCodeRequest AddCodeRequest

	if err := c.BindJSON(&addCodeRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, HttpErrorMessage{MESSAGE: "incorrect request body, should be AddCodeRequest body"})
		return
	}

	if len(readmeDB[readmeId]) < 1 {
		c.IndentedJSON(http.StatusNotFound, HttpErrorMessage{MESSAGE: "could not find readme"})
		return
	}

	if codeLanguageMap[addCodeRequest.CODE_LANGUAGE] {
		createdCodeString := "```" + addCodeRequest.CODE_LANGUAGE + "\n " + addCodeRequest.VALUE + "```" + "\n"

		readmeDB[readmeId] = append(readmeDB[readmeId], createdCodeString)

		c.IndentedJSON(http.StatusOK, gin.H{"message": createdCodeString})
		return
	}

	c.IndentedJSON(http.StatusBadRequest, HttpErrorMessage{MESSAGE: "Code language not supported"})
}

func addBlockquote(c *gin.Context) {
	readmeId := c.Param("id")
	message := c.Query("blockquote")

	if strings.TrimSpace(message) == "" {
		c.IndentedJSON(http.StatusBadRequest, HttpErrorMessage{MESSAGE: "blockquote cannot be empty"})
		return
	}

	createdBlockquote := "> " + message + "\n"

	if len(readmeDB[readmeId]) < 1 {
		c.IndentedJSON(http.StatusNotFound, HttpErrorMessage{MESSAGE: "could not find readme"})
		return
	}

	readmeDB[readmeId] = append(readmeDB[readmeId], createdBlockquote)

	c.IndentedJSON(http.StatusOK, gin.H{"message": createdBlockquote})
}

func addLink(c *gin.Context) {
	readmeId := c.Param("id")
	var addLinkRequest AddLinkRequest

	if err := c.BindJSON(&addLinkRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, HttpErrorMessage{MESSAGE: "incorrect request body, should be AddLinkRequest body"})
		return
	}

	if len(readmeDB[readmeId]) < 1 {
		c.IndentedJSON(http.StatusNotFound, HttpErrorMessage{MESSAGE: "could not find readme"})
		return
	}

	createdLink := "[" + addLinkRequest.DESCRIPTION + "]" + "(" + addLinkRequest.LINK + ")" + "\n"

	readmeDB[readmeId] = append(readmeDB[readmeId], createdLink)

	c.IndentedJSON(http.StatusOK, gin.H{"message": createdLink})
}

func addImage(c *gin.Context) {
	readmeId := c.Param("id")
	var addImageRequest AddLinkRequest

	if err := c.BindJSON(&addImageRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, HttpErrorMessage{MESSAGE: "incorrect request body, should be AddLinkRequest body"})
		return
	}

	if len(readmeDB[readmeId]) < 1 {
		c.IndentedJSON(http.StatusNotFound, HttpErrorMessage{MESSAGE: "could not find readme"})
		return
	}

	createdImage := "![" + addImageRequest.DESCRIPTION + "]" + "(" + addImageRequest.LINK + ")" + "\n"

	readmeDB[readmeId] = append(readmeDB[readmeId], createdImage)

	c.IndentedJSON(http.StatusOK, gin.H{"message": createdImage})
}

//TODO: finish addtable
//find largest column values, iterate through make sure dont go
//out of bounds on array
func addTable(c *gin.Context) {
	readmeId := c.Param("id")
	var addTableRequest AddTableRequest

	if err := c.BindJSON(&addTableRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, HttpErrorMessage{MESSAGE: "incorrect request body, should be AddTableRequest body"})
		return
	}

	if len(readmeDB[readmeId]) < 1 {
		c.IndentedJSON(http.StatusNotFound, HttpErrorMessage{MESSAGE: "could not find readme"})
		return
	}

	largestColumn := 0
	createdTableString := `|`

	//create column label row
	for _, cName := range addTableRequest.COLUMN_NAMES {
		if len(addTableRequest.COLUMN_VALUES[cName]) > largestColumn {
			largestColumn = len(addTableRequest.COLUMN_VALUES[cName])
		}
		createdTableString = createdTableString + cName + `|`
	}

	createdTableString = createdTableString + "\n" + `|`

	// create separated between column label and column values
	for i := 0; i < len(addTableRequest.COLUMN_NAMES); i++ {
		createdTableString = createdTableString + ` --- |`
	}

	createdTableString = createdTableString + "\n"

	// values in each column
	for i := 0; i < largestColumn; i++ {
		currentString := `|`
		for _, column_name := range addTableRequest.COLUMN_NAMES {
			if i < len(addTableRequest.COLUMN_VALUES[column_name]) {
				currentString = currentString + addTableRequest.COLUMN_VALUES[column_name][i] + `|`
			} else {
				currentString = currentString + " " + `|`
			}
		}
		createdTableString = createdTableString + currentString + "\n"
	}

	readmeDB[readmeId] = append(readmeDB[readmeId], createdTableString)

	c.IndentedJSON(http.StatusOK, gin.H{"message": createdTableString})
}
