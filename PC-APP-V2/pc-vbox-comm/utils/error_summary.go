package utils

import "bytes"

/* ================================================= Operation Error ================================================ */
type OpError struct {
    TransitionError         error
    EventError              error
}

func (oe *OpError) Error() string {
    var errStr bytes.Buffer

    if oe.TransitionError != nil {
        errStr.WriteString(oe.TransitionError.Error())
    }

    if oe.EventError != nil {
        errStr.WriteString(oe.EventError.Error())
    }
    return errStr.String()
}

func SummarizeErrors(transErr error, eventErr error) error {
    if transErr == nil && eventErr == nil {
        return nil
    }
    return &OpError{TransitionError: transErr, EventError: eventErr}
}

