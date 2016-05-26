//go:generate rm -vf gen/gen.go
//go:generate mkdir -p static
//go:generate go-bindata -pkg gen -o gen/gen.go ./dashboard/dist/... ./db/migrations/...

package autoscale
