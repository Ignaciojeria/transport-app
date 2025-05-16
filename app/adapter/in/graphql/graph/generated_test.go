package graph

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GraphQL Code Generation", func() {
	It("should validate that gqlgen generate was executed", func() {
		moduleRoot := findModuleRoot()
		modelsFile := filepath.Join(moduleRoot, "app", "adapter", "in", "graphql", "graph", "model", "models_gen.go")

		bak := modelsFile + ".bak"
		Expect(os.Rename(modelsFile, bak)).To(Succeed(), "Failed to rename %s to %s", modelsFile, bak)

		defer func() {
			_ = os.Remove(modelsFile)
			_ = os.Rename(bak, modelsFile)
		}()

		cmd := exec.Command("go", "run", "github.com/99designs/gqlgen", "generate")
		cmd.Dir = moduleRoot
		output, err := cmd.CombinedOutput()
		Expect(err).To(BeNil(), "gqlgen generate failed:\n%s", string(output))

		newData, err := os.ReadFile(modelsFile)
		Expect(err).To(BeNil(), "Missing regenerated models file")

		oldData, err := os.ReadFile(bak)
		Expect(err).To(BeNil(), "Missing backup models file")

		if normalizeCode(newData) != normalizeCode(oldData) {
			fmt.Println("===== DIFF DETECTED in models_gen.go =====")
			Fail("Run `gqlgen generate` and commit the updated models_gen.go")
		}
	})
})

func normalizeCode(data []byte) string {
	lines := strings.Split(string(data), "\n")
	filtered := []string{}
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if trim == "" ||
			strings.HasPrefix(trim, "// Code generated") ||
			strings.HasPrefix(trim, "//go:") {
			continue
		}
		filtered = append(filtered, trim)
	}
	return strings.Join(filtered, "\n")
}

func findModuleRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		next := filepath.Dir(dir)
		if next == dir {
			Fail("go.mod not found")
		}
		dir = next
	}
}
