package config

var yannotated = `# Handlers know how to send notifications to specific services.
handler:
  slack:
    # Slack "legacy" API token.
    token: ""
    # Slack channel.
    channel: ""
    # Title of the message.
    title: ""
  hipchat:
    # Hipchat token.
    token: ""
    # Room name.
    room: ""
    # URL of the hipchat server.
    url: ""
  mattermost:
    room: ""
    url: ""
    username: ""
  flock:
    # URL of the flock API.
    url: ""
  webhook:
    # Webhook URL.
    url: ""
  msteams:
    # MSTeams API Webhook URL.
    webhookurl: ""
  smtp:
    # Destination e-mail address.
    to: ""
    # Sender e-mail address .
    from: ""
    # Smarthost, aka "SMTP server"; address of server used to send email.
    smarthost: ""
    # Subject of the outgoing emails.
    subject: ""
    # Extra e-mail headers to be added to all outgoing messages.
    headers: {}
    # Authentication parameters.
    auth:
      # Username for PLAN and LOGIN auth mechanisms.
      username: ""
      # Password for PLAIN and LOGIN auth mechanisms.
      password: ""
      # Identity for PLAIN auth mechanism
      identity: ""
      # Secret for CRAM-MD5 auth mechanism
      secret: ""
    # If "true" forces secure SMTP protocol (AKA StartTLS).
    requireTLS: false
    # SMTP hello field (optional)
    hello: ""
# Resources to watch.
resource:
  deployment: false
  rc: false
  rs: false
  ds: false
  svc: false
  po: false
  job: false
  node: false
  clusterrole: false
  sa: false
  pv: false
  ns: false
  secret: false
  configmap: false
  ing: false
# For watching specific namespace, leave it empty for watching all.
# this config is ignored when watching namespaces
namespace: ""
`
