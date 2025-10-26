package bdd

import (
	"context"
	"os"
	"strconv"
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

	//create
	sc.Step(`^the API is running reacheable via http$`, theAPIIsRunningReacheableViaHttp)
	sc.Step(`^I POST "([^"]*)" with json:$`, w.iPOSTWithJSON)
	sc.Step(`^the response code should be (\d+)$`, w.theResponseCodeShouldBe)
	sc.Step(`^the response json at "([^"]*)" should be "([^"]*)"$`, w.jsonAtShouldBe)
	sc.Step(`^the response json has keys: "([^"]*)", "([^"]*)", "([^"]*)"$`, theResponseJsonHasKeys)

	//fetch
	sc.Step(`^a device exists with name "([^"]*)" and brand "([^"]*)"$`, w.aDeviceExistsWithNameAndBrand)
	sc.Step(`^I GET "([^"]*)"$`, w.iGET)
	sc.Step(`^there are more than (\d+) devices stored$`, func(x string) error {
		n, err := strconv.Atoi(x)
		if err != nil {
			return err
		}
		return w.thereAreMoreThanDevicesStored(n)
	})
	sc.Step(`^the response json should contain (\d+) devices$`, func(x string) error {
		n, err := strconv.Atoi(x)
		if err != nil {
			return err
		}
		return w.theResponseJSONShouldContainNDevices(n)
	})
	sc.Step(`^the response json should include "next_page" and "previous_page" fields$`, w.theResponseJSONShouldIncludeNextPrev)
	sc.Step(`^the API is running$`, theAPIIsRunning)

	sc.Step(`^the response json should contain (\d+) device[s]?$`, w.theResponseJSONShouldContainNDevices)
	sc.Step(`^the response json at "([^"]*)" should be "([^"]*)"$`, w.jsonAtShouldBe)

	sc.Step(`^the response json should contain (\d+) devices$`, w.theResponseJSONShouldContainNDevices)
	sc.Step(`^the response json at "([^"]*)" should be "([^"]*)"$`, w.responseJsonAtShouldBe)
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
