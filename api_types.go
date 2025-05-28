package plscli

type RegisterRequest struct {
	DeployName string `json:"deploy_name"`
}

type RegisterResponse struct {
	ClientId string `json:"client_id"`
}

type DeleteRequest struct {
	DeployName string `json:"deploy_name"`
	ClientId   string `json:"client_id"`
}

type DeleteResponse struct {
	ClientId string `json:"client_id"`
}

type LeaderRequest struct {
	DeployName string `json:"deploy_name"`
	ClientId   string `json:"client_id"`
}

type LeaderResponse struct {
	ClientId string `json:"client_id"`
}
