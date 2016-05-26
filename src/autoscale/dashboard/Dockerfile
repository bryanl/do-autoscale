FROM node:6.2.0

ENV DASHBOARD_DIR /src/dashboard
RUN mkdir -p $DASHBOARD_DIR

RUN npm install -g ember-cli bower

ADD . $DASHBOARD_DIR

WORKDIR $DASHBOARD_DIR
RUN npm set progress=false
RUN npm install --quiet
RUN bower install --allow-root --quiet

