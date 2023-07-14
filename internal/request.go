package request

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cloudflare/cfssl/cli/genkey"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/signer"
)

type Request struct {
	Signer        string `yaml:"signer"`
	Profile       string `yaml:"profile"`
	Outdir        string `yaml:"outdir"`
	Renew_seconds int    `yaml:"renew_seconds"`
	Perms         Perms  `yaml:"perms"`
	CSR           CSR    `yaml:"csr"`
	Reload        Reload `yaml:"reload"`
	certfile      string
	keyfile       string
	csrfile       string
}

func (r Request) String() string {
	return r.CSR.CN
}

func (r *Request) Init_files() error {
	if err := os.MkdirAll(r.Outdir, os.ModePerm); err != nil {
		return err
	}
	os.Chown(r.Outdir, r.Perms.uid(), r.Perms.gid())
	r.CSR.write(r.Outdir, r.Perms)
	r.certfile = fmt.Sprintf("%s/%s.cert", r.Outdir, r.CSR.CN)
	r.keyfile = fmt.Sprintf("%s/%s.key", r.Outdir, r.CSR.CN)
	r.csrfile = fmt.Sprintf("%s/%s.csr", r.Outdir, r.CSR.CN)
	return nil
}

func (r *Request) Set_defaults() {
	// TODO: we should ass in some config class to this
	if r.Signer == "" {
		r.Signer = "discovery"
	}
	if r.Outdir == "" {
		r.Outdir = "/etc/cfssl/csr"
	}
	if r.Renew_seconds == 0 {
		r.Renew_seconds = 84600
	}
	r.CSR.Key.Set_defaults()
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (r *Request) Certfile_exist() bool {
	return fileExists(r.certfile)
}
func (r *Request) Key_exist() bool {
	return fileExists(r.certfile)
}

func (r *Request) Gencert(s signer.Signer) error {
	var key, csrBytes []byte
	req := csr.CertificateRequest{
		KeyRequest: csr.NewKeyRequest(),
	}
	csrBytes, err := r.CSR.json()
	if err != nil {
		return err
	}
	err = json.Unmarshal(csrBytes, &req)
	if err != nil {
		return err
	}
	req.CN = r.CSR.CN
	req.Hosts = r.CSR.Hosts
	g := &csr.Generator{Validator: genkey.Validator}
	csrBytes, key, err = g.ProcessRequest(&req)
	if err != nil {
		key = nil
		return err
	}
	signReq := signer.SignRequest{
		Request: string(csrBytes),
		Hosts:   r.CSR.Hosts,
		Profile: r.Profile,
		Label:   r.Signer,
	}
	cert, err := s.Sign(signReq)
	if err != nil {
		return err
	}
	println("write file {}", r.certfile)
	err = os.WriteFile(r.certfile, cert, r.Perms.mode())
	if err != nil {
		return err
	}
	err = os.WriteFile(r.keyfile, key, r.Perms.mode())
	if err != nil {
		return err
	}
	err = os.WriteFile(r.csrfile, csrBytes, r.Perms.mode())
	if err != nil {
		return err
	}

	return nil
}

/*
func (r *Request) Sign(s signer.Signer) error {
	csr, err := r.CSR.json()
	if err != nil {
		return err
	}
	req := signer.SignRequest{
		Hosts:   r.CSR.Hosts,
		Request: string(csr),
		Subject: r.CSR.CN,
		Profile: r.Profile,
		Label:   r.Signer,
	}
	cert, err := s.Sign(req)
	if err != nil {
		return err
	}
	// write cert file
	return nil
}
*/
