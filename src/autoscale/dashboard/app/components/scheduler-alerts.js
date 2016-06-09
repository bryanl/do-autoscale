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
    const socket = this.get('websockets').socketFor(this.websocketURL);
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

  myOpenHandler: function () {
  },

  myMessageHandler: function (event) {
    console.log('Message: ' + event.data);
    var obj = JSON.parse(event.data);
    console.log(obj);
    this.get('notifications').push(obj);
  },

  actions: {
  }
});
