package main

import (
	"bytes"
	"cognix.ch/api/v2/core/proto"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"os"
	"strings"
)

type DockerServer struct {
	proto.UnsafeConfigMapServer
	path string
}

func (k *DockerServer) readConfig(name string) ([]string, error) {
	buf, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	buf = bytes.ReplaceAll(buf, []byte("\r"), []byte(""))

	rows := strings.Split(string(buf), "\n")
	return rows, nil

}
func (k *DockerServer) GetList(ctx context.Context, r *proto.ConfigMapList) (*proto.ConfigMapListResponse, error) {
	filename := fmt.Sprintf("%s/%s.env", k.path, r.Name)

	rows, err := k.readConfig(filename)
	if err != nil {
		return nil, err
	}

	result := proto.ConfigMapListResponse{
		Values: make([]*proto.ConfigMapRecord, 0),
	}
	for _, row := range rows {
		value := strings.Split(row, "=")
		if len(value) != 2 {
			continue
		}
		result.Values = append(result.Values, &proto.ConfigMapRecord{
			Key:   value[0],
			Value: value[1],
		})
	}
	return &result, nil
}

func (k *DockerServer) Save(ctx context.Context, r *proto.ConfigMapSave) (*empty.Empty, error) {
	filename := fmt.Sprintf("%s/%s.env", k.path, r.Name)
	rows, err := k.readConfig(filename)
	if err != nil {
		return nil, err
	}
	isExists := false
	for i, row := range rows {
		value := strings.Split(row, "=")
		if len(value) != 2 {
			continue
		}
		if value[0] == r.Value.Key {
			rows[i] = fmt.Sprintf("%s=%s", r.Value.Key, r.Value.Value)
			isExists = true
			break
		}
	}
	if !isExists {
		rows = append(rows, fmt.Sprintf("%s=%s", r.Value.Key, r.Value.Value))
	}
	return &empty.Empty{}, os.WriteFile(filename, []byte(strings.Join(rows, "\n")), 0644)
}

func (k *DockerServer) Delete(ctx context.Context, r *proto.ConfigMapDelete) (*empty.Empty, error) {
	filename := fmt.Sprintf("%s/%s.env", k.path, r.Name)
	rows, err := k.readConfig(filename)
	if err != nil {
		return nil, err
	}
	newRows := make([]string, 0)

	for _, row := range rows {
		value := strings.Split(row, "=")
		if len(value) == 2 && value[0] == r.Key {
			continue
		}
		newRows = append(newRows, row)
	}
	return &empty.Empty{}, os.WriteFile(filename, []byte(strings.Join(newRows, "\n")), 0644)
}

func NewDockerServer(root string) (proto.ConfigMapServer, error) {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return nil, fmt.Errorf("can not find docker config maps")
	}
	return &DockerServer{
		path: root,
	}, nil
}
