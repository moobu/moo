package raw

import (
	"testing"
	"time"

	"github.com/moobu/moo/builder"
	"github.com/moobu/moo/runtime/vanilla"
)

func TestWait(t *testing.T) {
	d := New()
	r := vanilla.Runnable{
		Bundle: &builder.Bundle{
			Entry: []string{"sleep"},
		},
		Args: []string{"3"},
	}
	p, err := d.Fork(&r)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.Wait(p); err != nil {
		t.Fatal(err)
	}
}

func TestKill(t *testing.T) {
	d := New()
	r := vanilla.Runnable{
		Bundle: &builder.Bundle{
			Entry: []string{"sleep"},
		},
		Args: []string{"10"},
	}
	p, err := d.Fork(&r)
	if err != nil {
		t.Fatal(err)
	}
	s := time.Now()
	go d.Wait(p)
	time.Sleep(time.Second)
	if err := d.Kill(p); err != nil {
		t.Fatal(err)
	}
	if time.Since(s) > time.Second*2 {
		t.Fatal("process should have been killed in 2 seconds")
	}
}
