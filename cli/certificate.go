package cli

import (
	"crypto/tls"
	"log"
	"os"
	"sync"
	"time"

	"github.com/qlyco/qwc/internal/res"
)

type certificateFiles struct {
	certFile    string
	keyFile     string
	lock        sync.Mutex
	certificate *tls.Certificate
	certTime    time.Time
	keyTime     time.Time
}

func loadCertificate(certFile string, keyFile string) (*certificateFiles, error) {
	files := &certificateFiles{certFile: certFile, keyFile: keyFile}
	certificate, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		return nil, err
	}

	files.certificate = &certificate
	files.certTime, files.keyTime = files.modTimes()

	return files, nil
}

func (files *certificateFiles) get(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	files.lock.Lock()
	defer files.lock.Unlock()

	certTime, keyTime := files.modTimes()

	if certTime.Equal(files.certTime) && keyTime.Equal(files.keyTime) {
		return files.certificate, nil
	}

	files.certTime, files.keyTime = certTime, keyTime
	certificate, err := tls.LoadX509KeyPair(files.certFile, files.keyFile)

	if err != nil {
		log.Printf(res.KeptCertificateMessage, err)
		return files.certificate, nil
	}

	files.certificate = &certificate
	log.Printf(res.LoadedCertificateMessage, files.certFile)

	return files.certificate, nil
}

func (files *certificateFiles) modTimes() (time.Time, time.Time) {
	return modTime(files.certFile), modTime(files.keyFile)
}

func modTime(filePath string) time.Time {
	info, err := os.Stat(filePath)

	if err != nil {
		return time.Time{}
	}

	return info.ModTime()
}
