package main

import (
	rpc "github.com/mikerjacobi/grpc/pb"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CacheService struct {
	store map[string][]byte
}

func newCacheService() *CacheService {
	return &CacheService{make(map[string][]byte)}
}

func (s *CacheService) Get(ctx context.Context, req *rpc.GetReq) (*rpc.GetResp, error) {
	val, ok := s.store[req.Key]
	if !ok {
		logrus.Infof("key %s not found", req.Key)
		return nil, status.Errorf(codes.NotFound, "Key not found %s", req.Key)
	}
	return &rpc.GetResp{Val: val}, nil
}
func (s *CacheService) Store(ctx context.Context, req *rpc.StoreReq) (*rpc.StoreResp,
	error) {
	s.store[req.Key] = req.Val
	return &rpc.StoreResp{}, nil
}

type AccountService struct {
	store map[string]*rpc.Account
}

func newAccountService() *AccountService {
	return &AccountService{map[string]*rpc.Account{}}
}
func (s *AccountService) Create(ctx context.Context, newAccount *rpc.Account) (*rpc.Account, error) {
	if newAccount.Username == "" || newAccount.Password == "" {
		logrus.Errorf("acct create params invalid %+v", newAccount)
		return nil, status.Errorf(codes.InvalidArgument, "invalid acct creation")
	}
	if s.store[newAccount.Username] != nil {
		logrus.Errorf("username taken %s", newAccount.Username)
		return nil, status.Errorf(codes.InvalidArgument, "username taken")
	}
	newAccount.AccountID = uuid.NewV4().String()
	s.store[newAccount.AccountID] = newAccount
	s.store[newAccount.Username] = newAccount
	resp := &rpc.Account{Username: newAccount.Username, AccountID: newAccount.AccountID}
	logrus.Infof("username created %s", newAccount.Username)
	return resp, nil
}
func (s *AccountService) Get(ctx context.Context, account *rpc.Account) (*rpc.Account, error) {
	if acct := s.store[account.Username]; acct != nil {
		logrus.Infof("get user by username: %s", acct.Username)
		return acct, nil
	}
	if acct := s.store[account.AccountID]; acct != nil {
		logrus.Infof("get user by accountID: %s", acct.Username)
		return acct, nil
	}
	logrus.Errorf("account not found %+v", account)
	return nil, status.Errorf(codes.NotFound, "acct not found")
}
func (s *AccountService) Authenticate(ctx context.Context, account *rpc.Account) (*rpc.AuthenticateAccountResp, error) {
	acct := s.store[account.Username]
	if acct == nil {
		logrus.Errorf("account not found %+v", account)
		return &rpc.AuthenticateAccountResp{false}, nil
	} else if acct.Password != account.Password {
		logrus.Errorf("password mismatch %+v", account)
		return &rpc.AuthenticateAccountResp{false}, nil
	}
	logrus.Infof("logged in: %s", acct.Username)
	return &rpc.AuthenticateAccountResp{true}, nil
}
