package integration

import (
	"context"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// StartServer builds and starts the server-test binary in the server-test dir.
// It returns the server directory and a cleanup function to stop the server.
func StartServer(t *testing.T) (serverDir string, cleanup func()) {
	t.Helper()
	// determine server-test dir relative to current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("pwd: %v", err)
	}
	candidates := []string{
		filepath.Join(cwd, "server-test"),
		filepath.Join(filepath.Dir(cwd), "server-test"),
		filepath.Join(filepath.Dir(filepath.Dir(cwd)), "server-test"),
	}
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && fi.IsDir() {
			serverDir = c
			break
		}
	}
	if serverDir == "" {
		t.Fatalf("could not locate server-test dir from %s", cwd)
	}

	build := exec.Command("go", "build", "-o", "server-test-bin", ".")
	build.Dir = serverDir
	build.Env = os.Environ()
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build server failed: %v\n%s", err, string(out))
	}

	cmd := exec.Command("./server-test-bin")
	cmd.Dir = serverDir
	// Force deterministic test environment: in-memory DB and known admin creds
	env := os.Environ()
	env = append(env, "FEATURE_EMAIL_VERIFICATION=false")
	env = append(env, "USE_IN_MEMORY_DB=true")
	env = append(env, "ADMIN_USER=testuser@unburdy.de")
	env = append(env, "ADMIN_PASSWORD=newpass123")
	env = append(env, "MOCK_EMAIL=true")
	cmd.Env = env
	if err := cmd.Start(); err != nil {
		t.Fatalf("start server failed: %v", err)
	}

	// wait for health
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client := &http.Client{}
	healthURL := "http://localhost:8080/api/v1/health"
	ok := false
	for {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, healthURL, nil)
		resp, err := client.Do(req)
		if err == nil && resp != nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				ok = true
				break
			}
		}
		select {
		case <-ctx.Done():
			break
		default:
			time.Sleep(200 * time.Millisecond)
		}
	}
	if !ok {
		_ = cmd.Process.Kill()
		t.Fatalf("server did not become healthy in time")
	}

	cleanup = func() {
		_ = cmd.Process.Kill()
		cmd.Wait()
		_ = os.Remove(filepath.Join(serverDir, "server-test-bin"))
	}
	return serverDir, cleanup
}
