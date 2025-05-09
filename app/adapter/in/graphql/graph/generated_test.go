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

		files := []string{
			filepath.Join(moduleRoot, "app", "adapter", "in", "graphql", "graph", "generated.go"),
			filepath.Join(moduleRoot, "app", "adapter", "in", "graphql", "graph", "model", "models_gen.go"),
		}

		// Crear backups y planear rollback
		rollback := func() {
			for _, f := range files {
				bak := f + ".bak"
				if _, err := os.Stat(bak); err == nil {
					_ = os.Remove(f)
					_ = os.Rename(bak, f)
				}
			}
		}
		defer rollback()

		for _, f := range files {
			bak := f + ".bak"
			if _, err := os.Stat(f); err == nil {
				Expect(os.Rename(f, bak)).To(Succeed(), "Failed to rename %s to %s", f, bak)
			}
		}

		// Ejecutar gqlgen
		cmd := exec.Command("go", "run", "github.com/99designs/gqlgen", "generate")
		cmd.Dir = moduleRoot
		output, err := cmd.CombinedOutput()
		Expect(err).To(BeNil(), "gqlgen generate failed:\n%s", string(output))

		// Comparar
		for _, f := range files {
			bak := f + ".bak"

			if _, err := os.Stat(bak); os.IsNotExist(err) {
				continue
			}

			newData, err := os.ReadFile(f)
			Expect(err).To(BeNil(), "Missing regenerated file: %s", f)

			oldData, err := os.ReadFile(bak)
			Expect(err).To(BeNil(), "Missing backup file: %s", bak)

			if normalizeCode(newData) != normalizeCode(oldData) {
				fmt.Println("===== DIFF DETECTED =====")
				fmt.Printf("❌ File %s is out of sync with generated version.\n", f)
				Fail("Run `gqlgen generate` and commit the updated generated files.")
			}
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
