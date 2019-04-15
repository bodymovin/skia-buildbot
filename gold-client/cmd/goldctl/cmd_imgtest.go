package main

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"go.skia.org/infra/gold-client/go/goldclient"
	"go.skia.org/infra/golden/go/jsonio"
	"go.skia.org/infra/golden/go/shared"
)

// imgTestEnv is the environment for the imgtest command ant its sub-commands.
type imgTestEnv struct {
	// Flags used by imgtest:init and imgtest:add.
	flagCommit       string // flag containing the commit hash
	flagKeysFile     string
	flagIssueID      string
	flagPatchsetID   string
	flagJobID        string
	flagInstanceID   string
	flagWorkDir      string
	flagPassFailStep bool
	flagFailureFile  string
	flagURL          string

	// Flags used by imgtest:add
	flagTestName string
	flagPNGFile  string
}

// getImgTestCmd returns the definition of the imgtest command.
func getImgTestCmd() *cobra.Command {
	env := &imgTestEnv{}

	// imgtest command and its sub commands
	imgTestCmd := &cobra.Command{
		Use:   "imgtest",
		Short: "Collect  and upload test results as images",
		Long: `
Collect and upload test results to the Gold backend.`,
	}

	// cmd: imgtest init
	imgTestInitCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a  testing environment",
		Long: `
Start a testing session during which tests are added. This initializes the environment.
It gathers whether the 'add' command returns a pass/fail value and the common
keys shared by all tests that are added via 'add'.`,
		Run: env.runImgTestInitCmd,
	}
	env.addCommonFlags(imgTestInitCmd, false)

	imgTestAddCmd := &cobra.Command{
		Use:   "add",
		Short: "Adds a test image to the results.",
		Long: `
Add images generated by the tests to the test results. This requires two arguments:
			 - The test name
			 - The path to the resulting PNG.
`,
		Run:  env.runImgTestAddCmd,
		Args: cobra.NoArgs,
	}
	env.addCommonFlags(imgTestAddCmd, true)
	imgTestAddCmd.Flags().StringVarP(&env.flagTestName, "test-name", "", "", "Unique name of the test, must not contain spaces.")
	imgTestAddCmd.Flags().StringVarP(&env.flagPNGFile, "png-file", "", "", "Path to the PNG file that contains the test results.")

	Must(imgTestAddCmd.MarkFlagRequired("test-name"))
	Must(imgTestAddCmd.MarkFlagRequired("png-file"))

	imgTestFinalizeCmd := &cobra.Command{
		Use:   "finalize",
		Short: "Finish adding tests and process results.",
		Long: `
All tests have been added. Upload images and generate and upload the JSON file that captures
test results.`,
		Run: env.runImgTestFinalizeCmd,
	}
	imgTestFinalizeCmd.Flags().StringVarP(&env.flagWorkDir, fstrWorkDir, "", "", "Work directory for intermediate results")
	Must(imgTestAddCmd.MarkFlagRequired("work-dir"))

	imgTestPassFailCmd := &cobra.Command{
		Use:   "passfail",
		Short: "Checks whether the results match expectations",
		Long: `
Check against Gold or local baseline whether the results match the expectations`,
		Run: env.runImgTestPassFailCmd,
	}

	// assemble the imgtest command.
	imgTestCmd.AddCommand(
		imgTestInitCmd,
		imgTestAddCmd,
		imgTestFinalizeCmd,
		imgTestPassFailCmd,
	)
	return imgTestCmd
}

func (i *imgTestEnv) addCommonFlags(cmd *cobra.Command, optional bool) {
	cmd.Flags().StringVarP(&i.flagInstanceID, "instance", "", "", "ID of the Gold instance.")
	cmd.Flags().StringVarP(&i.flagWorkDir, fstrWorkDir, "", "", "Work directory for intermediate results")
	cmd.Flags().BoolVarP(&i.flagPassFailStep, "passfail", "", false, "Whether the 'add' call returns a pass/fail for each test.")

	cmd.Flags().StringVarP(&i.flagCommit, "commit", "", "", "Git commit hash")
	cmd.Flags().StringVarP(&i.flagKeysFile, "keys-file", "", "", "JSON file containing key/value pairs commmon to all tests")
	cmd.Flags().StringVarP(&i.flagIssueID, "issue", "", "", "Gerrit issue if this is trybot run. ")
	cmd.Flags().StringVarP(&i.flagPatchsetID, "patchset", "", "", "Gerrit patchset number if this is a trybot run. ")
	cmd.Flags().StringVarP(&i.flagJobID, "jobid", "", "", "Job ID if this is a tryjob run. Current the BuildBucket id.")
	cmd.Flags().StringVarP(&i.flagFailureFile, "failure-file", "", "", "Path to the file where to write failure information")
	cmd.Flags().StringVarP(&i.flagURL, "url", "", "", "URL of the Gold instance. Used for testing, if empty the URL will be derived from the value of 'instance'")

	Must(cmd.MarkFlagRequired(fstrWorkDir))
	if !optional {
		Must(cmd.MarkFlagRequired("instance"))
		Must(cmd.MarkFlagRequired("commit"))
		Must(cmd.MarkFlagRequired("keys-file"))
	}
}

// TODO(kjlubick): Implement these stubbed out sub-commands.
func (i *imgTestEnv) runImgTestPassFailCmd(cmd *cobra.Command, args []string) { notImplemented(cmd) }

func (i *imgTestEnv) runImgTestInitCmd(cmd *cobra.Command, args []string) {
	auth, err := goldclient.LoadAuthOpt(i.flagWorkDir)
	ifErrLogExit(cmd, err)

	if auth == nil {
		logErrf(cmd, "Auth is empty - did you call goldctl auth first?")
		exitProcess(cmd, 1)
	}

	auth.SetDryRun(flagDryRun)

	keyMap, err := readKeysFile(i.flagKeysFile)
	ifErrLogExit(cmd, err)

	validation := shared.Validation{}
	issueID := validation.Int64Value("issue", i.flagIssueID, 0)
	patchsetID := validation.Int64Value("patchset", i.flagPatchsetID, 0)
	jobID := validation.Int64Value("jobid", i.flagJobID, 0)
	ifErrLogExit(cmd, validation.Errors())

	config := goldclient.GoldClientConfig{
		InstanceID:      i.flagInstanceID,
		WorkDir:         i.flagWorkDir,
		PassFailStep:    i.flagPassFailStep,
		OverrideGoldURL: i.flagURL,
	}
	goldClient, err := goldclient.NewCloudClient(auth, config)
	ifErrLogExit(cmd, err)

	// Define the meta data of the result that is shared by all tests.
	gr := jsonio.GoldResults{
		GitHash:       i.flagCommit,
		Key:           keyMap,
		Issue:         issueID,
		Patchset:      patchsetID,
		BuildBucketID: jobID,
	}

	logVerbose(cmd, "Loading hashes and baseline from Gold\n")
	err = goldClient.SetSharedConfig(gr)
	ifErrLogExit(cmd, err)

	logInfof(cmd, "Directory %s successfully loaded with configuration\n", i.flagWorkDir)
}

// runImgTestCommand processes and uploads test results to Gold.
func (i *imgTestEnv) runImgTestAddCmd(cmd *cobra.Command, args []string) {
	auth, err := goldclient.LoadAuthOpt(i.flagWorkDir)
	ifErrLogExit(cmd, err)

	if auth == nil {
		logErrf(cmd, "Auth is empty - did you call goldctl auth first?")
		exitProcess(cmd, 1)
	}

	auth.SetDryRun(flagDryRun)

	var goldClient goldclient.GoldClient

	if i.flagKeysFile != "" {
		// user has specified a full set of keys. This happens if they
		// did not (or could not) call init before the start of their test
		keyMap, err := readKeysFile(i.flagKeysFile)
		ifErrLogExit(cmd, err)

		validation := shared.Validation{}
		issueID := validation.Int64Value("issue", i.flagIssueID, 0)
		patchsetID := validation.Int64Value("patchset", i.flagPatchsetID, 0)
		jobID := validation.Int64Value("jobid", i.flagJobID, 0)
		ifErrLogExit(cmd, validation.Errors())

		// Define the meta data of the result that is shared by all tests.
		gr := jsonio.GoldResults{
			GitHash:       i.flagCommit,
			Key:           keyMap,
			Issue:         issueID,
			Patchset:      patchsetID,
			BuildBucketID: jobID,
		}

		config := goldclient.GoldClientConfig{
			InstanceID:      i.flagInstanceID,
			WorkDir:         i.flagWorkDir,
			PassFailStep:    i.flagPassFailStep,
			OverrideGoldURL: i.flagURL,
		}
		goldClient, err = goldclient.NewCloudClient(auth, config)
		ifErrLogExit(cmd, err)

		err = goldClient.SetSharedConfig(gr)
		ifErrLogExit(cmd, err)
	} else {
		// the user is presumed to have called init first, so we can just load
		// it
		goldClient, err = goldclient.LoadCloudClient(auth, i.flagWorkDir)
		ifErrLogExit(cmd, err)
	}

	pass, err := goldClient.Test(i.flagTestName, i.flagPNGFile)
	ifErrLogExit(cmd, err)

	if !pass {
		logErrf(cmd, "Test: %s FAIL\n", i.flagTestName)
		exitProcess(cmd, 1)
	}
	logInfof(cmd, "Test: %s PASS\n", i.flagTestName)
	exitProcess(cmd, 0)
}

func (i *imgTestEnv) runImgTestFinalizeCmd(cmd *cobra.Command, args []string) {
	auth, err := goldclient.LoadAuthOpt(i.flagWorkDir)
	ifErrLogExit(cmd, err)

	if auth == nil {
		logErrf(cmd, "Auth is empty - did you call goldctl auth first?")
		exitProcess(cmd, 1)
	}

	auth.SetDryRun(flagDryRun)

	// the user is presumed to have called init and tests first, so we just
	// have to load it from disk.
	goldClient, err := goldclient.LoadCloudClient(auth, i.flagWorkDir)
	ifErrLogExit(cmd, err)

	logVerbose(cmd, "Uploading the final JSON to Gold\n")
	err = goldClient.Finalize()
	ifErrLogExit(cmd, err)
	exitProcess(cmd, 0)
}

// readKeysFile is a helper function to read a JSON file with key/value pairs.
func readKeysFile(keysFile string) (map[string]string, error) {
	reader, err := os.Open(keysFile)
	if err != nil {
		return nil, err
	}

	ret := map[string]string{}
	err = json.NewDecoder(reader).Decode(&ret)
	return ret, err
}
