package route

type ResponseFeeder interface {
    FeedResponseForGet(path, payload string) error
    FeedResponseForPost(path, payload string) error
    FeedResponseForPut(path, payload string) error
    FeedResponseForDelete(path, payload string) error
}