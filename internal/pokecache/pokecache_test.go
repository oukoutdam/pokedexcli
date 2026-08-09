package pokecache

import (
	"bytes"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second

	tt := map[string]struct {
		key   string
		value []byte
	}{
		"first page": {
			key:   "https://pokeapi.co/api/v2/location-area",
			value: []byte(`{"count":1533,"next":"https://pokeapi.co/api/v2/location-area?offset=2&limit=2","previous":null,"results":[{"name":"canalave-city-area","co/api/v2/location-area/1/"},{"name":"eterna-city-area","url":"https://pokeapi.co/api/v2/location-area/2/"}]}`),
		},
		"second page": {
			key:   "https://pokeapi.co/api/v2/location-area?offset=2&limit=2",
			value: []byte(`{"count":1533,"next":"https://pokeapi.co/api/v2/location-area?offset=4&limit=2","previous":"https://pokeapi.co/api/v2/location-area?offse:[{"name":"floaroma-town-area","url":"https://pokeapi.co/api/v2/location-area/3/"},{"name":"floaroma-meadow-area","url":"https://pokeapi.co/api/v2/location-area/4/"}]}`),
		},
		"empty": {
			key:   "https://example.com/empty",
			value: []byte{},
		},
		"empty key": {
			key:   "",
			value: []byte("value for the empty key"),
		},
		"binary value": {
			key:   "https://pokeapi.co/api/v2/sprites/1.png",
			value: []byte{0x89, 0x50, 0x4e, 0x47, 0x00, 0xff},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			c := NewCache(interval)
			c.Add(tc.key, tc.value)

			got, ok := c.Get(tc.key)
			if !ok {
				t.Fatalf("expected key %q to be found", tc.key)
			}
			if !bytes.Equal(got, tc.value) {
				t.Fatalf("expected %v, got %v", tc.value, got)
			}
		})
	}
}

func TestGetMissing(t *testing.T) {
	const interval = 5 * time.Second

	c := NewCache(interval)
	c.Add("present-key", []byte("present"))

	got, ok := c.Get("absent-key")
	if ok {
		t.Fatalf("expected absent-key to be not found, got %v", got)
	}

	if len(got) != 0 {
		t.Fatalf("expected empty value on miss, got %v", got)
	}
}

func TestAddOverwrites(t *testing.T) {
	const interval = 5 * time.Second
	const sharedKey = "shared-key"

	c := NewCache(interval)
	c.Add(sharedKey, []byte("first-value"))
	c.Add(sharedKey, []byte("second-value"))

	got, ok := c.Get(sharedKey)
	if !ok {
		t.Fatal("expected key to be found after overwrite")
	}

	if !bytes.Equal(got, []byte("second-value")) {
		t.Fatalf("expected %q, got %q", []byte("second-value"), got)
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond

	c := NewCache(baseTime)
	key := "http://example.com"
	value := []byte("test-data")
	c.Add(key, value)

	if _, ok := c.Get(key); !ok {
		t.Fatal("expected to find key before interval passes")
	}

	time.Sleep(waitTime)
	if _, ok := c.Get(key); ok {
		t.Fatal("expected not to find key after wait time has passed")
	}
}
