package billing

import (
	"io"
	"strconv"
	"time"

	"github.com/coranlabs/CORAN_CONSOLE/backend/factory"
	"github.com/coranlabs/CORAN_CONSOLE/backend/logger"
	"github.com/jlaffaye/ftp"
)

// The ftp client is for CDR Pull method, that is the billing domain actively query CDR file from CHF
func FTPLogin() (*ftp.ServerConn, error) {
	// FTP server is for CDR transfer
	billingConfig := factory.WebuiConfig.Configuration.BillingServer
	addr := billingConfig.HostIPv4 + ":" + strconv.Itoa(billingConfig.Port)

	var c *ftp.ServerConn

	c, err := ftp.Dial(addr, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	err = c.Login("admin", "Coran")
	if err != nil {
		return nil, err
	}

	logger.BillingLog.Info("Login FTP server")
	return c, err
}

func PullCDRFile(c *ftp.ServerConn, fileName string) ([]byte, error) {
	r, err := c.Retr(fileName)
	if err != nil {
		logger.BillingLog.Warn("Fail to Pull CDR file: ", fileName)
		return nil, err
	}

	defer func() {
		if err = r.Close(); err != nil {
			logger.BillingLog.Error(err)
		}
	}()

	logger.BillingLog.Info("Pull CDR file success")

	if err = c.Quit(); err != nil {
		return nil, err
	}

	cdr, err_read := io.ReadAll(r)

	return cdr, err_read
}
