openssl x509 -in fabric-ca-cert.pem -noout -issuer -subject -serial -dates -nameopt RFC2253| sed 's/^/   /'
     openssl x509 -in fabric-ca-cert.pem -noout -text |
        awk '
           /Subject Alternative Name:/ {
              gsub(/^ */,"")
              printf $0"= "
              getline; gsub(/^ */,"")
              print
           }'| sed 's/^/   /'
     openssl x509 -in fabric-ca-cert.pem -noout -pubkey |
        openssl ecdsa -pubin -noout -text 2>/dev/null| sed 's/Private/Public/'
     openssl ecdsa -in fabric-ca-key.pem -text 2>/dev/null


//rsa|ecdsa


FABRIC_CA_CERT_FILE="$FABRIC_CA_CLIENT_HOME/msp/signcerts/cert.pem"
FABRIC_CA_KEY_FILE="$FABRIC_CA_CLIENT_HOME/msp/keystore/key.pem"

printAuth $FABRIC_CA_CERT_FILE $FABRIC_CA_KEY_FILE

printAuth() {
   local CLIENTCERT="$1"
   local CLIENTKEY="$2"

   : ${CLIENTCERT:="$HOME/fabric-ca/cert.pem"}
   : ${CLIENTKEY:="$HOME/fabric-ca/key.pem"}

   echo CERT:
   openssl x509 -in $FABRIC_CA_CERT_FILE -text 2>&1 | sed 's/^/    /'
   type=$(cat $FABRIC_CA_KEY_FILE | head -n1 | awk '{print tolower($2)}')
   test -z "$type" && type=rsa
   echo KEY:
   openssl $type -in $CLIENTKEY -text 2>/dev/null| sed 's/^/    /'
}
