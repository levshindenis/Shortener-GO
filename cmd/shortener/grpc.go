package main

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/levshindenis/sprint1/cmd/proto/shortener"
	"github.com/levshindenis/sprint1/internal/app/models"
	"github.com/levshindenis/sprint1/internal/app/storages/server"
)

// ShortenerServer поддерживает все необходимые методы сервера.
type ShortenerServer struct {
	pb.UnimplementedShortenerServer

	serv *server.Server
}

// AddURL - добавление короткого URL
func (s *ShortenerServer) AddURL(ctx context.Context, in *pb.AddURLRequest) (*pb.AddURLResponse, error) {
	var cookie string

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("UserID")
		if len(values) > 0 {
			cookie = values[0]
		}
	}

	if !s.serv.GetCookieStorage().InCookies(cookie) {
		s.serv.GetCookieStorage().SetValue(cookie)
	}

	address, flag, err := s.serv.MakeShortURL(in.LongURL)
	if err != nil {
		return nil, status.Error(codes.Aborted, "Something bad with MakeShortURL")
	}

	if !flag {
		if err = s.serv.GetStorage().SetData(address, in.LongURL, cookie); err != nil {
			return nil, status.Error(codes.Aborted, "Something bad with Save")
		}
		return &pb.AddURLResponse{ShortURL: s.serv.GetServerConfig().GetShortBaseURL() + "/" + address},
			status.Error(codes.OK, "Added")
	}

	return &pb.AddURLResponse{ShortURL: s.serv.GetServerConfig().GetShortBaseURL() + "/" + address},
		status.Error(codes.AlreadyExists, "Repeated URL")
}

// GetLongURL - возвращает длинный URL по короткому
func (s *ShortenerServer) GetLongURL(ctx context.Context, in *pb.GetLongURLRequest) (*pb.GetLongURLResponse, error) {
	result, isdeleted, err := s.serv.GetStorage().GetData(in.ShortURL, "key", "")
	if err != nil {
		return nil, status.Error(codes.Aborted, "Something bad with GetStorage")
	}

	if result == "" {
		return nil, status.Error(codes.InvalidArgument, "There is not longURL")
	}

	if isdeleted[0] {
		return nil, status.Error(codes.NotFound, "Deleted")
	}

	return &pb.GetLongURLResponse{LongURL: result}, status.Error(codes.OK, "Redirect")
}

// AddJsURL - добавление короткого URL в формате JSON
func (s *ShortenerServer) AddJsURL(ctx context.Context, in *pb.AddJsURLRequest) (*pb.AddJsURLResponse, error) {
	var cookie string

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("UserID")
		if len(values) > 0 {
			cookie = values[0]
		}
	}

	if !s.serv.GetCookieStorage().InCookies(cookie) {
		s.serv.GetCookieStorage().SetValue(cookie)
	}

	address, flag, err := s.serv.MakeShortURL(in.LongURL.GetUrl())
	if err != nil {
		return nil, status.Error(codes.Aborted, "Something bad with MakeShortURL")
	}

	if !flag {
		if err = s.serv.GetStorage().SetData(address, in.LongURL.GetUrl(), cookie); err != nil {
			return nil, status.Error(codes.Aborted, "Something bad with Save")
		}
		return &pb.AddJsURLResponse{ShortURL: &pb.JShortURL{Result: s.serv.GetServerConfig().GetShortBaseURL() + "/" + address}},
			status.Error(codes.OK, "Added")
	}

	return &pb.AddJsURLResponse{ShortURL: &pb.JShortURL{Result: s.serv.GetServerConfig().GetShortBaseURL() + "/" + address}},
		status.Error(codes.AlreadyExists, "Repeated URL")
}

// Ping - проверка соединения с БД
func (s *ShortenerServer) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	newCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := s.serv.GetDB().PingContext(newCtx); err != nil {
		return nil, status.Error(codes.Aborted, "Something bad with ping")
	}

	return &pb.PingResponse{}, status.Error(codes.OK, "Ping here")
}

// BatchURLs - множественное добавление коротких URL в JSON формате
func (s *ShortenerServer) BatchURLs(ctx context.Context, in *pb.BatchURLsRequest) (*pb.BatchURLsResponse, error) {
	var (
		cookie   string
		response pb.BatchURLsResponse
	)

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("UserID")
		if len(values) > 0 {
			cookie = values[0]
		}
	}

	if !s.serv.GetCookieStorage().InCookies(cookie) {
		s.serv.GetCookieStorage().SetValue(cookie)
	}

	for _, elem := range in.LongURLs {
		address, flag, err := s.serv.MakeShortURL(elem.OriginalUrl)
		if err != nil {
			return nil, status.Error(codes.Aborted, "Something bad with MakeShortURL")
		}

		if !flag {
			if err = s.serv.GetStorage().SetData(address, elem.OriginalUrl, cookie); err != nil {
				return nil, status.Error(codes.Aborted, "Something bad with Save")
			}
		}

		response.ShortURLs = append(response.ShortURLs,
			&pb.BShortURL{CorrelationId: elem.CorrelationId,
				ShortUrl: s.serv.GetServerConfig().GetShortBaseURL() + "/" + address})
	}

	return &response, status.Error(codes.OK, "Added")
}

// GetAllURLs - вернуть плоьзователю все сокращенные им URLs
func (s *ShortenerServer) GetAllURLs(ctx context.Context, in *pb.GetAllURLsRequest) (*pb.GetAllURLsResponse, error) {
	var (
		cookie   string
		response pb.GetAllURLsResponse
	)

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("UserID")
		if len(values) > 0 {
			cookie = values[0]
		}
	}

	if !s.serv.GetCookieStorage().InCookies(cookie) {
		return nil, status.Error(codes.Unauthenticated, "Bad cookie")
	}

	mystr, _, err := s.serv.GetStorage().GetData("", "all", cookie)
	if err != nil {
		return nil, status.Error(codes.Aborted, "Something bad with GetStorage")
	}
	if mystr == "" {
		return nil, status.Error(codes.NotFound, "No content")
	}

	myarr := strings.Split(mystr[:len(mystr)-1], "*")

	for i := 0; i < len(myarr); i += 2 {
		response.AllURLs = append(response.AllURLs,
			&pb.InfoURL{ShortUrl: s.serv.GetServerConfig().GetShortBaseURL() + "/" + myarr[i],
				OriginalUrl: myarr[i+1]})
	}

	return &response, status.Error(codes.OK, "Good")
}

// DeleteURLs - удаление записей из БД
func (s *ShortenerServer) DeleteURLs(ctx context.Context, in *pb.DeleteURLsRequest) (*pb.DeleteURLsResponse, error) {
	var cookie string

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("UserID")
		if len(values) > 0 {
			cookie = values[0]
		}
	}

	if !s.serv.GetCookieStorage().InCookies(cookie) {
		return nil, status.Error(codes.Unauthenticated, "Bad cookie")
	}

	for ind := range in.ShortURLs {
		s.serv.SetChan(models.DeleteValue{Value: in.ShortURLs[ind], Userid: cookie})
	}

	return &pb.DeleteURLsResponse{}, status.Error(codes.OK, "Deleted")
}

// Stats - возвращает статистику по количетсво пользователей и количетсву сокращенных URLs
func (s *ShortenerServer) Stats(ctx context.Context, in *pb.StatsRequest) (*pb.StatsResponse, error) {
	var (
		response pb.StatsResponse
		ip       string
	)

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("X-Real-IP")
		if len(values) > 0 {
			ip = values[0]
		}
	}

	if ip == "" || !s.serv.InCIDR(ip) {
		return nil, status.Error(codes.Canceled, "Bad IP")
	}

	data, err := s.serv.Stats()
	if err != nil {
		return nil, status.Error(codes.Aborted, "Something bad with Stats")
	}

	response.State.Users = int32(data.Users)
	response.State.Urls = int32(data.URLs)

	return &response, status.Error(codes.OK, "Stats ready")
}
