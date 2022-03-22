package main

import (
	"bufio"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

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

var id = 0

type HttpErrorMessage struct {
	MESSAGE string `json:"message" binding:"required"`
}

type HttpMessage struct {
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

// TODO: definition lists
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

// @title           ReadmeBuilder API
// @version         1.0
// @description     This is a API to be used for creating markdown files
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

// CreateReadme godoc
// @Summary Creates a readme
// @Description Creates a readme object can now add markdown elements
// @Accept json
// @Produce json
// @Param	name	query	string	false	"pass a value to create a user defined readmeId"
// @Success	201		{object}	HttpMessage	"returns a message with the readmeId"
// @Failure	409		{object}	HttpErrorMessage	"Readme already exists"
// @Router	/readme	[post]
func createReadme(c *gin.Context) {
	var readmeId = ""

	if c.Query("name") == "" {
		readmeId = uuid.NewString()
	} else {
		readmeId = c.Query("name")
	}

	if len(readmeDB[readmeId]) >= 1 {
		c.IndentedJSON(http.StatusConflict, HttpErrorMessage{MESSAGE: "Readme with that id already exists"})
		return
	}

	readmeDB[readmeId] = append(readmeDB[readmeId], "")

	c.IndentedJSON(http.StatusCreated, HttpMessage{MESSAGE: readmeId})
}

// CreateReadmeFile godoc
// @Summary Creates markdown file
// @Description From all of your previous operations takes the readme and generates the markdown file
// @Accept json
// @Produce json
// @Param	id	path	string	true	"readme id"
// @Success 200
// @Router 	/readme/{id}/file	[post]
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

// GetReadme godoc
// @Summary Returns a readme
// @Accept json
// @Produce json
// @Param	id	path	string	true	"readme id"
// @Success	200	{array}		string	"list of markdown strings"
// @Failure 404	{object}	HttpErrorMessage	"could not find readme"
// @Router	/readme/{id}		[get]
func getReadme(c *gin.Context) {
	readmeId := c.Param("id")

	if len(readmeDB[readmeId]) < 1 {
		c.IndentedJSON(http.StatusNotFound, HttpErrorMessage{MESSAGE: "could not find readme"})
		return
	}

	c.IndentedJSON(http.StatusOK, readmeDB[readmeId])
}

// change to read file from s3
func decodeReadme(c *gin.Context) {
	readmeId := c.Param("id")

	if len(readmeDB[readmeId]) < 1 {
		c.IndentedJSON(http.StatusNotFound, HttpErrorMessage{MESSAGE: "could not find readme"})
		return
	}

	readme := readmeDB[readmeId]

	currentReadmeDecoded := ``
	for _, line := range readme {
		decodedReadmeLine, err := strconv.Unquote(`"` + line + `"`)
		if err != nil {
			panic(err)
		}
		currentReadmeDecoded = currentReadmeDecoded + decodedReadmeLine
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": currentReadmeDecoded})
}

// AddHeader godoc
// @Summary Adds Header
// @Description Creates a string to be used for a markdown header
// @Accept json
// @Produce json
// @Param	id	path	string	true	"readme id"
// @Param	addHeader	body	AddHeaderRequest	true	"request body for header"
// @Success	200	{object}	HttpMessage	"returns the header markdown string"
// @Failure 404	{object}	HttpErrorMessage	"could not find readme"
// @Failure 400	{object}	HttpErrorMessage	"incorrect request body"
// @Router	/readme/{id}/header	[put]
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

	c.IndentedJSON(http.StatusOK, HttpMessage{MESSAGE: createdString})
}

// AddParagraph godoc
// @Summary Adds a paragraph
// @Description Updates readme to have a paragraph
// @Accept json
// @Produce json
// @Param	id	path	string	true	"readme id"
// @Param	paragraph	query	string	true	"paragraph you want to add to the readme"
// @Success	200	{object}	HttpMessage	"returns an paragraph markdown string"
// @Failure 404	{object}	HttpErrorMessage	"could not find readme"
// @Failure 400	{object}	HttpErrorMessage	"paragraph param cannot be empty"
// @Router	/readme/{id}/paragraph	[put]
func addParagraph(c *gin.Context) {
	readmeId := c.Param("id")
	paragraph := c.Query("paragraph")

	if len(readmeDB[readmeId]) < 1 {
		c.IndentedJSON(http.StatusNotFound, HttpErrorMessage{MESSAGE: "could not find readme"})
		return
	}

	if strings.TrimSpace(paragraph) == "" {
		c.IndentedJSON(http.StatusBadRequest, HttpErrorMessage{MESSAGE: "paragraph cannot be empty"})
		return
	}

	paragraph = paragraph + "\n"

	readmeDB[readmeId] = append(readmeDB[readmeId], paragraph)

	c.IndentedJSON(http.StatusOK, HttpMessage{MESSAGE: paragraph})
}

// AddCode godoc
// @Summary Adds code to readme
// @Description creates a string in markdown code block with the language specified
// @Accept json
// @Produce json
// @Param	id	path	string	true	"readme id"
// @Param	codeRequest	body	AddCodeRequest	true	"request for code markdown"
// @Success	200	{object}	HttpMessage	"returns the created code markdown string"
// @Failure 404	{object}	HttpErrorMessage	"could not find readme"
// @Failure 400	{object}	HttpErrorMessage	"incorrect request body"
// @Failure 400	{object}	HttpErrorMessage	"the code language is not suppored"
// @Router	/readme/{id}/code	[put]
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

		c.IndentedJSON(http.StatusOK, HttpMessage{MESSAGE: createdCodeString})
		return
	}

	c.IndentedJSON(http.StatusBadRequest, HttpErrorMessage{MESSAGE: "Code language not supported"})
}

// AddBlockquote godoc
// @Summary Add Blockquote
// @Description creates a blockquote markdown string
// @Accept json
// @Produce	json
// @Param	id	path	string	true	"readme id"
// @Param	paragraph	query	string	true	"string for paragraph markdown"
// @Success	200	{object}	HttpMessage	"returns created markdown blockquote string"
// @Failure 404	{object}	HttpErrorMessage	"could not find readme"
// @Failure 400	{object}	HttpErrorMessage	"blockquote can not be empty"
// @Router	/readme/{id}/blockquote	[put]
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

	c.IndentedJSON(http.StatusOK, HttpMessage{MESSAGE: createdBlockquote})
}

// AddLink godoc
// @Summary Add Link
// @Description	creates a markdown link string
// @Accept json
// @Produce json
// @Param	id	path	string	true	"readme id"
// @Param	addLinkRequest	body	AddLinkRequest	true	"request for adding link"
// @Success	200	{object}	HttpMessage	"returns created markdown link"
// @Failure 404	{object}	HttpErrorMessage	"could not find readme"
// @Failure 400	{object}	HttpErrorMessage	"incorrect request body"
// @Router	/readme/{id}/link	[put]
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

	c.IndentedJSON(http.StatusOK, HttpMessage{MESSAGE: createdLink})
}

// AddImage godoc
// @Summary Add Image
// @Description	creates a markdown image string
// @Accept json
// @Produce json
// @Param	id	path	string	true	"readme id"
// @Param	addLinkRequest	body	AddLinkRequest	true	"request body for adding image"
// @Success	200	{object}	HttpMessage	"returns created markdown image link"
// @Failure 404	{object}	HttpErrorMessage	"could not find readme"
// @Failure 400	{object}	HttpErrorMessage	"incorrect request body"
// @Router	/readme/{id}/image	[put]
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

	c.IndentedJSON(http.StatusOK, HttpMessage{MESSAGE: createdImage})
}

// AddTable godoc
// @Summary Add Table
// @Description	creates a markdown table as a string
// @Accept json
// @Produce json
// @Param	id	path	string	true	"readme id"
// @Param	addTableRequest	body	AddTableRequest	true	"request table body"
// @Success	200	{object}	HttpMessage	"returns table markdown string with values inserted"
// @Failure 404	{object}	HttpErrorMessage	"could not find readme"
// @Failure 400	{object}	HttpErrorMessage	"incorrect request body"
// @Router	/readme/{id}/table	[put]
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

	c.IndentedJSON(http.StatusOK, HttpMessage{MESSAGE: createdTableString})
}
