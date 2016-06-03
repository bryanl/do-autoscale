import Ember from 'ember';
import config from './config/environment';

const Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  this.route('templates');
  this.route('groups', function() {
    this.route('new');
    this.route('show', { path: '/:group_id'}, function() {
      this.route('status', { path: '/summary' });
      this.route('policy');
      this.route('delete');
    });
  });
});

export default Router;
