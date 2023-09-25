package msgbrokerctrl // import "github.com/erupshis/effective_mobile/internal/server/controllers/msgbrokerctrl"

type Controller struct{ ... }
    func Create(chIn <-chan msgbroker.Message, chError chan<- msgbroker.Message, ...) *Controller
