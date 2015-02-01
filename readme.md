	# http://golang.org/doc/install
	go get github.com/andres-erbsen/mit/{rsa-keygen,getcert,evals/evals}
	rsa-keygen mit.sk
	env MIT_USER=kerberos MIT_ID=123456789 MIT_PASSWORD=password getcert mit.sk > mit.cert
	env MIT_SK=mit.sk MIT_CERT=mit.cert evals -subject=6.01
