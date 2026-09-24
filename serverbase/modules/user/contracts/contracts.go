package contracts

import "embed"

// Files packages user template contracts with the module so tenant lifecycle
// hooks do not depend on the application's working directory.
//
//go:embed *.json
var Files embed.FS
