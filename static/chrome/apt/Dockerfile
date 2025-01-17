FROM browsers/base:7.4.2

ARG VERSION
ARG PACKAGE=google-chrome-stable
ARG INSTALL_DIR=chrome

LABEL browser=$PACKAGE:$VERSION

# https://pptr.dev/troubleshooting#chrome-doesnt-launch-on-linux
RUN \
    apt-get update && \
    apt-get -y --no-install-recommends install unzip \
    gconf-service libasound2 libatk1.0-0 libc6 libcairo2 \
    libcups2 libdbus-1-3 libexpat1 libfontconfig1 libgbm1 libgcc1 libgconf-2-4 \
    libgdk-pixbuf2.0-0 libglib2.0-0 libgtk-3-0 libnspr4 libpango-1.0-0 \
    libpangocairo-1.0-0 libstdc++6 libx11-6 libx11-xcb1 libxcb1 libxcomposite1 \
    libxcursor1 libxdamage1 libxext6 libxfixes3 libxi6 libxrandr2 libxrender1 \
    libxss1 libxtst6 ca-certificates fonts-liberation libnss3 lsb-release \
    xdg-utils

# https://stackoverflow.com/a/78880055
RUN CHROME_VERSION=$(curl -sSL https://googlechromelabs.github.io/chrome-for-testing/ | awk -F 'Version:' '/Stable/ {print $2}' | awk '{print $1}' | sed 's/<code>//g; s/<\/code>//g') && \
    CHROME_URL="https://storage.googleapis.com/chrome-for-testing-public/${CHROME_VERSION}/linux64/chrome-linux64.zip" && \
    echo "Fetching Chrome version: ${CHROME_VERSION}" && \
    curl -sSL ${CHROME_URL} -o /tmp/chrome-linux64.zip && \
    mkdir -p /opt/google/chrome && \
    mkdir -p /usr/local/bin && \
    unzip -q /tmp/chrome-linux64.zip -d /opt/google && \
    mv /opt/google/chrome-linux64/* /opt/google/chrome && \
    echo "ls -la /opt/google/chrome" && \
    ls -la /opt/google/chrome && \
    chown root:root /opt/google/chrome/chrome_sandbox && \
    chmod 4755 /opt/google/chrome/chrome_sandbox && \
    ln -s /opt/google/chrome/chrome /usr/local/bin/google-chrome && \
    google-chrome --version && \
    rm /tmp/chrome-linux64.zip

