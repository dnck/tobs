package tests

import (
	"net"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

func testTimescaleGetPassword(t testing.TB, user string) {
	var getpass *exec.Cmd

	if user == "" {
		t.Logf("Running 'ts-obs timescaledb get-password'")
		getpass = exec.Command("ts-obs", "timescaledb", "get-password", "-n", RELEASE_NAME, "--namespace", NAMESPACE)
	} else {
		t.Logf("Running 'ts-obs timescaledb get-password -U %v'\n", user)
		getpass = exec.Command("ts-obs", "timescaledb", "get-password", "-U", user, "-n", RELEASE_NAME, "--namespace", NAMESPACE)
	}

	out, err := getpass.CombinedOutput()
	if err != nil {
		t.Logf(string(out))
		t.Fatal(err)
	}
}

func testTimescalePortForward(t testing.TB, port string) {
	var portforward *exec.Cmd

	if port == "" {
		t.Logf("Running 'ts-obs timescaledb port-forward'")
		portforward = exec.Command("ts-obs", "timescaledb", "port-forward", "-n", RELEASE_NAME, "--namespace", NAMESPACE)
	} else {
		t.Logf("Running 'ts-obs timescaledb port-forward -p %v'\n", port)
		portforward = exec.Command("ts-obs", "timescaledb", "port-forward", "-p", port, "-n", RELEASE_NAME, "--namespace", NAMESPACE)
	}

	err := portforward.Start()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(4 * time.Second)

	if port == "" {
		port = "5432"
	}

	_, err = net.DialTimeout("tcp", "localhost:"+port, time.Second)
	if err != nil {
		t.Fatal(err)
	}

	portforward.Process.Signal(syscall.SIGINT)

}

func testTimescaleConnect(t testing.TB, master bool, user string) {
	var connect *exec.Cmd

	if master {
		t.Logf("Running 'ts-obs timescaledb connect -m'")
		connect = exec.Command("ts-obs", "timescaledb", "connect", "-m", "-n", RELEASE_NAME, "--namespace", NAMESPACE)
	} else {
		if user == "" {
			t.Logf("Running 'ts-obs timescaledb connect'\n")
			connect = exec.Command("ts-obs", "timescaledb", "connect", "-n", RELEASE_NAME, "--namespace", NAMESPACE)
		} else {
			t.Logf("Running 'ts-obs timescaledb connect -U %v'\n", user)
			connect = exec.Command("ts-obs", "timescaledb", "connect", "-U", user, "-n", RELEASE_NAME, "--namespace", NAMESPACE)
		}
	}

	err := connect.Start()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Second)
	connect.Process.Signal(syscall.SIGINT)

	deletepod := exec.Command("kubectl", "delete", "pods", "psql")
	deletepod.Run()
}

func TestTimescale(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TimescaleDB tests")
	}

	testTimescaleGetPassword(t, "")
	testTimescaleGetPassword(t, "admin")

	testTimescalePortForward(t, "")
	testTimescalePortForward(t, "5432")
	testTimescalePortForward(t, "1789")
	testTimescalePortForward(t, "1030")
	testTimescalePortForward(t, "2389")

	testTimescaleConnect(t, true, "")
	testTimescaleConnect(t, false, "")
	testTimescaleConnect(t, false, "postgres")
	testTimescaleConnect(t, false, "postgres")
	testTimescaleConnect(t, false, "admin")
	testTimescaleConnect(t, false, "admin")
}