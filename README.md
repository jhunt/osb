OSB - A CLI for the Open Service Broker API
===========================================

```
USAGE: ./osb [options] <command> [options]

Options:

  -h, --help         Show the help screen.

  -T, --trace        Trace HTTP requests and responses as they happen.
                     Can also be enabled by setting OSB_TRACE=yes.

  --data             Path to the OSB data file, for storing instance and
                     binding information required by future bind, unbind,
                     and deprovision requests.

  -e, --endpoint     The URL to the backend service broker to interact with.
                     Can also be specified via the OSB_URL variable.

  -U, --username     The username for service broker HTTP Basic Auth.
                     Can also be specified via the OSB_USERNAME variable.

  -P, --password     The password for service broker HTTP Basic Auth.
                     Can also be specified via the OSB_PASSWORD variable.

  -k, --skip-verify  Do not validate X.509 TLS certificates.
                     Can also be specified via OSB_SKIP_VERIFY.

  -t, --timeout      Timeout (in seconds) for HTTP reuests.
                     Can also be specified via OSB_TIMEOUT.

  --json             Emit JSON responses, and nothing else.
                     Useful for scripting!

Commands:

  list           List known instance and binding details, from ~/.osbrc.
  catalog        Retrieve the service catalog from the service broker.

  provision      Provision a new instance of a service/plan.
  deprovision    Remove a provsioned instance.

  bind           Bind a provisioned instance, to get credentials.
  unbind         Unbind an instance, releasing bound credentials.

```
