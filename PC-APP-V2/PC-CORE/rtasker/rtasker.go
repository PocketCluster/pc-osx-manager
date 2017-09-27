package rtasker

import (
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/service"
)

/*
 * RouteTasker is a short-lived task that is triggered by route path and continues to live for a certain period.
 * Examplary, well-suited tasks are installation, compose job management.
 */

type RouteTasker struct {
    route.Router
    service.ServiceSupervisor
}
