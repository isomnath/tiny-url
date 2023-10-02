package errs

import (
	"errors"
)

type RepositoryErrors error

var ErrURLInsertFailure RepositoryErrors = errors.New("failed to insert url mapping")
var ErrURLFetchFailure RepositoryErrors = errors.New("failed to fetch url mapping")
var ErrURLFetchFailureKeyDoesNotExist RepositoryErrors = errors.New("failed to fetch url mapping as shortened url does not exist in the set")
var ErrHashInsertFailure RepositoryErrors = errors.New("failed to insert hash map")
var ErrHashFetchFailure RepositoryErrors = errors.New("failed to fetch hash map")
var ErrHashFetchFailureKeyDoesNotExist RepositoryErrors = errors.New("failed to fetch hash mapping as hash does not exist in the set")
var ErrMetricsDTIncrFailure RepositoryErrors = errors.New("failed to increment domain transformation counter")
var ErrMetricsDTFetchFailure RepositoryErrors = errors.New("failed to fetch list of most transformed domains")
var ErrMetricsDRIncrFailure RepositoryErrors = errors.New("failed to increment domain traffic counter")
var ErrMetricsDRFetchFailure RepositoryErrors = errors.New("failed to fetch list of domains with most traffic")
var ErrSeriesCountInitFailure RepositoryErrors = errors.New("failed to set atomic counter with initial value")
var ErrSeriesCounterFailure RepositoryErrors = errors.New("atomic counter could not be incremented")
