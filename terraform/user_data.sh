#!/bin/bash
set -e

# Update system
yum update -y

# Install Docker
yum install -y docker
systemctl start docker
systemctl enable docker

# Add ec2-user to docker group
usermod -a -G docker ec2-user

# Install Docker Compose
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Create application directory
mkdir -p /opt/cryple
cd /opt/cryple

# Create docker-compose.yml for production
cat > docker-compose.yml << 'EOF'
version: "3.9"

services:
  cryple:
    image: philippeberto/cryple-api:0.0.1
    ports:
      - "8080:8080"
    environment:
      - POSTGRES_HOST=${db_host}
      - POSTGRES_USER=${db_user}
      - POSTGRES_PASSWORD=${db_password}
      - POSTGRES_DB=${db_name}
      - POSTGRES_PORT=5432
      - POSTGRES_DRIVER=postgres
      - POSTGRES_TIMEOUT=30
      - POSTGRES_IDLE_CONNECTION=10
      - POSTGRES_LIFE_TIME=3600
      - POSTGRES_OPEN_CONNECTION=25
      - POSTGRES_MIGRATION=true
      - DEBUG=false
      - PORT=8080
      - APP_SERVICE=cryple_general
      - APP_NAME=cryple
      - METRICS_PORT=80
      - METRICS_ENABLE=false
      - TRACE_ENABLE=false
      - ENABLE_CORS=true
      - CORS_ALLOW_ORIGINS=*
    restart: unless-stopped
EOF

# Create systemd service for the application
cat > /etc/systemd/system/cryple.service << 'EOF'
[Unit]
Description=Cryple Application
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/cryple
ExecStart=/usr/local/bin/docker-compose up -d
ExecStop=/usr/local/bin/docker-compose down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
EOF

# Enable and start the service
systemctl daemon-reload
systemctl enable cryple.service

# Wait a bit for Docker to be fully ready
sleep 30

# Start the application
systemctl start cryple.service

# Create a simple health check script
cat > /opt/cryple/health-check.sh << 'EOF'
#!/bin/bash
response=$(curl -s -o /dev/null -w "%%{http_code}" http://localhost:8080/users/check || echo "000")
if [ "$response" -eq 400 ] || [ "$response" -eq 200 ]; then
    echo "Application is healthy (HTTP $response)"
    exit 0
else
    echo "Application is unhealthy (HTTP $response)"
    exit 1
fi
EOF

chmod +x /opt/cryple/health-check.sh

# Log the deployment
echo "Cryple application deployed successfully at $(date)" >> /var/log/cryple-deploy.log
