package worker

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kataras/iris"
	"github.com/megaease/easegateway/pkg/api"
	"github.com/megaease/easegateway/pkg/logger"
	"github.com/megaease/easegateway/pkg/object/meshcontroller/registrycenter"
)

func (w *Worker) nacosAPIs() []*apiEntry {
	APIs := []*apiEntry{
		{
			Path:    meshNacosPrefix + "/ns/instance/list",
			Method:  "GET",
			Handler: w.nacosInstanceList,
		},
		{
			Path:    meshNacosPrefix + "/ns/instance",
			Method:  "POST",
			Handler: w.nacosRegister,
		},
		{
			Path:    meshNacosPrefix + "/ns/instance",
			Method:  "DELETE",
			Handler: w.emptyHandler,
		},
		{
			Path:    meshNacosPrefix + "/ns/instance/beat",
			Method:  "PUT",
			Handler: w.emptyHandler,
		},
		{
			Path:    meshNacosPrefix + "/ns/instance",
			Method:  "PUT",
			Handler: w.emptyHandler,
		},
		{
			Path:    meshNacosPrefix + "/ns/instance",
			Method:  "GET",
			Handler: w.nacosInstance,
		},
		{
			Path:    meshNacosPrefix + "/ns/service/list",
			Method:  "GET",
			Handler: w.nacosServiceList,
		},
		{
			Path:    meshNacosPrefix + "/ns/service",
			Method:  "GET",
			Handler: w.nacosService,
		},
	}

	return APIs
}

func (w *Worker) nacosRegister(ctx iris.Context) {
	err := w.registryServer.CheckRegistryURL(ctx)
	if err != nil {
		api.HandleAPIError(ctx, http.StatusBadRequest,
			fmt.Errorf("parse request url parameters failed: %v", err))
		return
	}

	serviceSpec := w.service.GetServiceSpec(w.serviceName)
	if serviceSpec == nil {
		err := fmt.Errorf("registry to unknown service: %s", w.serviceName)
		api.HandleAPIError(ctx, http.StatusBadRequest, err)
		return
	}

	w.registryServer.Register(serviceSpec, w.ingressServer.Ready, w.egressServer.Ready)
}

func (w *Worker) nacosInstanceList(ctx iris.Context) {
	serviceName := ctx.Params().Get("serviceName")
	if len(serviceName) == 0 {
		api.HandleAPIError(ctx, http.StatusBadRequest,
			fmt.Errorf("empty serviceName in url parameters"))
		return
	}
	var (
		err         error
		serviceInfo *registrycenter.ServiceRegistryInfo
	)

	if serviceInfo, err = w.registryServer.DiscoveryService(serviceName); err != nil {
		logger.Errorf("discovery service: %s, err: %v ", serviceName, err)
		api.HandleAPIError(ctx, http.StatusInternalServerError, err)
		return
	}

	nacosSvc := w.registryServer.ToNacosService(serviceInfo)

	buff, err := json.Marshal(nacosSvc)
	if err != nil {
		logger.Errorf("json marshal nacosService: %#v err: %v", nacosSvc, err)
		api.HandleAPIError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Header("Content-Type", registrycenter.ContentTypeJSON)
	ctx.Write(buff)
}

func (w *Worker) nacosInstance(ctx iris.Context) {
	serviceName := ctx.Params().Get("serviceName")
	if len(serviceName) == 0 {
		api.HandleAPIError(ctx, http.StatusBadRequest,
			fmt.Errorf("empty serviceName in url parameters"))
		return
	}
	var (
		err         error
		serviceInfo *registrycenter.ServiceRegistryInfo
	)

	if serviceInfo, err = w.registryServer.DiscoveryService(serviceName); err != nil {
		logger.Errorf("discovery service: %s, err: %v ", serviceName, err)
		api.HandleAPIError(ctx, http.StatusInternalServerError, err)
		return
	}

	nacosIns := w.registryServer.ToNacosInstanceInfo(serviceInfo)

	buff, err := json.Marshal(nacosIns)
	if err != nil {
		logger.Errorf("json marshal nacosInstance: %#v err: %v", nacosIns, err)
		api.HandleAPIError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Header("Content-Type", registrycenter.ContentTypeJSON)
	ctx.Write(buff)
}

func (w *Worker) nacosServiceList(ctx iris.Context) {
	var (
		err          error
		serviceInfos []*registrycenter.ServiceRegistryInfo
	)
	if serviceInfos, err = w.registryServer.Discovery(); err != nil {
		logger.Errorf("discovery services err: %v ", err)
		api.HandleAPIError(ctx, http.StatusInternalServerError, err)
		return
	}
	serviceList := w.registryServer.ToNacosServiceList(serviceInfos)

	buff, err := json.Marshal(serviceList)
	if err != nil {
		logger.Errorf("json marshal serviceList: %#v err: %v", serviceList, err)
		api.HandleAPIError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Header("Content-Type", registrycenter.ContentTypeJSON)
	ctx.Write(buff)
}

func (w *Worker) nacosService(ctx iris.Context) {
	serviceName := ctx.Params().Get("serviceName")
	if len(serviceName) == 0 {
		api.HandleAPIError(ctx, http.StatusBadRequest,
			fmt.Errorf("empty serviceName in url parameters"))
		return
	}
	var (
		err         error
		serviceInfo *registrycenter.ServiceRegistryInfo
	)

	if serviceInfo, err = w.registryServer.DiscoveryService(serviceName); err != nil {
		logger.Errorf("discovery service: %s, err: %v ", serviceName, err)
		api.HandleAPIError(ctx, http.StatusInternalServerError, err)
		return
	}

	nacosSvcDetail := w.registryServer.ToNacosServiceDetail(serviceInfo)

	buff, err := json.Marshal(nacosSvcDetail)
	if err != nil {
		logger.Errorf("json marshal nacosSvcDetail: %#v err: %v", nacosSvcDetail, err)
		api.HandleAPIError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Header("Content-Type", registrycenter.ContentTypeJSON)
	ctx.Write(buff)
}
