package ecsta_test

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/fujiwara/ecsta"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var ignore = cmpopts.IgnoreUnexported(
	types.KeyValuePair{},
	types.Attachment{},
	types.Attribute{},
	types.ManagedAgent{},
	types.NetworkInterface{},
	types.Container{},
	types.EphemeralStorage{},
	types.ContainerOverride{},
	types.TaskOverride{},
	types.Task{},
)

func TestMarshalUnmarshalTask(t *testing.T) {
	b, err := os.ReadFile("testdata/task.json")
	if err != nil {
		t.Error(err)
	}
	var task types.Task
	if err := ecsta.UnmarshalJSONForAPI(b, &task); err != nil {
		t.Error(err)
	}
	if cpu := aws.ToString(task.Cpu); cpu != "256" {
		t.Errorf("unexpected cpu: %s", cpu)
	}
	if len(task.Containers) != 2 {
		t.Errorf("unexpected number of containers: %d", len(task.Containers))
	}
	if addr := aws.ToString(task.Containers[0].NetworkInterfaces[0].PrivateIpv4Address); addr != "10.3.1.230" {
		t.Errorf("unexpected private ipv4 address: %s", addr)
	}

	b2, err := ecsta.MarshalJSONForAPI(&task)
	if err != nil {
		t.Error(err)
	}
	var task2 types.Task
	if err := ecsta.UnmarshalJSONForAPI(b2, &task2); err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(task, task2, ignore); diff != "" {
		t.Error("not equal task", diff)
	}
}