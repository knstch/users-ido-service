package public

import (
	"fmt"

	"github.com/knstch/knstch-libs/svcerrs"
)

var (
	ErrAccessDenied = fmt.Errorf("access denied: %w", svcerrs.ErrForbidden)
)
