go env -w GOPRIVATE=github.com/AgileExecutives/*
git config --global url."git@github.com:".insteadOf "https://github.com/"
go clean -modcache

go get github.com/AgileExecutives/ae-framework/serverbase@main
go get github.com/AgileExecutives/ae-framework/shared-modules/organization@main
go get github.com/AgileExecutives/ae-framework/shared-modules/pdf@main
go get github.com/AgileExecutives/ae-framework/shared-modules/saas-base@main
go get github.com/AgileExecutives/ae-framework/shared-modules/static@main

go mod tidy