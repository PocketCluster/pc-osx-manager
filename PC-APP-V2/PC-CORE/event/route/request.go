package route

// route event
type Request struct {
    method     string
    path       string
    request    string
}

func RouteRequestGet(path string) Request {
    return Request {
        method:     get,
        path:       path,
    }
}

func RouteRequestPost(path, request string) Request {
    return Request {
        method:     post,
        path:       path,
        request:    request,
    }
}

func RouteRequestPut(path, request string) Request {
    return Request {
        method:     put,
        path:       path,
        request:    request,
    }
}

func RouteRequestDelete(path string) Request {
    return Request {
        method:     deleteh,
        path:       path,
    }
}
