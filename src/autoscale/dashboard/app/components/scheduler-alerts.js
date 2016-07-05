/* jshint node: true */

import config from '../config/environment';
import Ember from 'ember';

export default Ember.Component.extend({
  websockets: Ember.inject.service(),
  socketRef: null,
  websocketURL: config.APP.websocketURL,

  willRender() {
    this.websocketConnect();
  },

  websocketConnect() {
    // TODO get this URL from somewhere else
    var self = this;

    var websocketURL = this.websocketURL;
    if (websocketURL.length == 0) {
      var l = window.location;
      websocketURL = "wss://" + l.host + "/api/notifications";
    }

    const socket = this.get('websockets').socketFor(websocketURL);
    socket.on('open', this.myOpenHandler, this);
    socket.on('message', this.myMessageHandler, this);
    socket.on('close', () => {
      setTimeout(() => {
        console.info("reconnecting to websocket");
        self.websocketConnect();
      }, Math.floor(Math.random() * 3000));
    }, this);

    this.set('socketRef', socket);
  },

  notifications: [],
  notificationEntries: Ember.computed('notifications', function () {
    var list = [];
    const notifications = this.get('notifications');
    for (let notif of notifications) {
      var createdAt = new Date(Date.parse(notif.createdAt));
      var createdAtUTC = new Date(createdAt.getUTCFullYear(), createdAt.getUTCMonth(), createdAt.getUTCDate(), createdAt.getUTCHours(), createdAt.getUTCMinutes(), createdAt.getUTCSeconds());
      var dateStr = createdAtUTC.toLocaleString('en-US', { hour12: false });

      var action = "grew";
      if (notif.delta < 0) {
        action = "shrank"
      }
      const msg = `${notif.name} ${action} to ${notif.count}`;

      list.push({ id: notif.groupID, msg: msg });
    }

    return list;
  }),

  badgeHasAlerts: Ember.computed('notifications', function() {
    const notifications = this.get('notifications');
    if (notifications.length > 0) {
      return "badge-with-alerts";
    }

    return "badge-with-no-alerts";
  }),

  myOpenHandler: function (event) {
  },

  myMessageHandler: function (event) {
    var notifications = this.get('notifications');
    var obj = JSON.parse(event.data);
    notifications.addObject(obj);
    this.notifyPropertyChange('notifications');
  },

  actions: {
  }
});
