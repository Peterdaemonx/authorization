//go:build integration
// +build integration

package spanner

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"testing"

	gspanner "cloud.google.com/go/spanner"
	"gitlab.cmpayments.local/creditcard/authorization/database/fixture"
)

const (
	spanner_test_database      = "authorizations"
	spanner_test_instance      = "authorization-instance"
	spanner_test_project       = "cc-acquiring-development"
	spanner_test_configuration = "cc-acquiring-test-configuration"
	grpcPort                   = 65000
	restPort                   = 65100
)

var (
	db = fmt.Sprintf("projects/%s/instances/%s/databases/%s", spanner_test_project, spanner_test_instance, spanner_test_database)
	//nolint:deadcode,varcheck,unused
	migrateDb = fmt.Sprintf("spanner://projects/%s/instances/%s/databases/%s?x-clean-statements=true", spanner_test_project, spanner_test_instance, spanner_test_database)
)

//nolint:deadcode,unused
func setup() {
	fmt.Println("Create Spanner emulator config...")
	os.Setenv("SPANNER_EMULATOR_HOST", fmt.Sprintf("localhost:%v", grpcPort))

	createConfigCmd := exec.Command("gcloud", "config", "configurations", "create", spanner_test_configuration)
	output, err := createConfigCmd.CombinedOutput()
	if err != nil {
		cmdErr := exec.Command("gcloud", "config", "configurations", "activate", spanner_test_configuration).Start()
		if cmdErr != nil {
			fmt.Println("Error doing config stuff: " + cmdErr.Error() + " - " + string(output))
			os.Exit(1)
		}
	}
	// Print the output
	fmt.Println(string(output))
	fmt.Println("Disable auth credentials...")
	disableAuthCredsCmd := exec.Command("gcloud", "config", "set", "auth/disable_credentials", "true")
	output, err = disableAuthCredsCmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error disabling credentials: " + err.Error() + " - " + string(output))
		os.Exit(1)
	}
	// Print the output
	fmt.Println(string(output))
	fmt.Println("Set gcloud spanner project...")
	setProjectTestEmulatorTestCmd := exec.Command("gcloud", "config", "set", "project", spanner_test_project)
	output, err = setProjectTestEmulatorTestCmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error setting project: " + err.Error() + " - " + string(output))
		os.Exit(1)
	}
	// Print the output
	fmt.Println(string(output))
	fmt.Println("Set api endpoint...")
	setApiEndpointOverrideCmd := exec.Command("gcloud", "config", "set", "api_endpoint_overrides/spanner", fmt.Sprintf("http://localhost:%v/", restPort))
	output, err = setApiEndpointOverrideCmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error setting API endpoint: " + err.Error() + " - " + string(output))
		os.Exit(1)
	}
	//Print the output
	fmt.Println(string(output))
	fmt.Println("Start Spanner emulator...")
	var startGcloudEmulator *exec.Cmd
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		startGcloudEmulator = exec.Command("docker", "run", "-d", "--name", "test-container", "-p", fmt.Sprintf("%v:9010", grpcPort), "-p", fmt.Sprintf("%v:9020", restPort), "gcr.io/cloud-spanner-emulator/emulator")
	} else {
		startGcloudEmulator = exec.Command("gcloud", "emulators", "spanner", "start", fmt.Sprintf("--host-port=localhost:%v", grpcPort), fmt.Sprintf("--rest-port=%v", restPort))
	}
	err = startGcloudEmulator.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// Print the output
	fmt.Println(string(output))
	fmt.Println("creating spanner instance...")
	createSpannerInstance := exec.Command("gcloud", "spanner", "instances", "create", spanner_test_instance, "--config=emulator-config", `--description="Creditcard acquiring"`, "--nodes=1")
	output, err = createSpannerInstance.CombinedOutput()
	if err != nil {
		fmt.Println("Error creating instance: " + err.Error() + " - " + string(output))
		os.Exit(1)
	}
	// Print the output
	fmt.Println(string(output))

	fmt.Println("creating spanner database...")
	createDatabase := exec.Command("gcloud", "spanner", "databases", "create", spanner_test_database, "--instance="+spanner_test_instance)
	output, err = createDatabase.CombinedOutput()
	if err != nil {
		fmt.Println("Error creating database: " + err.Error() + " - " + string(output))
		return
	}
	// Print the output
	fmt.Println(string(output))

	fmt.Println("doing migrations...")
	migrationCmd := exec.Command("migrate", "-path", "../../../database/migrations/", "-database", migrateDb, "up")
	output, err = migrationCmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running migrations %s: "+err.Error()+" - "+string(output), migrationCmd.String())
		os.Exit(1)
	}

	// Print the output
	fmt.Println(string(output))
}

//nolint:deadcode,varcheck,unused
func teardown() {
	fmt.Println("Teardown spanner test emulation...")
	createConfigCmd := exec.Command("gcloud", "config", "configurations", "create", "default")
	stdout, err := createConfigCmd.Output()
	if err != nil {
		//nolint:deadcode,varcheck,unused,errcheck
		err := exec.Command("gcloud", "config", "configurations", "activate", "default").Run()
		if err != nil {
			return
		}
	}
	// Print the output
	fmt.Println(string(stdout))

	deleteTestEmulator := exec.Command("gcloud", "config", "configurations", "delete", spanner_test_configuration)
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		stopContainer := exec.Command("docker", "stop", "test-container")
		output, err := stopContainer.Output()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(output))

		removeContainer := exec.Command("docker", "rm", "test-container")
		output, err = removeContainer.Output()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(output))
	}
	//nolint:deadcode,varcheck,unused,ineffassign,staticcheck
	stdout, err = deleteTestEmulator.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

var Client *gspanner.Client

func TestMain(m *testing.M) {
	// setup the gcloud emulator.
	setup()
	// set the env variable SPANNER_EMULATOR_HOST to indicate the spanner needs to use the emulator.
	os.Setenv("SPANNER_EMULATOR_HOST", fmt.Sprintf("localhost:%v", grpcPort))
	var err error
	// Connect to the spanner client.
	Client, err = NewSpannerClient(context.Background(), db, 4)
	if err != nil {
		log.Fatalln(err.Error())
	}
	// Close the client at the end of the statement.
	defer Client.Close()

	fmt.Println("Seeding data...")
	err = fixture.Seed(context.Background(), Client)
	if err != nil {
		log.Fatalf("error seeding database: %s", err)
	}

	// run all the tests.
	retCode := m.Run()

	// removes the gcloud configuration or removes the docker containers when on Windows or MacOS
	teardown()
	os.Exit(retCode)
}
