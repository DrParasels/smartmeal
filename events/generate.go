package events

//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.12
//go:generate protoc --go_out=. --go_opt=module=github.com/yourname/smartmeal/events proto/meal_created.proto
