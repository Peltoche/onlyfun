package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Peltoche/onlyfun/assets"
	"github.com/Peltoche/onlyfun/internal/server"
	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/logger"
	"github.com/Peltoche/onlyfun/internal/tools/response"
	"github.com/Peltoche/onlyfun/internal/tools/router"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/Peltoche/onlyfun/internal/web/html"
	"github.com/spf13/afero"
)

var (
	ErrConflictTLSConfig = errors.New("can't use --self-signed-cert and --tls-key at the same time")
	ErrDevFlagRequire    = errors.New("this flag require the --dev flag setup")
)

type flags struct {
	LogLevel       string
	Folder         string
	TLSCert        string
	TLSKey         string
	HTTPHost       string
	HTTPHostnames  []string
	HTTPPort       int
	MemoryFS       bool
	SelfSignedCert bool
	Debug          bool
	Dev            bool
	HotReload      bool
	PrintVersion   bool
	PrintHelp      bool
}

func NewConfigFromFlags(flags *flags) (server.Config, error) {
	if flags.HotReload && !flags.Dev {
		return server.Config{}, fmt.Errorf("--hot-reload: %w", ErrDevFlagRequire)
	}

	if flags.MemoryFS && !flags.Dev {
		return server.Config{}, fmt.Errorf("--memory-fs: %w", ErrDevFlagRequire)
	}

	var logLevel slog.Level
	switch strings.ToLower(flags.LogLevel) {
	case "info":
		logLevel = slog.LevelInfo
	case "warn", "warning":
		logLevel = slog.LevelWarn
	case "err", "error":
		logLevel = slog.LevelError
	default:
		return server.Config{}, errors.New("invalid log level")
	}

	if flags.Debug {
		logLevel = slog.LevelDebug
	}

	var fs afero.Fs
	var storagePath string
	if flags.MemoryFS {
		fs = afero.NewMemMapFs()
		loadRequiredFilesIntoMemFS(fs)
		storagePath = ":memory:"
	} else {
		fs = afero.NewOsFs()
		storagePath = path.Join(flags.Folder, "db.sqlite")
	}

	err := fs.MkdirAll(flags.Folder, 0o755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return server.Config{}, fmt.Errorf("failed to create %q: %w", flags.Folder, err)
	}

	if flags.SelfSignedCert {
		if flags.TLSCert != "" || flags.TLSKey != "" {
			return server.Config{}, ErrConflictTLSConfig
		}

		flags.TLSCert, flags.TLSKey, err = generateSelfSignedCertificate(flags.HTTPHostnames, flags.Folder, fs)
		if err != nil {
			return server.Config{}, fmt.Errorf("failed to generate the self-signed certificate: %w", err)
		}
	}

	isTLSEnabled := flags.TLSCert != "" || flags.TLSKey != ""

	return server.Config{
		FS: fs,
		Listener: router.Config{
			Addr:      net.JoinHostPort(flags.HTTPHost, strconv.Itoa(flags.HTTPPort)),
			TLS:       isTLSEnabled,
			Secure:    !flags.Dev,
			CertFile:  flags.TLSCert,
			KeyFile:   flags.TLSKey,
			HostNames: flags.HTTPHostnames,
		},
		Storage: sqlstorage.Config{
			Path: storagePath,
		},
		Assets: assets.Config{
			HotReload: flags.HotReload,
		},
		Tools: tools.Config{
			Response: response.Config{
				PrettyRender: flags.Dev,
			},
			Log: logger.Config{
				Level:  logLevel,
				Output: os.Stderr,
			},
		},
		Folder: server.Folder(flags.Folder),
		HTML: html.Config{
			PrettyRender: flags.Dev,
			HotReload:    flags.HotReload,
		},
	}, nil
}

func generateSelfSignedCertificate(hostnames []string, folderPath string, fs afero.Fs) (string, string, error) {
	sslfolder := path.Join(folderPath, "ssl")
	certificatePath := path.Join(sslfolder, "cert.pem")
	privateKeyPath := path.Join(sslfolder, "key.pem")

	err := fs.MkdirAll(sslfolder, 0o700)
	if err != nil {
		return "", "", fmt.Errorf("failed to create the SSL folder: %w", err)
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", fmt.Errorf("failed to GenerateKey: %w", err)
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate serial number: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Duck Corp"},
		},
		DNSNames:  hostnames,
		NotBefore: time.Now(),
		// NotAfter:  time.Now().Add(3 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to create certificate: %w", err)
	}

	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if pemCert == nil {
		return "", "", errors.New("failed to encode certificate to PEM")
	}

	if err := afero.WriteFile(fs, certificatePath, pemCert, 0o644); err != nil {
		return "", "", fmt.Errorf("failed to write the certificate into the data folder: %w", err)
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", fmt.Errorf("unable to marshal private key: %w", err)
	}

	pemKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	if pemKey == nil {
		return "", "", errors.New("failed to encode key to PEM")
	}

	if err := afero.WriteFile(fs, privateKeyPath, pemKey, 0o600); err != nil {
		return "", "", fmt.Errorf("failed to write the certificate into the data folder: %w", err)
	}

	return certificatePath, privateKeyPath, nil
}

func loadRequiredFilesIntoMemFS(fs afero.Fs) error {
	err := loadFileInFS(fs, "/proc/uptime")
	if err != nil {
		return fmt.Errorf("failed to load /proc/uptime: %w", err)
	}

	err = loadFileInFS(fs, "/etc/hostname")
	if err != nil {
		return fmt.Errorf("failed to load /etc/hostname: %w", err)
	}

	return nil
}

func loadFileInFS(fs afero.Fs, path string) error {
	rawFile, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %q: %w", path, err)
	}

	err = afero.WriteFile(fs, path, rawFile, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write %q: %w", path, err)
	}

	return nil
}
