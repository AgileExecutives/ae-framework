package testutils

import (
	"github.com/AgileExecutives/ae-framework/serverbase/modules/email/repo"
)

// NewInMemoryEmailRepo returns an EmailRepo suitable for unit tests.
func NewInMemoryEmailRepo() repo.EmailRepo {
	return repo.NewInMemoryEmailRepo()
}
