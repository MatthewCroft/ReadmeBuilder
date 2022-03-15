package main

import (
	"bufio"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var markdownSyntaxMap = map[string]string{
	"SMALL_HEADING":  "### ",
	"MEDIUM_HEADING": "## ",
	"LARGE_HEADING":  "# ",
	"BLOCKQUOTE":     "> ",
	"CODE":           "```",
	"LINK":           "[]() ",
	"IMAGE":          "![]() ",
}

var codeLanguageMap = map[string]bool{
	"go":   true,
	"java": true,
	"json": true,
}

type HttpErrorMessage struct {
	MESSAGE string `json:"message"`
}

type MarkDownRequest struct {
	MDTYPE  string   `json:"mdtype"`
	MDVALUE []string `json:"mdvalue"`
}

type ReadmeResponse struct {
	NAME   string   `json:"name"`
	VALUES []string `json:"values"`
}

type AddHeaderRequest struct {
	HEADER_TYPE string `json:"header_type"`
	VALUE       string `json:"value"`
}

type AddCodeRequest struct {
	CODE_LANGUAGE string `json:"code_language"`
	VALUE         string `json:"value"`
}

type AddLinkRequest struct {
	DESCRIPTION string `json:"description"`
	LINK        string `json:"link"`
}

var readmeDB = make(map[string][]string)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func setupRouter() *gin.Engine {
	router := gin.New()

	router.POST("/readme", createReadme)
	router.GET("/readme/:id", getReadme)
	router.PUT("/readme/:id/header", addHeader)
	router.PUT("/readme/:id/code", addCode)
	router.PUT("/readme/:id/blockquote", addBlockquote)
	router.PUT("/readme/:id/link", addLink)
	router.PUT("/readme/:id/image", addImage)
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

	headerMarkdown := markdownSyntaxMap[addHeaderRequest.HEADER_TYPE]

	createdString := headerMarkdown + addHeaderRequest.VALUE + "\n"

	readmeDB[readmeId] = append(readmeDB[readmeId], createdString)

	c.IndentedJSON(http.StatusOK, createdString)
}

func addCode(c *gin.Context) {
	readmeId := c.Param("id")
	var addCodeRequest AddCodeRequest

	if err := c.BindJSON(&addCodeRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, HttpErrorMessage{MESSAGE: "incorrect request body, should be AddCodeRequest body"})
		return
	}

	if codeLanguageMap[addCodeRequest.CODE_LANGUAGE] {
		createdCodeString := markdownSyntaxMap["CODE"] + addCodeRequest.CODE_LANGUAGE + "\n " + addCodeRequest.VALUE + "```" + "\n"

		readmeDB[readmeId] = append(readmeDB[readmeId], createdCodeString)

		c.IndentedJSON(http.StatusOK, createdCodeString)
		return
	}

	c.IndentedJSON(http.StatusBadRequest, HttpErrorMessage{MESSAGE: "Code language not supported"})
}

func addBlockquote(c *gin.Context) {
	readmeId := c.Param("id")
	message := c.Query("blockquote")

	createdBlockquote := markdownSyntaxMap["BLOCKQUOTE"] + message + "\n"

	if len(readmeDB[readmeId]) < 1 {
		c.IndentedJSON(http.StatusNotFound, HttpErrorMessage{MESSAGE: "could not find readme"})
		return
	}

	readmeDB[readmeId] = append(readmeDB[readmeId], createdBlockquote)

	c.IndentedJSON(http.StatusOK, createdBlockquote)
}

func addLink(c *gin.Context) {
	readmeId := c.Param("id")
	var addLinkRequest AddLinkRequest

	if err := c.BindJSON(&addLinkRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, HttpErrorMessage{MESSAGE: "incorrect request body, should be AddLinkRequest body"})
	}

	if len(readmeDB[readmeId]) < 1 {
		c.IndentedJSON(http.StatusNotFound, HttpErrorMessage{MESSAGE: "could not find readme"})
		return
	}

	createdLink := "[" + addLinkRequest.DESCRIPTION + "]" + "(" + addLinkRequest.LINK + ")" + "\n"

	readmeDB[readmeId] = append(readmeDB[readmeId], createdLink)

	c.IndentedJSON(http.StatusOK, createdLink)
}

func addImage(c *gin.Context) {
	readmeId := c.Param("id")
	var addImageRequest AddLinkRequest

	if err := c.BindJSON(&addImageRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, HttpErrorMessage{MESSAGE: "incorrect request body, should be AddLinkRequest body"})
	}

	if len(readmeDB[readmeId]) < 1 {
		c.IndentedJSON(http.StatusNotFound, HttpErrorMessage{MESSAGE: "could not find readme"})
		return
	}

	createdImage := "![" + addImageRequest.DESCRIPTION + "]" + "(" + addImageRequest.LINK + ")" + "\n"

	readmeDB[readmeId] = append(readmeDB[readmeId], createdImage)

	c.IndentedJSON(http.StatusOK, createdImage)
}
