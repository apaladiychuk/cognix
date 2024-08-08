package main

import (
	"cognix.ch/api/v2/core/proto"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
)

const nsSecret = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

type K8SServer struct {
	proto.UnsafeConfigMapServer
	namespace string
	client    *kubernetes.Clientset
}

func (k *K8SServer) GetList(ctx context.Context, r *proto.ConfigMapList) (*proto.ConfigMapListResponse, error) {
	configMap, err := k.client.CoreV1().ConfigMaps(k.namespace).Get(ctx, r.GetName(), v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	result := proto.ConfigMapListResponse{
		Values: make([]*proto.ConfigMapRecord, 0),
	}
	for key, value := range configMap.Data {
		result.Values = append(result.Values, &proto.ConfigMapRecord{
			Key:   key,
			Value: value,
		})
	}
	return &result, nil
}

func (k *K8SServer) Save(ctx context.Context, r *proto.ConfigMapSave) (*empty.Empty, error) {
	configMap, err := k.client.CoreV1().ConfigMaps(k.namespace).Get(ctx, r.GetName(), v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	configMap.Data[r.GetValue().GetKey()] = r.GetValue().GetValue()
	if _, err = k.client.CoreV1().ConfigMaps(k.namespace).Update(ctx, configMap, v1.UpdateOptions{}); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (k *K8SServer) Delete(ctx context.Context, r *proto.ConfigMapDelete) (*empty.Empty, error) {
	configMap, err := k.client.CoreV1().ConfigMaps(k.namespace).Get(ctx, r.GetName(), v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	delete(configMap.Data, r.GetKey())

	if _, err = k.client.CoreV1().ConfigMaps(k.namespace).Update(ctx, configMap, v1.UpdateOptions{}); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func NewK8SServer() (proto.ConfigMapServer, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	ns, err := getCurrentNamespace()
	if err != nil {
		return nil, err
	}
	return &K8SServer{
		client:    client,
		namespace: ns,
	}, nil
}

func getCurrentNamespace() (string, error) {
	ns, err := os.ReadFile(nsSecret)
	if err != nil {
		return "", err
	}
	return string(ns), nil
}
