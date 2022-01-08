package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	HeaderPrefix         = "X-Later-"
	HeaderRequestUrl     = HeaderPrefix + "Request-Url"
	HeaderWhen           = HeaderPrefix + "When"
	HeaderCron           = HeaderPrefix + "Cron"
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
		RequestUrl:     requestUrl.String(),
		RequestMethod:  c.Request.Method,
		ResponseUrl:    responseUrlStr,
		ResponseMethod: c.GetHeader(HeaderResponseMethod),
		Header:         c.Request.Header,
		Form:           c.Request.Form,
		Body:           body,
		JobEntity:      newHttpRequestJobEntity(t, ""),
	}
	job.CronExpression = c.GetHeader(HeaderCron)
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

	skip, _ := strconv.Atoi(c.Query("limit"))
	limit, _ := strconv.Atoi(c.Param("limit"))
	jobs, err := a.worker.storage.Read(skip, limit)
	if err != nil {
		log.Err(err).Msg("Error reading requests")
		c.String(http.StatusInternalServerError, "Error reading jobs")
		return
	}
	response := GetJobsResponse{
		Skip:  skip,
		Limit: limit,
	}
	response.Jobs = make([]JobResponse, 0)

	for _, job := range jobs {
		jobResponse := JobResponse{
			Id:                 job.Id,
			ScheduledTimestamp: job.ScheduledTimestamp,
			ExecutedTimestamp:  job.ExecutedTimestamp,
			ExecutedStatus:     job.ExecutedStatus,
			Description:        job.Description,
			CronExpression:     job.CronExpression,
			DataDiscriminator:  job.DataDiscriminator,
			ParentJobId:        job.ParentJobId,
			Status:             job.Status,
			RetryCount:         job.RetryCount,
		}
		response.Jobs = append(response.Jobs, jobResponse)
	}
	c.JSON(http.StatusOK, response)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func handleRequests(worker *worker) {
	a := &apiServer{
		worker: worker,
	}

	r := gin.Default()
	// CORS for https://foo.com and https://github.com origins, allowing:
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	r.Use(CORSMiddleware())

	r.Any("/later", a.createJob)
	r.GET("/jobs/:skip/:limit", a.read)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
