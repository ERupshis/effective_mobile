package msgsavectrl // import "github.com/erupshis/effective_mobile/internal/server/controllers/msgsavectrl"

type Controller struct{ ... }
    func Create(chIn <-chan datastructs.ExtraDataFilling, chError chan<- msgbroker.Message, ...) *Controller
