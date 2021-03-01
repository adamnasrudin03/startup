package payment

import (
	"bwa-startup/users"
	"strconv"

	midtrans "github.com/veritrans/go-midtrans"
)

type Service interface {
	GetPaymentURL(transaction Transaction, user users.User)  (string, error)
}

type service struct {
}

func NewService() *service {
	return &service{}
}


func (s *service) GetPaymentURL(transaction Transaction, user users.User) (string, error) {
	midclient := midtrans.NewClient()
    midclient.ServerKey = "YOUR-VT-SERVER-KEY"
    midclient.ClientKey = "YOUR-VT-CLIENT-KEY"
    midclient.APIEnvType = midtrans.Sandbox

    snapGateway := midtrans.SnapGateway {
        Client: midclient,
    }

	snapReq := &midtrans.SnapReq{
        CustomerDetail: &midtrans.CustDetail{
            Email: user.Email,
			FName: user.Name,
        },
		TransactionDetails: midtrans.TransactionDetails{
			OrderID: strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
    } 
	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil
}