package sys

import (
	"github.com/NatPro/medienstunde/internal/module/lockmod"
	"github.com/NatPro/medienstunde/internal/module/usermod"
	"github.com/NatPro/medienstunde/internal/public/errmsg"
	"github.com/NatPro/medienstunde/internal/sys/systruct"
	"github.com/NatPro/medienstunde/internal/toolbox/net"
)

// Context is the context of the given request and contains everything that is need inside the request.
// all plugins are stored here to mock them later in an easier way
type Context struct {
	systruct.Data

	// plug-ins to hide third party + database calls
	Net     net.Connection
	Usermod usermod.Usermod
	Lockmod lockmod.Lockmod

	// session specific
	User usermod.User
}

// NewContext creates a new 'Context'
func NewContext(s *serverSetup, net net.Connection, log systruct.Logger) (*Context, error) {
	defaultData := systruct.Data{
		Config: s.Config,
		Log:    log,
		ErrMsg: errmsg.LangMap[net.Local],
	}

	return &Context{
		Data: defaultData,

		Net:     net,
		Usermod: usermod.New(s.UserGateCreator.Create(), defaultData),
		Lockmod: lockmod.New(s.LogGateCreator.Create(), defaultData),
	}, nil
}

func (c *Context) Close(commit bool) (err error) {
	var maybeErr error

	maybeErr = c.Lockmod.Close(commit)
	if maybeErr != nil {
		err = maybeErr
	}

	maybeErr = c.Usermod.Close(commit)
	if maybeErr != nil {
		err = maybeErr
	}

	return err
}
