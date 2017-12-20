Feature: Manage proxies wit API.
    As an authenticated admin user, I need to be able to manage proxies list.

    Scenario: API list must not be available w/out valid JWT token
        Given request JWT token is not set
        When I request "/apis" API path with "GET" method
        Then I should receive 401 response code

        Given request header "Authorization" is set to "Basic YWRtaW46YWRtaW4="
        When I request "/apis" API path with "GET" method
        Then I should receive 401 response code

        Given request header "Authorization" is set to "Bearer InvalidJWT"
        When I request "/apis" API path with "GET" method
        Then I should receive 401 response code

    Scenario: APIs list and create must be available for user with correct admin token
        Given request JWT token is valid admin token
        When I request "/apis" API path with "GET" method
        Then I should receive 200 response code
        And response JSON body is an array of length 0

        Given request JSON payload '{"name":"example","active":true,"proxy":{"preserve_host":false,"listen_path":"/example/*","upstreams": {"balancing": "roundrobin", "targets": [{"target": "http://localhost:9089/hello-world"}]},"strip_path":false,"append_path":false,"enable_load_balancing":false,"methods":["GET"]},"health_check":{"url":"https://example.com/status"}}'
        When I request "/apis" API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/example"

        When I request "/apis" API path with "GET" method
        Then I should receive 200 response code
        And response JSON body is an array of length 1
        And response JSON body has "0.name" path with value 'example'
        And response JSON body has "0.active" path with value 'true'

        Given request JSON payload '{"name":"posts","active":true,"proxy":{"preserve_host":false,"listen_path":"/posts/*","upstreams": {"balancing": "roundrobin", "targets": [{"target": "http://localhost:9089/posts"}]},"strip_path":true,"append_path":false,"enable_load_balancing":false,"methods":["ALL"],"hosts":["hellofresh.*"]},"plugins":[{"name":"cors","enabled":true,"config":{"domains":["*"],"methods":["GET","POST","PUT","PATCH","DELETE"],"request_headers":["Origin","Authorization","Content-Type"],"exposed_headers":["X-Debug-Token","X-Debug-Token-Link"]}},{"name":"rate_limit","enabled":true,"config":{"limit":"10-S","policy":"local"}},{"name":"oauth2","enabled":true,"config":{"server_name":"local"}},{"name":"compression","enabled":true}]}'
        When I request "/apis" API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/posts"

        When I request "/apis" API path with "GET" method
        Then I should receive 200 response code
        And response JSON body is an array of length 2

    Scenario: API fails to create routes with the same name
        Given request JWT token is valid admin token
        And request JSON payload '{"name":"example","active":true,"proxy":{"preserve_host":false,"listen_path":"/example/*","upstreams": {"balancing": "roundrobin", "targets": [{"target": "http://localhost:9089/hello-world"}]},"strip_path":false,"append_path":false,"enable_load_balancing":false,"methods":["GET"]},"health_check":{"url":"https://example.com/status"}}'
        When I request "/apis" API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/example"

        Given request JSON payload '{"name":"example","active":true,"proxy":{"preserve_host":false,"listen_path":"/example/*","upstreams": {"balancing": "roundrobin", "targets": [{"target": "http://localhost:9089/hello-world"}]},"strip_path":false,"append_path":false,"enable_load_balancing":false,"methods":["GET"]},"health_check":{"url":"https://example.com/status"}}'
        When I request "/apis" API path with "POST" method
        Then I should receive 409 response code
        And the response should contain "api name is already registered"

        Given request JSON payload '{"name":"example1","active":true,"proxy":{"preserve_host":false,"listen_path":"/example/*","upstreams": {"balancing": "roundrobin", "targets": [{"target": "http://localhost:9089/hello-world"}]},"strip_path":false,"append_path":false,"enable_load_balancing":false,"methods":["GET"]},"health_check":{"url":"https://example.com/status"}}'
        When I request "/apis" API path with "POST" method
        Then I should receive 409 response code
        And the response should contain "api listen path is already registered"

    Scenario: API must return existing routes and response with error for non-existent
        Given request JWT token is valid admin token
        When I request "/apis/example" API path with "GET" method
        Then I should receive 404 response code

        Given request JSON payload '{"name":"example","active":true,"proxy":{"preserve_host":false,"listen_path":"/example/*","upstreams": {"balancing": "roundrobin", "targets": [{"target": "http://localhost:9089/hello-world"}]},"strip_path":false,"append_path":false,"enable_load_balancing":false,"methods":["GET"]},"health_check":{"url":"https://example.com/status"}}'
        When I request "/apis" API path with "POST" method
        Then I should receive 201 response code

        When I request "/apis/example" API path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "name" path with value 'example'

    Scenario: API must update existing routes with new path value
        Given request JWT token is valid admin token
        And request JSON payload '{"name":"example","active":true,"proxy":{"preserve_host":false,"listen_path":"/example/*","upstreams": {"balancing": "roundrobin", "targets": [{"target": "http://localhost:9089/hello-world"}]},"strip_path":false,"append_path":false,"enable_load_balancing":false,"methods":["GET"]},"health_check":{"url":"https://example.com/status"}}'
        When I request "/apis" API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/example"

        When I request "/apis/example" API path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "name" path with value 'example'
        And response JSON body has "proxy.listen_path" path with value '/example/*'

        Given request JSON payload '{"name":"example","active":true,"proxy":{"preserve_host":false,"listen_path":"/example1/*","upstreams": {"balancing": "roundrobin", "targets": [{"target": "http://localhost:9089/hello-world"}]},"strip_path":false,"append_path":false,"enable_load_balancing":false,"methods":["GET"]},"health_check":{"url":"https://example.com/status"}}'
        When I request "/apis/example" API path with "PUT" method
        Then I should receive 200 response code

        When I request "/apis/example" API path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "name" path with value 'example'
        And response JSON body has "proxy.listen_path" path with value '/example1/*'

        When I request "/apis/foo-bar" API path with "GET" method
        Then I should receive 404 response code

    Scenario: API must delete existing routes
        Given request JWT token is valid admin token
        And request JSON payload '{"name":"example","active":true,"proxy":{"preserve_host":false,"listen_path":"/example/*","upstreams": {"balancing": "roundrobin", "targets": [{"target": "http://localhost:9089/hello-world"}]},"strip_path":false,"append_path":false,"enable_load_balancing":false,"methods":["GET"]},"health_check":{"url":"https://example.com/status"}}'
        When I request "/apis" API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/example"

        Given request JSON payload '{"name":"posts","active":true,"proxy":{"preserve_host":false,"listen_path":"/posts/*","upstreams": {"balancing": "roundrobin", "targets": [{"target": "http://localhost:9089/posts"}]},"strip_path":true,"append_path":false,"enable_load_balancing":false,"methods":["ALL"],"hosts":["hellofresh.*"]},"plugins":[{"name":"cors","enabled":true,"config":{"domains":["*"],"methods":["GET","POST","PUT","PATCH","DELETE"],"request_headers":["Origin","Authorization","Content-Type"],"exposed_headers":["X-Debug-Token","X-Debug-Token-Link"]}},{"name":"rate_limit","enabled":true,"config":{"limit":"10-S","policy":"local"}},{"name":"oauth2","enabled":true,"config":{"server_name":"local"}},{"name":"compression","enabled":true}]}'
        When I request "/apis" API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/posts"

        When I request "/apis" API path with "GET" method
        Then I should receive 200 response code
        And response JSON body is an array of length 2

        When I request "/apis/example" API path with "DELETE" method
        Then I should receive 204 response code

        When I request "/apis" API path with "GET" method
        Then I should receive 200 response code
        And response JSON body is an array of length 1

        When I request "/apis/example" API path with "DELETE" method
        Then I should receive 404 response code

        When I request "/apis/posts" API path with "DELETE" method
        Then I should receive 204 response code

        When I request "/apis" API path with "GET" method
        Then I should receive 200 response code
        And response JSON body is an array of length 0

        When I request "/apis/posts" API path with "DELETE" method
        Then I should receive 404 response code
