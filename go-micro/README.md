NAME:
   micro - A microservices toolkit

USAGE:
   micro [global options] command [command options] [arguments...]
   
VERSION:
   0.2.1
   
COMMANDS:
    api		Run the micro API
    bot		Run the micro bot
    registry	Query registry
    query	Query a service method using rpc
    stream	Query a service method using streaming rpc
    health	Query the health of a service
    stats	Query the stats of a service
    list	List items in registry
    register	Register an item in the registry
    deregister	Deregister an item in the registry
    get		Get item from registry
    sidecar	Run the micro sidecar
    new		Create a new Micro service by specifying a directory path relative to your $GOPATH
    run		Run the micro runtime
    web		Run the micro web app

GLOBAL OPTIONS:
   --client 									Client for go-micro; rpc [$MICRO_CLIENT]
   --client_request_timeout 							Sets the client request timeout. e.g 500ms, 5s, 1m. Default: 5s [$MICRO_CLIENT_REQUEST_TIMEOUT]
   --client_retries "0"								Sets the client retries. Default: 1 [$MICRO_CLIENT_RETRIES]
   --client_pool_size "0"							Sets the client connection pool size. Default: 0 [$MICRO_CLIENT_POOL_SIZE]
   --client_pool_ttl 								Sets the client connection pool ttl. e.g 500ms, 5s, 1m. Default: 1m [$MICRO_CLIENT_POOL_TTL]
   --server_name 								Name of the server. go.micro.srv.example [$MICRO_SERVER_NAME]
   --server_version 								Version of the server. 1.1.0 [$MICRO_SERVER_VERSION]
   --server_id 									Id of the server. Auto-generated if not specified [$MICRO_SERVER_ID]
   --server_address 								Bind address for the server. 127.0.0.1:8080 [$MICRO_SERVER_ADDRESS]
   --server_advertise 								Used instead of the server_address when registering with discovery. 127.0.0.1:8080 [$MICRO_SERVER_ADVERTISE]
   --server_metadata [--server_metadata option --server_metadata option]	A list of key-value pairs defining metadata. version=1.0.0 [$MICRO_SERVER_METADATA]
   --broker 									Broker for pub/sub. http, nats, rabbitmq [$MICRO_BROKER]
   --broker_address 								Comma-separated list of broker addresses [$MICRO_BROKER_ADDRESS]
   --registry 									Registry for discovery. consul, mdns [$MICRO_REGISTRY]
   --registry_address 								Comma-separated list of registry addresses [$MICRO_REGISTRY_ADDRESS]
   --selector 									Selector used to pick nodes for querying [$MICRO_SELECTOR]
   --server 									Server for go-micro; rpc [$MICRO_SERVER]
   --transport 									Transport mechanism used; http [$MICRO_TRANSPORT]
   --transport_address 								Comma-separated list of transport addresses [$MICRO_TRANSPORT_ADDRESS]
   --enable_tls									Enable TLS [$MICRO_ENABLE_TLS]
   --tls_cert_file 								TLS Certificate file [$MICRO_TLS_CERT_FILE]
   --tls_key_file 								TLS Key file [$MICRO_TLS_KEY_FILE]
   --tls_client_ca_file 							TLS CA file to verify clients against [$MICRO_TLS_CLIENT_CA_FILE]
   --api_address 								Set the api address e.g 0.0.0.0:8080 [$MICRO_API_ADDRESS]
   --proxy_address 								Proxy requests via the HTTP address specified [$MICRO_PROXY_ADDRESS]
   --sidecar_address 								Set the sidecar address e.g 0.0.0.0:8081 [$MICRO_SIDECAR_ADDRESS]
   --web_address 								Set the web UI address e.g 0.0.0.0:8082 [$MICRO_WEB_ADDRESS]
   --register_ttl "0"								Register TTL in seconds [$MICRO_REGISTER_TTL]
   --register_interval "0"							Register interval in seconds [$MICRO_REGISTER_INTERVAL]
   --api_handler 								Specify the request handler to be used for mapping HTTP requests to services; {api, proxy, rpc} [$MICRO_API_HANDLER]
   --api_namespace 								Set the namespace used by the API e.g. com.example.api [$MICRO_API_NAMESPACE]
   --sidecar_handler 								Specify the request handler to be used for mapping HTTP requests to services; {proxy, rpc} [$MICRO_SIDECAR_HANDLER]
   --sidecar_namespace 								Set the namespace used by the Sidecar e.g. com.example.srv [$MICRO_SIDECAR_NAMESPACE]
   --web_namespace 								Set the namespace used by the Web proxy e.g. com.example.web [$MICRO_WEB_NAMESPACE]
   --api_cors 									Comma separated whitelist of allowed origins for CORS [$MICRO_API_CORS]
   --web_cors 									Comma separated whitelist of allowed origins for CORS [$MICRO_WEB_CORS]
   --sidecar_cors 								Comma separated whitelist of allowed origins for CORS [$MICRO_SIDECAR_CORS]
   --enable_stats								Enable stats [$MICRO_ENABLE_STATS]
   --help, -h									show help
   --version									print the version
