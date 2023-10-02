package errs

import (
	"fmt"
)

type ServiceErrors error

var ErrURLDoesNotExist ServiceErrors = fmt.Errorf("shortened url does not exist")
var ErrURLGenerationHashCollisionURL ServiceErrors = fmt.Errorf("hash collision has occurred: same hash generated for different url, please retry")
var ErrFailedToFetchURLFromURLStore ServiceErrors = fmt.Errorf("failed to fetch shortened url from url store")
var ErrFailedToFetchHashFromHashStore ServiceErrors = fmt.Errorf("failed to fetch url hash from hash store")
var ErrFailedToFetchNextValueInSeries ServiceErrors = fmt.Errorf("failed to fetch next series value from series store")
var ErrFailedToSaveHashInStore ServiceErrors = fmt.Errorf("failed to save hash in hash store")
var ErrFailedToSaveURLInStore ServiceErrors = fmt.Errorf("failed to save shortened URL in URL store")
