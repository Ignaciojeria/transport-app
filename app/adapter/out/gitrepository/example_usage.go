package gitrepository

import (
	"context"
	"fmt"
	"log"
	"os"
)

// ExampleCreateRepository demuestra cómo usar la función CreateRepositoryIfNotExists
func ExampleCreateRepository() {
	// Crear un repositorio en una nueva ubicación
	newRepoPath := "/tmp/mi-nuevo-repositorio"
	repo, err := CreateRepositoryIfNotExists(newRepoPath)
	if err != nil {
		log.Printf("Error creating repository: %v", err)
		return
	}

	fmt.Printf("Repositorio creado/abierto exitosamente en: %s\n", newRepoPath)
	fmt.Printf("Repositorio: %+v\n", repo)

	// Intentar crear el mismo repositorio otra vez (debería abrirlo)
	repo2, err := CreateRepositoryIfNotExists(newRepoPath)
	if err != nil {
		log.Printf("Error opening existing repository: %v", err)
		return
	}

	fmt.Printf("Repositorio existente abierto: %+v\n", repo2)
}

// ExampleGitRepositoryAdapter demuestra cómo usar el adaptador con el cliente inyectado
func ExampleGitRepositoryAdapter(gitRepoAdapter *GitRepositoryAdapter) {
	if gitRepoAdapter == nil {
		log.Println("GitRepositoryAdapter is not initialized (GIT_REPOSITORY_PATH not set)")
		return
	}
	ctx := context.Background()

	// Ejemplo 1: Obtener información del repositorio
	fmt.Println("=== Información del Repositorio ===")
	branch, err := gitRepoAdapter.GetCurrentBranch()
	if err != nil {
		log.Printf("Error getting current branch: %v", err)
	} else {
		fmt.Printf("Rama actual: %s\n", branch)
	}

	commit, err := gitRepoAdapter.GetCurrentCommit()
	if err != nil {
		log.Printf("Error getting current commit: %v", err)
	} else {
		fmt.Printf("Commit actual: %s\n", commit)
	}

	// Ejemplo 2: Verificar si el repositorio está sucio
	fmt.Println("\n=== Estado del Repositorio ===")
	isDirty, err := gitRepoAdapter.IsDirty()
	if err != nil {
		log.Printf("Error checking if dirty: %v", err)
	} else {
		if isDirty {
			fmt.Println("El repositorio tiene cambios sin commitear")
		} else {
			fmt.Println("El repositorio está limpio")
		}
	}

	// Ejemplo 3: Agregar archivos y hacer commit
	fmt.Println("\n=== Agregar y Commit ===")
	err = gitRepoAdapter.AddStage(".")
	if err != nil {
		log.Printf("Error adding files: %v", err)
	} else {
		fmt.Println("Archivos agregados al staging")
	}

	// Hacer commit
	commitOptions := &CommitOptions{
		All: true,
	}
	err = gitRepoAdapter.Commit("feat: agregar nueva funcionalidad", commitOptions)
	if err != nil {
		log.Printf("Error making commit: %v", err)
	} else {
		fmt.Println("Commit realizado exitosamente")
	}

	// Ejemplo 4: Hacer push (requiere configuración de autenticación)
	fmt.Println("\n=== Push ===")
	pushOptions := &PushOptions{
		SSHKeyPath: "/home/gitpod/.ssh/id_rsa", // SSH key en Gitpod
		Progress:   os.Stdout,
	}

	err = gitRepoAdapter.Push(ctx, pushOptions)
	if err != nil {
		log.Printf("Error pushing changes: %v", err)
	} else {
		fmt.Println("Push realizado exitosamente")
	}

	// Ejemplo 5: Hacer pull
	fmt.Println("\n=== Pull ===")
	pullOptions := &PullOptions{
		SSHKeyPath: "/home/gitpod/.ssh/id_rsa", // SSH key en Gitpod
		Progress:   os.Stdout,
	}

	err = gitRepoAdapter.Pull(ctx, pullOptions)
	if err != nil {
		log.Printf("Error pulling changes: %v", err)
	} else {
		fmt.Println("Pull realizado exitosamente")
	}
}
