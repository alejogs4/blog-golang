go test -run Unit ./... -coverprofile coverage.out
go tool cover -func coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}'