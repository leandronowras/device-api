package bdd

import (
	"context"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

func InitializeScenario(sc *godog.ScenarioContext) {
	w := &apiWorld{}

	sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return ctx, w.theAPIIsRunning()
	})
	sc.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		w.stopServer()
		return ctx, nil
	})

	sc.Step(`^the API is running reacheable via http$`, theAPIIsRunningReacheableViaHttp)
	sc.Step(`^I POST "([^"]*)" with json:$`, w.iPOSTWithJSON)
	sc.Step(`^the response code should be (\d+)$`, w.theResponseCodeShouldBe)
	sc.Step(`^the response json at "([^"]*)" should be "([^"]*)"$`, w.jsonAtShouldBe)
	sc.Step(`^the response json has keys: "([^"]*)", "([^"]*)", "([^"]*)"$`, theResponseJsonHasKeys)
}

func TestMain(m *testing.M) {
	opts := godog.Options{
		Output: colors.Colored(os.Stdout),
		Format: "pretty",
		Paths:  []string{"features"},
	}
	code := godog.TestSuite{
		Name:                "devices",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()
	os.Exit(code)
}
