package local

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"

	"github.com/marmotedu/errors"
)

var errCtxGetFatherCtxFromMainByType = context.ErrContext(errors.New("cannot query father context from `Main` context by type"))
