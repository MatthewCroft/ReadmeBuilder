package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// router.POST("/readme", createReadme)
// router.GET("/readme/:id", getReadme)
// router.PUT("/readme/:id/header", addHeader)
// router.PUT("/readme/:id/code", addCode)
// router.PUT("/readme/:id/blockquote", addBlockquote)
// router.PUT("/readme/:id/link", addLink)
// router.PUT("/readme/:id/image", addImage)
// router.POST("/readme/:id/file", createReadmeFile)

func TestGetReadme(t *testing.T) {
	router := setupRouter()
	r := httptest.NewRecorder()
	w := httptest.NewRecorder()

	req1, _ := http.NewRequest("POST", "/readme?name=1", nil)
	router.ServeHTTP(r, req1)

	emptyReadme := string(`[""]`)

	req, _ := http.NewRequest("GET", "/readme/1", nil)
	router.ServeHTTP(w, req)

	require.JSONEq(t, emptyReadme, w.Body.String())
}

func TestGetReadmeReturnsNotFoundReadme(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/readme/INVALID", nil)
	router.ServeHTTP(w, req)

	require.JSONEq(t, string(`{"message": "could not find readme"}`), w.Body.String())
}

func TestAddHeader(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	r := httptest.NewRecorder()

	req1, _ := http.NewRequest("POST", "/readme?name=3", nil)
	router.ServeHTTP(w, req1)

	var headerRequest = []byte(`{
			"header_type": "SMALL_HEADING",
			"value": "My first header"
		}`)

	req2, _ := http.NewRequest("PUT", "/readme/3/header", bytes.NewBuffer(headerRequest))
	router.ServeHTTP(r, req2)

	require.JSONEq(t, string(`{ "message": "### My first header\n"}`), r.Body.String())
}

func TestAddHeaderReturnsNotFoundReadme(t *testing.T) {
	router := setupRouter()
	r := httptest.NewRecorder()

	var headerRequest = []byte(`{
		"header_type": "SMALL_HEADING",
		"value": "My first header"
	}`)

	req2, _ := http.NewRequest("PUT", "/readme/4/header", bytes.NewBuffer(headerRequest))
	router.ServeHTTP(r, req2)

	require.JSONEq(t, string(`{ "message":"could not find readme"}`), r.Body.String())
}

func TestAddHeaderReturnsIncorrectRequestBody(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	r := httptest.NewRecorder()

	req1, _ := http.NewRequest("POST", "/readme?name=4", nil)
	router.ServeHTTP(w, req1)

	var headerRequest1 = []byte(`{
	}`)

	req2, _ := http.NewRequest("PUT", "/readme/4/header", bytes.NewBuffer(headerRequest1))
	router.ServeHTTP(r, req2)

	require.JSONEq(t, string(`{ "message":"incorrect request body, should be AddHeaderRequest body"}`), r.Body.String())
}

// TODO: fix test to return bytes, having trouble with ```
// func TestAddCode(t *testing.T) {
// 	router := setupRouter()
// 	w := httptest.NewRecorder()
// 	r := httptest.NewRecorder()

// 	value := "`" + "```go\n func createReadme(c *gin.Context) {\r\n\treadmeName := c.Query(\"name\")\r\n\r\n\treadmeDB[readmeName] = append(readmeDB[readmeName], \"\")\r\n}\r\n```\n" + "`"

// 	req1, _ := http.NewRequest("POST", "/readme?name=5", nil)
// 	router.ServeHTTP(w, req1)

// 	headerRequest := []byte(`{
// 		"code_language": "go",
// 		"value": `"func createReadme(c *gin.Context) {\r\n\treadmeName := c.Query(\"name\")\r\n\r\n\treadmeDB[readmeName] = append(readmeDB[readmeName], \"\")\r\n}\r\n"`
// 	}`)

// 	req2, _ := http.NewRequest("PUT", "/readme/5/code", bytes.NewBuffer(headerRequest))
// 	router.ServeHTTP(r, req2)

// 	// message := "  \"```" + "go\n func createReadme(c *gin.Context) {\r\n\treadmeName := c.Query(\"name\")\r\n\r\n\treadmeDB[readmeName] = append(readmeDB[readmeName], \"\")\r\n}\r\n" + "```\""
// 	// ans := "{" + `"message": "` + bytes[]{"```go\n func createReadme(c *gin.Context) {\r\n\treadmeName := c.Query(\"name\")\r\n\r\n\treadmeDB[readmeName] = append(readmeDB[readmeName], \"\")\r\n}\r\n```\n" + "}"

// 	// fmt.Println(r.Body)
// 	// fmt.Println(ans)

// 	if gin.H{"": ""} != r.Body {
// 		t.Fatalf("Code was not added to the readme correctly")
// 	}
// }

func TestAddBlockquote(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	r := httptest.NewRecorder()

	req1, _ := http.NewRequest("POST", "/readme?name=7", nil)
	router.ServeHTTP(w, req1)

	req2, _ := http.NewRequest("PUT", "/readme/7/blockquote?blockquote=This is an important Quote", nil)
	router.ServeHTTP(r, req2)

	require.JSONEq(t, string(`{ "message": "> This is an important Quote\n" }`), r.Body.String())
}

func TestAddBlockquoteReturnsNotFoundReadme(t *testing.T) {
	router := setupRouter()
	r := httptest.NewRecorder()

	req1, _ := http.NewRequest("PUT", "/readme/10/blockquote?blockquote=This is an important Quote", nil)
	router.ServeHTTP(r, req1)

	require.JSONEq(t, string(`{ "message":"could not find readme" }`), r.Body.String())
}

func TestAddBlockquoteReturnsEmptyBlockquote(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	r := httptest.NewRecorder()

	req1, _ := http.NewRequest("POST", "/readme?name=11", nil)
	router.ServeHTTP(w, req1)

	req2, _ := http.NewRequest("PUT", "/readme/11/blockquote?blockquote=", nil)
	router.ServeHTTP(r, req2)

	require.JSONEq(t, string(`{ "message":"blockquote cannot be empty" }`), r.Body.String())
}

func TestAddLink(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	r := httptest.NewRecorder()

	var headerRequest = []byte(`{
		"link": "https://go.dev/doc/",
		"description": "Go Dev"
	}`)

	req1, _ := http.NewRequest("POST", "/readme?name=12", nil)
	router.ServeHTTP(w, req1)

	req2, _ := http.NewRequest("PUT", "/readme/12/link", bytes.NewBuffer(headerRequest))
	router.ServeHTTP(r, req2)

	require.JSONEq(t, string(`{ "message": "[Go Dev](https://go.dev/doc/)\n" }`), r.Body.String())
}

func TestAddLinkReturnsReadmeNotFound(t *testing.T) {
	router := setupRouter()
	r := httptest.NewRecorder()

	var headerRequest = []byte(`{
		"link": "https://go.dev/doc/",
		"description": "Go Dev"
	}`)

	req1, _ := http.NewRequest("PUT", "/readme/INVALID/link", bytes.NewBuffer(headerRequest))
	router.ServeHTTP(r, req1)

	require.JSONEq(t, string(`{"message":"could not find readme"}`), r.Body.String())
}

func TestAddLinkReturnsIncorrectRequestBody(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	r := httptest.NewRecorder()

	var headerRequest = []byte(`{
	}`)

	req1, _ := http.NewRequest("POST", "/readme?name=13", nil)
	router.ServeHTTP(w, req1)

	req2, _ := http.NewRequest("PUT", "/readme/13/link", bytes.NewBuffer(headerRequest))
	router.ServeHTTP(r, req2)

	require.JSONEq(t, string(`{"message":"incorrect request body, should be AddLinkRequest body"}`), r.Body.String())
}

func TestAddImage(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	r := httptest.NewRecorder()

	var headerRequest = []byte(`{
		"link": "https://imgur.com/grhk1rU",
		"description": "Metamask"
	}`)

	req1, _ := http.NewRequest("POST", "/readme?name=14", nil)
	router.ServeHTTP(w, req1)

	req2, _ := http.NewRequest("PUT", "/readme/14/link", bytes.NewBuffer(headerRequest))
	router.ServeHTTP(r, req2)

	require.JSONEq(t, string(`{ "message":"[Metamask](https://imgur.com/grhk1rU)\n" }`), r.Body.String())
}

func TestAddImageReturnsReadmeNotFound(t *testing.T) {
	router := setupRouter()
	r := httptest.NewRecorder()

	var headerRequest = []byte(`{
		"link": "https://imgur.com/grhk1rU",
		"description": "Metamask"
	}`)

	req2, _ := http.NewRequest("PUT", "/readme/15/link", bytes.NewBuffer(headerRequest))
	router.ServeHTTP(r, req2)

	require.JSONEq(t, string(`{ "message":"could not find readme" }`), r.Body.String())
}

func TestAddImageReturnsIncorrectRequestBody(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	r := httptest.NewRecorder()

	var headerRequest = []byte(`{
	}`)

	req1, _ := http.NewRequest("POST", "/readme?name=16", nil)
	router.ServeHTTP(w, req1)

	req2, _ := http.NewRequest("PUT", "/readme/16/link", bytes.NewBuffer(headerRequest))
	router.ServeHTTP(r, req2)

	require.JSONEq(t, string(`{ "message":"incorrect request body, should be AddLinkRequest body" }`), r.Body.String())
}

func TestAddTable(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	r := httptest.NewRecorder()

	var headerRequest = []byte(`{
		"column_names": ["c1", "c2"],
		"column_values": {
			"c1": ["value1", "value2"],
			"c2": ["value3", "value4"]
		}
	}`)

	req1, _ := http.NewRequest("POST", "/readme?name=17", nil)
	router.ServeHTTP(w, req1)

	req2, _ := http.NewRequest("PUT", "/readme/17/table", bytes.NewBuffer(headerRequest))
	router.ServeHTTP(r, req2)

	require.JSONEq(t, string(`{"message":"|c1|c2|\n| --- | --- |\n|value1|value3|\n|value2|value4|\n"}`), r.Body.String())
}

func TestAddTableReturnsIncorrectRequestBody(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	r := httptest.NewRecorder()

	var headerRequest = []byte(`{
	}`)

	req1, _ := http.NewRequest("POST", "/readme?name=18", nil)
	router.ServeHTTP(w, req1)

	req2, _ := http.NewRequest("PUT", "/readme/18/table", bytes.NewBuffer(headerRequest))
	router.ServeHTTP(r, req2)

	require.JSONEq(t, string(`{"message":"incorrect request body, should be AddTableRequest body"}`), r.Body.String())
}

func TestAddTableReturnsReadmeNotFound(t *testing.T) {
	router := setupRouter()
	r := httptest.NewRecorder()

	var headerRequest = []byte(`{
		"column_names": ["c1", "c2"],
		"column_values": {
			"c1": ["value1", "value2"],
			"c2": ["value3", "value4"]
		}
	}`)

	req2, _ := http.NewRequest("PUT", "/readme/19/table", bytes.NewBuffer(headerRequest))
	router.ServeHTTP(r, req2)

	require.JSONEq(t, string(`{"message":"could not find readme"}`), r.Body.String())
}

func TestAddParagraph(t *testing.T) {
	router := setupRouter()
	r := httptest.NewRecorder()
	w := httptest.NewRecorder()

	req1, _ := http.NewRequest("POST", "/readme?name=20", nil)
	router.ServeHTTP(r, req1)

	req2, _ := http.NewRequest("PUT", "/readme/20/paragraph?paragraph=I want this to be a paragraph that will be created for the app", nil)
	router.ServeHTTP(w, req2)

	require.JSONEq(t, string(`{"message":"I want this to be a paragraph that will be created for the app\n"}`), w.Body.String())
}

func TestAddParagraphReturnsReadmeNotFound(t *testing.T) {
	router := setupRouter()
	r := httptest.NewRecorder()

	req1, _ := http.NewRequest("PUT", "/readme/21/paragraph?paragraph=I want this to be a paragraph that will be created for the app", nil)
	router.ServeHTTP(r, req1)

	require.JSONEq(t, string(`{"message":"could not find readme"}`), r.Body.String())
}

func TestAddParagraphReturnsEmptyParagraph(t *testing.T) {
	router := setupRouter()
	r := httptest.NewRecorder()
	w := httptest.NewRecorder()

	req1, _ := http.NewRequest("POST", "/readme?name=22", nil)
	router.ServeHTTP(r, req1)

	req2, _ := http.NewRequest("PUT", "/readme/22/paragraph?paragraph=  ", nil)
	router.ServeHTTP(w, req2)

	require.JSONEq(t, string(`{"message": "paragraph cannot be empty"}`), w.Body.String())
}

func TestHell(t *testing.T) {
	temp := `Hello\nWorld`

	fmt.Println(temp)

	temp2, err := strconv.Unquote(`"` + temp + `"`)
	if err != nil {
		panic(err)
	}
	fmt.Println(temp2)
}
