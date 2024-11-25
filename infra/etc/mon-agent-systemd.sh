cat << EOF > /etc/systemd/system/mon-agent.service
[Unit]
Description=mon-agent:$version
After=network-online.target

[Service]
Type=simple
ExecStart=/usr/local/bin/mon-agent
EnvironmentFile=/etc/mon-agent.conf
Restart=always
RestartSec=5
TimeoutStopSec=35

[Install]
WantedBy=multi-user.target
EOF

chmod 664 /etc/systemd/system/mon-agent.service
systemctl daemon-reload
systemctl enable mon-agent.service
systemctl start mon-agent.service
