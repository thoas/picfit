vcl 4.1;

import std;

backend default {
    .host = "127.0.0.1";
    .port = "3001";
}

acl purgers {
    "127.0.0.1";
}

sub vcl_backend_response {
    if (beresp.status >= 400) {
        set beresp.ttl = 1s;
        set beresp.uncacheable = true;
        return (deliver);
    }

	if (beresp.http.cache-control ~ "must-revalidate") {
		set beresp.ttl = 1s;
		set beresp.grace = 0s;
		set beresp.keep = 1w;
	} else {
		set beresp.ttl = 1w;
	}

    return (deliver);
}

sub vcl_hit {
    if (req.method == "PURGE") {
        return (pass);
    }

    if (req.method == "DELETE") {
        return (pass);
    }

    if (req.method == "BAN") {
        return (pass);
    }

    if (req.method == "POST") {
        return (pass);
    }
}

sub vcl_miss {
    if (req.method == "PURGE") {
        return (pass);
    }

    if (req.method == "DELETE") {
        return (pass);
    }

    if (req.method == "BAN") {
        return (pass);
    }

    if (req.method == "POST") {
        return (pass);
    }
}

# Remove all cookies
sub vcl_backend_response {
    unset beresp.http.set-cookie;
}

sub vcl_pass {
    if (req.method == "PURGE") {
        return(synth(502, "PURGE on a passed object"));
    }
}

sub vcl_recv {
    if (req.url ~ "^/(get|stats|healthcheck|debug|sys)") {
        return (pass);
    }

    if (req.method == "DELETE" || req.method == "POST") {
		return (pass);
	}

    unset req.http.cookie;

    if (req.method == "PURGE") {
        if (!client.ip ~ purgers) {
            return(synth(405, "Method not allowed"));
        }
        return (purge);
    }

	if (req.method == "BAN") {
		if (!client.ip ~ purgers) {
			return(synth(403, "Not allowed."));
		}

		ban("req.http.host == " + req.http.host +
				" && req.url == " + req.url);

		return(synth(200, "Ban added"));
	}
}

sub vcl_deliver {
    # Was a HIT or a MISS?
    if (obj.hits > 0) {
        set resp.http.X-Cache = "HIT";
    } else {
        set resp.http.X-Cache = "MISS";
    }

    # And add the number of hits in the header:
    set resp.http.X-Cache-Hits = obj.hits;
}
