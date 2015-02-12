vcl 4.0;

backend default {
    .host = "127.0.0.1";
    .port = "8080";
}

# Remove all cookies
sub vcl_recv {
    unset req.http.cookie;
}

# Remove all cookies
sub vcl_backend_response {
    unset beresp.http.set-cookie;
}
