package service

import (
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
)

type WebServerBinCMDService interface {
	Exec(servername string, arg ...string) (isSuccessful bool, stdout, stderr string, err error)
}

type webServerBinCMDService struct {
	eps epv1.WebServerBinCMDEndpoints
}

func (w *webServerBinCMDService) Exec(servername string, arg ...string) (bool, string, string, error) {
	resp, err := w.eps.EndpointExec()(GetContext(), &v1.ExecuteRequest{
		ServerName: servername,
		Args:       arg,
	})
	if resp != nil {
		response := resp.(*v1.ExecuteResponse)
		return response.Successful, string(response.StandardOutput), string(response.StandardError), err
	}
	return false, "", "", err
}

func newWebServerBinCMDService(factory *factory) WebServerBinCMDService {
	return &webServerBinCMDService{eps: factory.eps.WebServerBinCMD()}
}
