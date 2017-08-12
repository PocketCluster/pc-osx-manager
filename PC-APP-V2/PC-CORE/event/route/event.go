package route

// route event
type Event struct {
    method     string
    path       string
    payload    string
}

func RouteEventGet(path string) Event {
    return Event{
        method:     get,
        path:       path,
    }
}

func RouteEventPost(path, payload string) Event {
    return Event{
        method:     post,
        path:       path,
        payload:    payload,
    }
}

func RouteEventPut(path, payload string) Event {
    return Event{
        method:     put,
        path:       path,
        payload:    payload,
    }
}

func RouteEventDelete(path string) Event {
    return Event{
        method:     deleteh,
        path:       path,
    }
}
