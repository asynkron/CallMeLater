package server

import (
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
		c.String(http.StatusBadRequest, "JobStatusFailed to parse request url")
		return
	}
	when, err := time.ParseDuration(c.GetHeader(HeaderWhen))
	if err != nil {
		log.Err(err).Msg("failed to parse when")
		c.String(http.StatusBadRequest, "failed to parse when")
		return
	}
	tmp := c.GetHeader(HeaderResponseUrl)
	var responseUrlStr string
	if tmp != "" {
		responseUrl, err := url.Parse(tmp)
		if err != nil {
			log.Err(err).Msg("failed to parse response url")
			c.String(http.StatusBadRequest, "failed to parse response url")
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

	c.String(http.StatusAccepted, "Job Created")

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

func (a *apiServer) read(c *gin.Context) {
	_, err := a.worker.storage.Read(0, 0)
	if err != nil {
		log.Err(err).Msg("Error reading requests")
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func handleRequests(worker *worker) {
	a := &apiServer{
		worker: worker,
	}

	r := gin.Default()
	r.Any("/later", a.createJob)
	r.GET("/jobs/:skip/:limit", a.read)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
