DEPS := \
	abodemine-ca-key.pem \
	abodemine-ca.csr \
	abodemine-ca.pem \
	abodemine.local-peer-chain.pem \
	abodemine.local-peer-key.pem \
	abodemine.local-peer.csr \
	abodemine.local-peer.pem \
	abodemine.local-peer.pfx \
	abodemine.test-peer-chain.pem \
	abodemine.test-peer-key.pem \
	abodemine.test-peer.csr \
	abodemine.test-peer.pem \
	abodemine.test-peer.pfx \
	alpha.test-peer-chain.pem \
	alpha.test-peer-key.pem \
	alpha.test-peer.csr \
	alpha.test-peer.pem \
	alpha.test-peer.pfx \
	omega.test-peer-chain.pem \
	omega.test-peer-key.pem \
	omega.test-peer.csr \
	omega.test-peer.pem \
	omega.test-peer.pfx

all: $(DEPS)

clean:
	rm -f $(DEPS)

%-ca-key.pem %-ca.csr %-ca.pem:
	sh ${ABODEMINE_WORKSPACE}/infra/cfssl/gen-ca.sh $*-ca
	chmod 0644 $*-ca-key.pem $*-ca.csr $*-ca.pem

	sudo cp $*-ca.pem /usr/local/share/ca-certificates/$*-ca.crt
	sudo update-ca-certificates

abodemine.local-peer-chain.pem abodemine.local-peer-key.pem abodemine.local-peer.csr abodemine.local-peer.pem abodemine.local-peer.pfx:
	sh ${ABODEMINE_WORKSPACE}/infra/cfssl/gen-cert.sh abodemine-ca peer abodemine.local
	chmod 0644 abodemine.local-peer-chain.pem abodemine.local-peer-key.pem abodemine.local-peer.csr abodemine.local-peer.pem abodemine.local-peer.pfx

abodemine.test-peer-chain.pem abodemine.test-peer-key.pem abodemine.test-peer.csr abodemine.test-peer.pem abodemine.test-peer.pfx:
	sh ${ABODEMINE_WORKSPACE}/infra/cfssl/gen-cert.sh abodemine-ca peer abodemine.test
	chmod 0644 abodemine.test-peer-chain.pem abodemine.test-peer-key.pem abodemine.test-peer.csr abodemine.test-peer.pem abodemine.test-peer.pfx

alpha.test-peer-chain.pem alpha.test-peer-key.pem alpha.test-peer.csr alpha.test-peer.pem alpha.test-peer.pfx:
	sh ${ABODEMINE_WORKSPACE}/infra/cfssl/gen-cert.sh abodemine-ca peer alpha.test
	chmod 0644 alpha.test-peer-chain.pem alpha.test-peer-key.pem alpha.test-peer.csr alpha.test-peer.pem alpha.test-peer.pfx

omega.test-peer-chain.pem omega.test-peer-key.pem omega.test-peer.csr omega.test-peer.pem omega.test-peer.pfx:
	sh ${ABODEMINE_WORKSPACE}/infra/cfssl/gen-cert.sh abodemine-ca peer omega.test
	chmod 0644 omega.test-peer-chain.pem omega.test-peer-key.pem omega.test-peer.csr omega.test-peer.pem omega.test-peer.pfx
