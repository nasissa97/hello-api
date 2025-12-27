package main

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"

	"hello-api/config"
	"hello-api/handlers/rest"

	"github.com/cucumber/godog"
	"github.com/go-redis/redis"
	"github.com/go-resty/resty/v2"
	"github.com/ory/dockertest"
)

type apiFeature struct {
	client   *resty.Client
	server   *httptest.Server
	word     string
	language string
}

var (
	pool     *dockertest.Pool
	database *dockertest.Resource
)

func (api *apiFeature) iTranslateItTo(arg1 string) error {
	api.language = arg1
	return nil
}

func (api *apiFeature) theResponseShouldBe(arg1 string) error {
	url := fmt.Sprintf("%s/%s", api.server.URL, api.word)

	resp, err := api.client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"language": api.language,
		}).
		SetResult(&rest.Resp{}).
		Get(url)
	if err != nil {
		return err
	}

	res := resp.Result().(*rest.Resp)
	if res.Translation != arg1 {
		return fmt.Errorf("translation should be set to %s", arg1)
	}
	return nil
}

func (api *apiFeature) theWord(arg1 string) error {
	api.word = arg1
	return nil
}

func InitializeTestSuite(sc *godog.TestSuiteContext) {
	var err error

	sc.BeforeSuite(func() {
		pool, err = dockertest.NewPool("")
		if err != nil {
			panic(fmt.Sprintf("unable to create connection pool %v", err))
		}
		wd, err := os.Getwd()
		if err != nil {
			panic(fmt.Sprintf("unable to get working directory %v", err))
		}

		mount := fmt.Sprintf("%s/data/:/data/", filepath.Dir(wd))
		fmt.Println(mount)
		resource, err := pool.RunWithOptions(&dockertest.RunOptions{
			Repository: "redis",
			Mounts:     []string{mount},
		})
		addr := fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp"))
		if err != nil {
			panic(fmt.Sprintf("unable to create container: %s", err))
		}
		if err := pool.Retry(func() error {
			client := redis.NewClient(&redis.Options{Addr: addr})

			if err := client.Ping().Err(); err != nil {
				return err
			}

			err = client.Set("hello:german", "hallo", 0).Err()
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			panic(fmt.Sprintf("failed to seed redis: %v", err))
		}
		if err := resource.Expire(600); err != nil {
			panic("unable to set expiration on container")
		}
		database = resource
	})

	sc.AfterSuite(func() {
		_ = database.Close()
	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	client := resty.New()
	api := &apiFeature{
		client: client,
	}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		cfg := config.Configuration{}
		cfg.LoadFromEnv()
		cfg.DatabasePort = database.GetPort("6379/tcp")
		cfg.DatabaseURL = "localhost"
		fmt.Printf("%+v\n", cfg)
		mux := API(cfg)
		server := httptest.NewServer(mux)
		api.server = server
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		api.server.Close()
		return ctx, nil
	})

	ctx.Step(`^I translate it to "([^"]*)"$`, api.iTranslateItTo)
	ctx.Step(`^the response should be "([^"]*)"$`, api.theResponseShouldBe)
	ctx.Step(`^the word "([^"]*)"$`, api.theWord)
}
