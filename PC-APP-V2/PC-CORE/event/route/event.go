package route

// route event
type Request struct {
    method     string
    path       string
    request    string
}

func RouteRequestEventGet(path string) Request {
    return Request {
        method:     get,
        path:       path,
    }
}

func RouteRequestEventPost(path, payload string) Request {
    return Request {
        method:     post,
        path:       path,
        request:    payload,
    }
}

func RouteRequestEventPut(path, payload string) Request {
    return Request {
        method:     put,
        path:       path,
        request:    payload,
    }
}

func RouteRequestEventDelete(path string) Request {
    return Request {
        method:     deleteh,
        path:       path,
    }
}
