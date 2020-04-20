package parent

/*
  Parent implementations which use a local checkout to create changes.
*/

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"go.skia.org/infra/autoroll/go/config_vars"
	"go.skia.org/infra/autoroll/go/repo_manager/common/git_common"
	"go.skia.org/infra/autoroll/go/repo_manager/common/gitiles_common"
	"go.skia.org/infra/autoroll/go/revision"
	"go.skia.org/infra/go/git"
	"go.skia.org/infra/go/skerr"
)

const (
	// rollBranch is the git branch which is used to create rolls.
	rollBranch = "roll_branch"
)

// GitCheckoutConfig provides configuration for a Parent which uses a local
// checkout to create changes.
type GitCheckoutConfig struct {
	BaseConfig
	git_common.GitCheckoutConfig
}

// See documentation for util.Validator interface.
func (c GitCheckoutConfig) Validate() error {
	if err := c.BaseConfig.Validate(); err != nil {
		return skerr.Wrap(err)
	}
	if err := c.GitCheckoutConfig.Validate(); err != nil {
		return skerr.Wrap(err)
	}
	return nil
}

// GitCheckoutParent is a base for implementations of Parent which use a local
// Git checkout.
type GitCheckoutParent struct {
	*baseParent
	*git_common.Checkout
	createRoll     GitCheckoutCreateRollFunc
	getLastRollRev GitCheckoutGetLastRollRevFunc
	uploadRoll     GitCheckoutUploadRollFunc
}

// GitCheckoutCreateRollFunc generates commit(s) in the local Git checkout to
// be used in the next roll and returns the hash of the commit to be uploaded.
// GitCheckoutParent handles creation of the roll branch.
type GitCheckoutCreateRollFunc func(context.Context, *git.Checkout, *revision.Revision, *revision.Revision, []*revision.Revision, string) (string, error)

// GitCheckoutUploadRollFunc uploads a CL using the given commit hash and
// returns its ID.
type GitCheckoutUploadRollFunc func(context.Context, *git.Checkout, string, string, []string, bool) (int64, error)

// GitCheckoutGetLastRollRevFunc retrieves the last-rolled revision ID from the
// local Git checkout. GitCheckoutParent handles updating the checkout itself.
type GitCheckoutGetLastRollRevFunc func(context.Context, *git.Checkout) (string, error)

// NewGitCheckout returns a base for implementations of Parent which use
// a local checkout to create changes.
func NewGitCheckout(ctx context.Context, c GitCheckoutConfig, reg *config_vars.Registry, client *http.Client, serverURL, workdir string, co *git.Checkout, getLastRollRev GitCheckoutGetLastRollRevFunc, createRoll GitCheckoutCreateRollFunc, uploadRoll GitCheckoutUploadRollFunc) (*GitCheckoutParent, error) {
	// Create a baseParent.
	base, err := newBaseParent(ctx, c.BaseConfig, serverURL)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	// Create the local checkout.
	checkout, err := git_common.NewCheckout(ctx, c.GitCheckoutConfig, reg, workdir, co)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	return &GitCheckoutParent{
		baseParent:     base,
		Checkout:       checkout,
		createRoll:     createRoll,
		getLastRollRev: getLastRollRev,
		uploadRoll:     uploadRoll,
	}, nil
}

// See documentation for Parent interface.
func (p *GitCheckoutParent) Update(ctx context.Context) (string, error) {
	if _, _, err := p.Checkout.Update(ctx); err != nil {
		return "", skerr.Wrap(err)
	}
	return p.getLastRollRev(ctx, p.Checkout.Checkout)
}

// See documentation for Parent interface.
func (p *GitCheckoutParent) CreateNewRoll(ctx context.Context, from, to *revision.Revision, rolling []*revision.Revision, emails []string, cqExtraTrybots string, dryRun bool) (int64, error) {
	// Create the roll branch.
	_, upstreamBranch, err := p.Checkout.Update(ctx)
	if err != nil {
		return 0, skerr.Wrap(err)
	}
	_, _ = p.Git(ctx, "branch", "-D", rollBranch) // Fails if the branch does not exist.
	if _, err := p.Git(ctx, "checkout", "-b", rollBranch, "-t", fmt.Sprintf("origin/%s", upstreamBranch)); err != nil {
		return 0, skerr.Wrap(err)
	}
	if _, err := p.Git(ctx, "reset", "--hard", upstreamBranch); err != nil {
		return 0, skerr.Wrap(err)
	}

	// Generate the commit message.
	// TODO(borenet): This should probably move into parentChildRepoManager.
	commitMsg, err := p.buildCommitMsg(from, to, rolling, emails, cqExtraTrybots, nil)
	if err != nil {
		return 0, skerr.Wrap(err)
	}

	// Run the provided function to create the changes for the roll.
	hash, err := p.createRoll(ctx, p.Checkout.Checkout, from, to, rolling, commitMsg)
	if err != nil {
		return 0, skerr.Wrap(err)
	}

	// Ensure that createRoll generated at least one commit downstream of
	// p.baseCommit, and that it did not leave uncommitted changes.
	commits, err := p.RevList(ctx, "--ancestry-path", "--first-parent", fmt.Sprintf("%s..%s", upstreamBranch, hash))
	if err != nil {
		return 0, skerr.Wrap(err)
	}
	if len(commits) == 0 {
		return 0, skerr.Fmt("createRoll generated no commits!")
	}
	if _, err := p.Git(ctx, "diff", "--quiet"); err != nil {
		return 0, skerr.Wrapf(err, "createRoll left uncommitted changes")
	}
	out, err := p.Git(ctx, "ls-files", "--others", "--exclude-standard")
	if err != nil {
		return 0, skerr.Wrap(err)
	}
	if len(strings.Fields(out)) > 0 {
		return 0, skerr.Fmt("createRoll left untracked files:\n%s", out)
	}

	// Upload the CL.
	return p.uploadRoll(ctx, p.Checkout.Checkout, upstreamBranch, hash, emails, dryRun)
}

// VersionFileGetLastRollRevFunc returns a GitCheckoutGetLastRollRevFunc which
// reads the given file path from the repo and returns its full contents as the
// last-rolled revision ID.
func VersionFileGetLastRollRevFunc(path, dep string) GitCheckoutGetLastRollRevFunc {
	return func(ctx context.Context, co *git.Checkout) (string, error) {
		contents, err := ioutil.ReadFile(filepath.Join(co.Dir(), path))
		if err != nil {
			return "", skerr.Wrap(err)
		}
		// TODO(borenet): It's kind of weird that we're using a gitiles
		// helper function here.
		return gitiles_common.GetPinnedRev(path, dep, string(contents))
	}
}

var _ Parent = &GitCheckoutParent{}