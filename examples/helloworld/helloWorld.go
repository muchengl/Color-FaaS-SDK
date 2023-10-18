package main

import cfaas "Color-FaaS-SDK/pkg"

func helloWorld(ctx cfaas.FaaSContext, req cfaas.FuncRequest) (cfaas.FuncResponse, error) {
	return cfaas.FuncResponse("hello world, your msg is : " + string(req)), nil
}

func main() {
	cfaas.Run(helloWorld)
}
