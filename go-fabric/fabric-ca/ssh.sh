
TESTCASE="intermediateca-test"
TDIR=/tmp/$TESTCASE
FABRIC_CA="$GOPATH/src/github.com/hyperledger/fabric-ca"
SCRIPTDIR="$FABRIC_CA/scripts/fvt"
ROOT_CA_ADDR=localhost
TLSDIR="$TDIR/tls"
NUMINTCAS=1

function setupTLScerts() {
   rm -rf $TLSDIR
   mkdir -p $TLSDIR
   rm -rf /tmp/CAs $TLSDIR/rootTlsCa* $TLSDIR/subTlsCa*
   export HOME=$TLSDIR
   # Root TLS CA
   $SCRIPTDIR/utils/pki -f newca -a rootTlsCa -t ec -l 256 -d sha256 \
                        -n "/C=US/ST=NC/L=RTP/O=IBM/O=Hyperledger/OU=FVT/CN=localhost/" -S "IP:127.0.0.1" \
                        -K "digitalSignature,nonRepudiation,keyEncipherment,dataEncipherment,keyAgreement,keyCertSign,cRLSign" \
                        -E "serverAuth,clientAuth,codeSigning,emailProtection,timeStamping" \
                        -e 20370101000000Z -s 20160101000000Z -p rootTlsCa- >/dev/null 2>&1
#    # Sub TLS CA
   $SCRIPTDIR/utils/pki -f newsub -b subTlsCa -a rootTlsCa -t ec -l 256 -d sha256 \
                        -n "/C=US/ST=NC/L=RTP/O=IBM/O=Hyperledger/OU=FVT/CN=subTlsCa/" -S "IP:127.0.0.1" \
                        -K "digitalSignature,nonRepudiation,keyEncipherment,dataEncipherment,keyAgreement,keyCertSign,cRLSign" \
                        -E "serverAuth,clientAuth,codeSigning,emailProtection,timeStamping" \
                        -e 20370101000000Z -s 20160101000000Z -p subTlsCa- >/dev/null 2>&1
#    # EE TLS certs
#     i=0;while test $((i++)) -lt $NUMINTCAS; do
#     rm -rf $TLSDIR/intFabCaTls${i}*
#    $SCRIPTDIR/utils/pki -f newcert -a subTlsCa -t ec -l 256 -d sha512 \
#                         -n "/C=US/ST=NC/L=RTP/O=IBM/O=Hyperledger/OU=FVT/CN=intFabCaTls${i}/" -S "IP:127.0.${i}.1" \
#                         -K "digitalSignature,nonRepudiation,keyEncipherment,dataEncipherment,keyAgreement,keyCertSign,cRLSign" \
#                         -E "serverAuth,clientAuth,codeSigning,emailProtection,timeStamping" \
#                         -e 20370101000000Z -s 20160101000000Z -p intFabCaTls${i}- >/dev/null 2>&1 <<EOF
# EOF
#    done
#  cat $TLSDIR/rootTlsCa-cert.pem $TLSDIR/subTlsCa-cert.pem > $TLSDIR/tlsroots.pem
}

setupTLScerts
