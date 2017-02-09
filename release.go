package wormhole

import (
	"fmt"
	"os"
	"strings"

	"github.com/superfly/wormhole/messages"

	git "srcd.works/go-git.v4"
	"srcd.works/go-git.v4/plumbing"
)

func computeRelease(releaseIDVar, releaseDescVar string) (*messages.Release, error) {
	release := &messages.Release{}
	if releaseID := os.Getenv(releaseIDVar); releaseID != "" {
		release.ID = releaseID
	}

	if releaseDesc := os.Getenv(releaseDescVar); releaseDesc != "" {
		release.Description = releaseDesc
	}

	var branches []string
	if _, err := os.Stat(".git"); !os.IsNotExist(err) {
		release.VCSType = "git"
		repo, err := git.PlainOpen(".")
		if err != nil {
			return nil, fmt.Errorf("Could not open repository: %s", err.Error())
		}
		head, err := repo.Head()
		if err != nil {
			return nil, fmt.Errorf("Could not get repo head: %s", err.Error())
		}

		oid := head.Hash()
		release.VCSRevision = oid.String()
		tip, err := repo.Commit(oid)
		if err != nil {
			return nil, fmt.Errorf("Could not get current commit: %s", err.Error())
		}

		refs, err := repo.References()
		if err != nil {
			return nil, fmt.Errorf("Could not get current refs: %s", err.Error())
		}
		refs.ForEach(func(ref *plumbing.Reference) error {
			if ref.IsBranch() && head.Hash().String() == ref.Hash().String() {
				branch := strings.TrimPrefix(ref.Name().String(), "refs/heads/")
				branches = append(branches, branch)
			}
			return nil
		})

		author := tip.Author
		release.VCSRevisionAuthorEmail = author.Email
		release.VCSRevisionAuthorName = author.Name
		release.VCSRevisionTime = author.When
		release.VCSRevisionMessage = tip.Message
	}
	if release.ID == "" && release.VCSRevision != "" {
		release.ID = release.VCSRevision
	}
	if release.Description == "" && release.VCSRevisionMessage != "" {
		release.Description = release.VCSRevisionMessage
	}
	// TODO: be smarter about branches, and maybe let users override this
	if len(branches) > 0 {
		release.Branch = branches[0]
	}
	return release, nil
}