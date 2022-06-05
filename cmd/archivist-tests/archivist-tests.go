package main

import (
	"github.com/voodooEntity/archivist"
)

func main() {
	archivist.Init("debug", "file", "out.lot")
	archivist.Error("Test")
}
