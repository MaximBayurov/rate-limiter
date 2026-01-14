package bucket

type Bucket interface {
	Allow() bool
	CurrentLoad() int
}
