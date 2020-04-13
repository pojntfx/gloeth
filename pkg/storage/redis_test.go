package storage

import (
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/go-redis/redis/v7"
)

func getTestData() [][2]string {
	return [][2]string{
		{
			"n1",
			"n4",
		},
		{
			"n2",
			"n5",
		},
		{
			"n3",
			"n6",
		},
		{
			"n6",
			"n5",
		},
		{
			"n5",
			"n4",
		},
		{
			"n8",
			"n7",
		},
		{
			"n7",
			"n4",
		},
		{
			"n7",
			"n5",
		},
		{
			"n10",
			"n9",
		},
		{
			"n9",
			"n6",
		},
		{
			"n9",
			"n7",
		},
	}
}

func flushTestData(client *redis.Client) error {
	return client.FlushAll().Err()
}

func applyTestData(testData [][2]string, client *redis.Client) error {
	for _, node := range testData {
		if err := client.Set(fmt.Sprintf("node:%v:%v", node[0], node[1]), true, 0).Err(); err != nil {
			return err
		}
	}

	return nil
}

func getAllTestData(client *redis.Client) ([][2]string, error) {
	keys, _, err := client.Scan(0, "node:*", 1000).Result()
	if err != nil {
		return nil, err
	}

	out := [][2]string{}
	for _, key := range keys {
		line := [2]string{}

		from := strings.Split(key, ":")
		line[0] = from[1]
		line[1] = from[2]

		out = append(out, line)
	}

	return out, nil
}

func getRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func TestNewRedis(t *testing.T) {
	raddr, err := net.ResolveTCPAddr("tcp", "localhost:6379")
	if err != nil {
		t.Error(err)
	}

	type args struct {
		addr     *net.TCPAddr
		password string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"New",
			args{
				raddr,
				"",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRedis(tt.args.addr, tt.args.password); got.client == nil {
				t.Errorf("NewRedis() client = %v, want !nil", got.client)
			}
		})
	}
}

func TestRedis_Apply(t *testing.T) {
	client := getRedisClient()
	deletions := [][2]string{}
	additions := getTestData()

	type fields struct {
		client *redis.Client
	}
	type args struct {
		deletions [][2]string
		additions [][2]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Apply",
			fields{
				client,
			},
			args{
				deletions,
				additions,
			},
			false,
		},
		{
			"Apply (deletions only)",
			fields{
				client,
			},
			args{
				deletions,
				additions,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Redis{
				client: tt.fields.client,
			}

			if len(tt.args.deletions) != 0 {
				if err := applyTestData(tt.args.deletions, r.client); err != nil {
					t.Error(err)
				}
			}

			if err := r.Apply(tt.args.deletions, tt.args.additions); (err != nil) != tt.wantErr {
				t.Errorf("Redis.Apply() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(tt.args.deletions) != 0 {
				data, err := getAllTestData(tt.fields.client)
				if err != nil {
					t.Error(err)
				}

				if len(data) != 0 {
					t.Errorf("Redis.Apply() deletions did not fully delete, got %v, want %v", data, [][2]string{})
				}
			}

			if err := flushTestData(r.client); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestRedis_GetAll(t *testing.T) {
	client := getRedisClient()
	testData := getTestData()

	type fields struct {
		client *redis.Client
	}
	tests := []struct {
		name      string
		fields    fields
		dataToAdd [][2]string
		want      [][2]string
		wantErr   bool
	}{
		{
			"GetAll",
			fields{
				client,
			},
			testData,
			testData,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := &Redis{
				client: tt.fields.client,
			}

			if err := applyTestData(tt.dataToAdd, r.client); err != nil {
				t.Error(err)
			}

			got, err := r.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			actualMatchLength := 0
			expectedMatchLength := len(tt.want)
			for _, ael := range got {
				for _, eel := range tt.want {
					if (ael[0] == eel[1] && ael[1] == eel[0]) || (ael[0] == eel[0] && ael[1] == eel[1]) {
						actualMatchLength = actualMatchLength + 1
					}
				}
			}

			if actualMatchLength != expectedMatchLength {
				t.Errorf("len(matches(Redis.GetAll())) = %v, want %v", actualMatchLength, expectedMatchLength)
			}

			if err := flushTestData(r.client); err != nil {
				t.Error(err)
			}
		})
	}
}
