package route

// route event
type Request struct {
    method     string
    path       string
    request    string
}

func RouteRequestGet(path string) Request {
    return Request {
        method: RouteMethodGet,
        path:   path,
    }
}

func RouteRequestPost(path, request string) Request {
    return Request {
        method:  RouteMethodPost,
        path:    path,
        request: request,
    }
}

func RouteRequestPut(path, request string) Request {
    return Request {
        method:  RouteMethodPut,
        path:    path,
        request: request,
    }
}

func RouteRequestDelete(path string) Request {
    return Request {
        method: RouteMethodDeleteh,
        path:   path,
    }
}
