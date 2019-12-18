# image-resizer
 
config variables:

	"relative-path" - relative path to endpoint
	"http-port" - port which must be used by http server
	"download-queue-size" - load balanse variable. 
	describes how much cuncurrent requests service can handle. after this count exceeds - all clients would receive error
	"incoming-timeout" - maximum processing await timer. 
	if service cant process image after this timer expires - client would receive error.
	"download-timeout" - timeout for image downloading from third-party resource.

deployments/docker-compose.yml - currently creates 2 resizer services and nginx load-balancer.
