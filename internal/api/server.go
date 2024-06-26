package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/consumer"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/tasks"
	"github.com/labstack/echo/v4"
)

const (
	DefaultQueueName = "default"
	Protocol         = "https://"
)

type ServerAPIInterface struct {
	Inspector *asynq.Inspector
}

func (msi *ServerAPIInterface) GetGetlinks(ctx echo.Context, params GetGetlinksParams) error {
	siteid := params.Siteid
	if siteid == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "siteid cannot be empty"})
	}

	manager, ok := ctx.Get(common.CrawlManagerKeyStr).(*crawler.CrawlManager)
	if !ok || manager == nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "CrawlManager not found in context"})
	}

	links, err := consumer.RetrieveAndUnmarshalLinks(ctx.Request().Context(), manager, siteid)
	if err != nil {
		return err
	}

	output := consumer.CreateOutput(siteid, links)

	jsonOutput, err := consumer.MarshalOutput(output)
	if err != nil {
		return err
	}

	return ctx.Blob(http.StatusOK, "application/json", jsonOutput)
}

func (msi *ServerAPIInterface) GetMatchlinks(ctx echo.Context) error {
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
	allTasks := append(activeTasks, pendingTasks...)

	// Convert the tasks to JSON and write it to the response
	tasksJSON, err := json.Marshal(allTasks)
	if err != nil {
		return err
	}

	_, err = ctx.Response().Write(tasksJSON)
	return err
}

// PostMatchlinks starts the article posting process.
func (msi *ServerAPIInterface) PostMatchlinks(ctx echo.Context) error {
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
	manager, ok := ctx.Get(common.CrawlManagerKeyStr).(*crawler.CrawlManager)
	if !ok || manager == nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "CrawlManager not found in context"})
	}
	redisDetails := manager.Client.Options()
	redisAddr := redisDetails.Addr
	redisAuth := redisDetails.Password

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redisAddr,
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

//nolint:revive
func (msi *ServerAPIInterface) GetMatchlinksId(ctx echo.Context, id string) error {
	queue := DefaultQueueName
	info, err := msi.Inspector.GetTaskInfo(queue, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Queue or Task not found"})
		}
		return err
	}

	// Convert the task info to JSON and write it to the response
	taskJSON, err := json.Marshal(info)
	if err != nil {
		return err
	}

	_, err = ctx.Response().Write(taskJSON)
	return err
}

//nolint:revive
func (msi *ServerAPIInterface) DeleteMatchlinksId(ctx echo.Context, id string) error {
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
func (msi *ServerAPIInterface) GetPing(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Pong"})
}
