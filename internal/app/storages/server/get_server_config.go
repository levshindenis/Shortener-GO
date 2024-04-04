package server

func (serv *Server) GetServerConfig(param string) string {
	switch param {
	case "address":
		return serv.sc.GetStartAddress()
	case "baseURL":
		return serv.sc.GetShortBaseURL()
	case "file":
		return serv.sc.GetFilePath()
	case "db":
		return serv.sc.GetDBAddress()
	default:
		return ""
	}
}
