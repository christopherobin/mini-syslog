# mini-syslog

A minimalistic syslog daemon I use when debugging infrastructure, output is
similar to papertrail.

## Installation

Either:

```
go get github.com/christopherobin/mini-syslog
```

Or download from the release tab on github and put the binary in your path.

## Usage

```
usage: mini-syslog [<flags>]

Mini-Syslog with papertrail-like output

Flags:
      --help               Show context-sensitive help (also try --help-long and --help-man).
  -b, --bind=":1514"       What IP:Port to bind to
  -p, --protocol=udp       Which protocol to use
      --color              Enable colors
  -t, --template=TEMPLATE  Override output template
```

## Templating

The template system uses go templates and takes the actual syslog message
as a map, standards fields provided in the message are:

* app_name
* proc_id
* msg_id
* structured_data
* client
* severity
* timestamp
* version
* hostname
* message
* priority
* facility

A couple helper functions are also provided:

* __severity__: Transforms severity integer to colored string
* __bold__: Makes text bold
* __blue__: Makes text blue
* __gray__: Makes text gray
* __green__: Makes text green
* __red__: Makes text red
* __yellow__: Makes text yellow

### Default template

```
{{ .timestamp }} [{{ severity .severity }}] {{ yellow .hostname }} {{ blue .app_name }}: {{ .message }}
```

## License

MIT