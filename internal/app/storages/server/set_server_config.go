package server

func (serv *Server) SetServerConfig(value string, param string) {
	switch param {
	case "address":
		serv.sc.SetStartAddress(value)
	case "baseURL":
		serv.sc.SetShortBaseURL(value)
	case "file":
		serv.sc.SetFilePath(value)
	case "db":
		serv.sc.SetDBAddress(value)
	default:
		break
	}
}
