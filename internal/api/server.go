package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/jonesrussell/page-prowler/internal/tasks"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

const (
	DefaultQueueName = "default"
	Protocol         = "https://"
)

type ApiServerInterface struct {
	Inspector *asynq.Inspector
}

func (msi *ApiServerInterface) GetMatchlinks(ctx echo.Context) error {
	queue := DefaultQueueName
	activeTasks, err := msi.Inspector.ListActiveTasks(queue)
	if err != nil {
		return err
	}
	pendingTasks, err := msi.Inspector.ListPendingTasks(queue)
	if err != nil {
		return err
	}

	// Merge activeTasks and pendingTasks
	tasks := append(activeTasks, pendingTasks...)

	// Convert the tasks to JSON and write it to the response
	tasksJson, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	_, err = ctx.Response().Write(tasksJson)
	return err
}

// PostMatchlinks starts the article posting process.
func (msi *ApiServerInterface) PostMatchlinks(ctx echo.Context) error {
	var req PostMatchlinksJSONBody
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	// Validate the input parameters
	if req.URL == nil || *req.URL == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "URL cannot be empty"})
	}
	if req.SearchTerms == nil || *req.SearchTerms == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "SearchTerms cannot be empty"})
	}
	if req.CrawlSiteID == nil || *req.CrawlSiteID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "CrawlSiteID cannot be empty"})
	}
	if req.MaxDepth == nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "MaxDepth cannot be null"})
	}

	// Default Debug to false if it is nil
	if req.Debug == nil {
		req.Debug = new(bool)
		*req.Debug = false
	}

	// Ensure the URL is correctly formatted
	url := strings.TrimSpace(*req.URL)
	if !strings.HasPrefix(url, Protocol) {
		url = Protocol + url
	}

	// Create a new asynq.Client using the same Redis connection details
	redisHost := viper.GetString("REDIS_HOST")
	redisPort := viper.GetString("REDIS_PORT")
	redisAuth := viper.GetString("REDIS_AUTH")
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisAuth,
	})

	payload := &tasks.CrawlTaskPayload{
		URL:         url,
		SearchTerms: *req.SearchTerms,
		CrawlSiteID: *req.CrawlSiteID,
		MaxDepth:    *req.MaxDepth,
		Debug:       *req.Debug,
	}

	tid, err := tasks.EnqueueCrawlTask(client, payload)
	if err != nil {
		log.Println("Error enqueuing crawl task: ", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Crawling started successfully", "task_id": tid})
}

func (msi *ApiServerInterface) GetMatchlinksId(ctx echo.Context, id string) error {
	queue := DefaultQueueName
	info, err := msi.Inspector.GetTaskInfo(queue, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Queue or Task not found"})
		}
		return err
	}

	// Convert the task info to JSON and write it to the response
	taskJson, err := json.Marshal(info)
	if err != nil {
		return err
	}

	_, err = ctx.Response().Write(taskJson)
	return err
}

func (msi *ApiServerInterface) DeleteMatchlinksId(ctx echo.Context, id string) error {
	queue := DefaultQueueName
	_, err := msi.Inspector.GetTaskInfo(queue, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Task not found"})
		}
		return err
	}
	err = msi.Inspector.DeleteTask(queue, id)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Task deleted successfully"})
}

// GetPing handles the ping request.
func (msi *ApiServerInterface) GetPing(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Pong"})
}
