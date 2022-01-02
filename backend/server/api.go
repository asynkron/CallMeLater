package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	HeaderPrefix         = "X-Later-"
	HeaderRequestUrl     = HeaderPrefix + "Request-Url"
	HeaderWhen           = HeaderPrefix + "When"
	HeaderResponseUrl    = HeaderPrefix + "Response-Url"
	HeaderResponseMethod = HeaderPrefix + "Response-Method"
)

type apiServer struct {
	worker *worker
}

func (a *apiServer) createJob(c *gin.Context) {
	requestUrl, err := url.Parse(c.GetHeader(HeaderRequestUrl))
	if err != nil {
		log.Err(err).Msg("JobStatusFailed to parse request url")
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	when, err := time.ParseDuration(c.GetHeader(HeaderWhen))
	if err != nil {
		log.Err(err).Msg("failed to parse when")
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	tmp := c.GetHeader(HeaderResponseUrl)
	var responseUrlStr string
	if tmp != "" {
		responseUrl, err := url.Parse(tmp)
		if err != nil {
			log.Err(err).Msg("failed to parse response url")
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		responseUrlStr = responseUrl.String()
	}

	body, err := ioutil.ReadAll(c.Request.Body)

	t := time.Now().Add(when)

	var job = &HttpRequestJob{
		Id:                 uuid.New().String(),
		RequestUrl:         requestUrl.String(),
		RequestMethod:      c.Request.Method,
		ResponseUrl:        responseUrlStr,
		ResponseMethod:     c.GetHeader(HeaderResponseMethod),
		ScheduledTimestamp: t,
		Header:             c.Request.Header,
		Form:               c.Request.Form,
		Body:               body,
	}
	job.InitDefaults()

	a.saveRequest(job)

	c.SetAccepted()

	log.Info().Msg("Request accepted")
}

func (a *apiServer) saveRequest(rd *HttpRequestJob) {
	err := a.worker.storage.Create(rd)
	if err != nil {
		log.Err(err).Msg("Error saving request")
		return
	}
	log.Info().Msg("Saved Request")

	a.worker.executableJobs <- rd
}

func (a *apiServer) read(w http.ResponseWriter, r *http.Request) {
	_, err := a.worker.storage.Read(0, 0)
	if err != nil {
		log.Err(err).Msg("Error reading requests")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error reading requests")
		return
	}
}

func handleRequests(worker *worker) {
	a := &apiServer{
		worker: worker,
	}

	r := gin.Default()
	r.Any("/later", a.createJob)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
