package state

type opError struct {
    TransitionError         error
    EventError              error
}

func (oe *opError) Error() string {
    var errStr string = ""

    if oe.TransitionError != nil {
        errStr += oe.TransitionError.Error()
    }

    if oe.EventError != nil {
        errStr += oe.EventError.Error()
    }
    return errStr
}

func summarizeErrors(transErr error, eventErr error) *opError {
    if transErr == nil && eventErr == nil {
        return nil
    }
    return &opError{TransitionError: transErr, EventError: eventErr}
}

