//go:generate rm -vf gen/gen.go
//go:generate mkdir -p static
//go:generate go-bindata -pkg gen -o gen/gen.go ./static/... ./db/migrations/...

package autoscale
