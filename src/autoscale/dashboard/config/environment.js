/* jshint node: true */

module.exports = function(environment) {
  var ENV = {
    modulePrefix: 'dashboard',
    environment: environment,
    baseURL: '/',
    locationType: 'hash',
    EmberENV: {
      FEATURES: {
        // Here you can enable experimental features on an ember canary build
        // e.g. 'with-controller': true
      }
    },

    APP: {
      // Here you can pass flags/options to your application instance
      // when it is created
      websocketURL: "ws://localhost:8888/api/notifications"
    }
  };

  if (environment === 'development') {
    // This needs to be changed based on your current dev configuration
    ENV.APP.websocketURL = "ws://192.81.208.211:8888/api/notifications";

    // ENV.APP.LOG_RESOLVER = true;
    // ENV.APP.LOG_ACTIVE_GENERATION = true;
    // ENV.APP.LOG_TRANSITIONS = true;
    // ENV.APP.LOG_TRANSITIONS_INTERNAL = true;
    // ENV.APP.LOG_VIEW_LOOKUPS = true;
  }

  if (environment === 'test') {
    // Testem prefers this...
    ENV.baseURL = '/';
    ENV.locationType = 'none';

    // keep test console output quieter
    ENV.APP.LOG_ACTIVE_GENERATION = false;
    ENV.APP.LOG_VIEW_LOOKUPS = false;

    ENV.APP.rootElement = '#ember-testing';
  }

  if (environment === 'production') {
    ENV.baseURL = "/dashboard/";
    // This needs to be changed based on your current production configuration
    ENV.APP.websocketURL = "wss://autoscale-dev.x.pifft.com/api/notifications";
  }

  return ENV;
};
