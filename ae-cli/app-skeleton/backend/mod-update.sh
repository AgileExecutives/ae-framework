go env -w GOPRIVATE=github.com/AgileExecutives/*
git config --global url."git@github.com:".insteadOf "https://github.com/"
go mod tidy