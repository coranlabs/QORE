package main

import "github.com/coranlabs/CORAN_LIB_NAS/internal/tools/generator"

func main() {
	generator.ParseSpecs()

	generator.GenerateNasMessage()

	generator.GenerateNasEncDec()

	generator.GenerateTestLarge()
}
